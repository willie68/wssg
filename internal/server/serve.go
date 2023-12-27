package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/willie68/wssg/internal/generator"
	"github.com/willie68/wssg/internal/logging"
)

// Server this is the http server with update capatibiliries
type Server struct {
	rootFolder string
	log        *logging.Logger
	gen        generator.Generator
	output     string
}

// New creates a new http server with auto update capatibilities
func New(rootFolder string, gen generator.Generator) Server {
	s := Server{
		rootFolder: rootFolder,
		log:        logging.New().WithName("server"),
		gen:        gen,
	}
	s.init()
	return s
}

func (s *Server) init() {
	s.output = filepath.Join(s.rootFolder, s.gen.GenConfig().Output)
	// Starting the file watcher on the output folder
}

// StartWatcher starting a file system watcher
func (s *Server) StartWatcher() error {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				s.log.Infof("event: %v", event)
				if event.Has(fsnotify.Write) {
					s.log.Infof("modified file: %s", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				s.log.Errorf("error: %v", err)
			}
		}
	}()

	err = filepath.Walk(s.rootFolder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}
		if info.IsDir() {
			s.log.Infof("adding watch path: %s", path)
			watcher.Add(path)
		}
		return nil
	})
	return err
}

// Serve starting the http server serving the files
func (s *Server) Serve() error {

	err := s.StartWatcher()
	if err != nil {
		return err
	}

	fileServer := http.FileServer(http.Dir(s.output))
	http.Handle("/", fileServer)
	s.log.Info("start serving site. use http://localhost:8080/index.html")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		return err
	}
	return nil
}
