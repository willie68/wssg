package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/willie68/wssg/internal/config"
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
	sy         sync.Mutex
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
	output := filepath.Join(s.rootFolder, s.gen.GenConfig().Output)
	output, err := filepath.Abs(output)
	if err != nil {
		s.log.Errorf("error converting relativ path to absolute: %v", err)
	}
	s.output = output
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
		if info.IsDir() && (strings.HasPrefix(info.Name(), ".") && (info.Name() != config.WssgFolder)) {
			return filepath.SkipDir
		}
		absPath, err := filepath.Abs(path)
		if (s.output == absPath) || (err != nil) {
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
			if strings.HasPrefix(fn, "_") {
				continue
			}
			s.log.Infof("event: %v, file: %s", event, fn)
			if event.Has(fsnotify.Write) {
				s.generate(event.Name)
			}
			absPath, err := filepath.Abs(event.Name)
			if err == nil && s.output == absPath {
				if event.Has(fsnotify.Remove) {
					s.generate(event.Name)
				}
			}
		case err, ok := <-s.watcher.Errors:
			if !ok {
				return
			}
			s.log.Errorf("error: %v", err)
		}
	}
}

func (s *Server) generate(name string) {
	s.sy.Lock()
	s.log.Infof("modified file: %s", name)
	err := s.gen.Execute()
	if err != nil {
		s.log.Errorf("error generate site: %v", err)
	}
	s.sy.Unlock()
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
	http.HandleFunc("/_refresh", func(w http.ResponseWriter, r *http.Request) {
		refresh := struct {
			Refresh bool `json:"refresh"`
		}{
			Refresh: false,
		}
		if s.gen.IsRefreshed() {
			refresh.Refresh = true
		}
		dst, err := json.Marshal(refresh)
		if err != nil {
			msg := fmt.Sprintf("error output refresh: %v", err)
			s.log.Error(msg)
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(msg))
			if err != nil {
				s.log.Errorf("error output refresh: %v", err)
			}
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(dst)
		if err != nil {
			s.log.Errorf("error output refresh: %v", err)
		}
	})
	s.log.Info("start serving site. use http://localhost:8080/index.html for the result. Stopping server with ctrl+c.")
	open("http://localhost:8080/index.html")
	return http.ListenAndServe(":8080", nil)
}

// https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
