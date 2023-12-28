package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

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
	watcher    *fsnotify.Watcher
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
	var err error
	// Create new watcher.
	s.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// Start listening for events.
	go s.doEvent()
	absPath, err := filepath.Abs(s.rootFolder)
	if err != nil {
		return err
	}
	err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}
		if info.IsDir() {
			s.log.Infof("adding watch path: %s", path)
			s.watcher.Add(path)
		}
		return nil
	})
	return err
}

func (s *Server) doEvent() {
	var sy sync.Mutex
	for {
		select {
		case event, ok := <-s.watcher.Events:
			if !ok {
				return
			}
			fn := filepath.Base(event.Name)
			if strings.HasPrefix(fn, ".") {
				continue
			}
			s.log.Infof("event: %v, file: %s", event, fn)
			if event.Has(fsnotify.Write) {
				sy.Lock()
				s.log.Infof("modified file: %s", event.Name)
				err := s.gen.Execute()
				if err != nil {
					s.log.Errorf("error generate site: %v", err)
				}
				sy.Unlock()
			}
		case err, ok := <-s.watcher.Errors:
			if !ok {
				return
			}
			s.log.Errorf("error: %v", err)
		}
	}
}

// Serve starting the http server serving the files
func (s *Server) Serve() error {
	err := s.StartWatcher()
	if err != nil {
		return err
	}
	defer s.watcher.Close()

	fileServer := http.FileServer(http.Dir(s.output))
	http.Handle("/", fileServer)
	s.log.Info("start serving site. use http://localhost:8080/index.html for the result. Stopping server with ctrl+c.")
	return http.ListenAndServe(":8080", nil)
}
