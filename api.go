package chatterbox

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/schnapper79/chatterbox/types"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Router       *mux.Router
	ModelPath    string
	PathToLLama  string
	LoadedModels map[string]*Runner
	Server       *http.Server
	usedPorts    map[int]bool
}

func (s *Server) loadModelHandler(w http.ResponseWriter, r *http.Request) {
	// Access the model value from the path
	vars := mux.Vars(r)
	modelname := vars["model"]
	if _, ok := s.LoadedModels[modelname]; ok {
		http.Error(w, "Model already loaded", http.StatusBadRequest)
		return
	}

	req := types.NewModelRequestWithDefaults()
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.ModelName = modelname

	if _, ok := s.usedPorts[req.Port]; ok {
		http.Error(w, "Port already in use", http.StatusBadRequest)
		return
	}

	ctx, Cancel := context.WithCancel(context.Background())

	newRunner := NewRunner(ctx, Cancel, s.PathToLLama, s.ModelPath, req)
	//Load model
	err = newRunner.Run()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Logging
	go func() {
		for msg := range newRunner.LogChan {
			logger.Info(msg)
		}
	}()

	// Error handling
	go func() {
		for err := range newRunner.ErrorChan {
			logger.Error(err)
		}
		Cancel()
		delete(s.usedPorts, s.LoadedModels[modelname].Config.Port)
		delete(s.LoadedModels, modelname)
	}()

	s.usedPorts[req.Port] = true
	s.LoadedModels[req.ModelName] = newRunner

	w.Write([]byte("Model loaded"))
}
func (s *Server) unloadModelHandler(w http.ResponseWriter, r *http.Request) {
	// Access the model value from the path
	vars := mux.Vars(r)
	modelname := vars["model"]

	//Unload model
	if _, ok := s.LoadedModels[modelname]; !ok {
		http.Error(w, "Model not loaded", http.StatusBadRequest)
		return
	}

	delete(s.usedPorts, s.LoadedModels[modelname].Config.Port)
	delete(s.LoadedModels, modelname)

}

func (s *Server) getLoadedModelsHandler(w http.ResponseWriter, r *http.Request) {
	models := map[string]*types.Model_Request{}
	for k := range s.LoadedModels {
		models[k] = s.LoadedModels[k].Config
	}
	json.NewEncoder(w).Encode(models)
}

func (s *Server) getAvailableModelsHandler(w http.ResponseWriter, r *http.Request) {
	//read available models from disk
	models := []string{}

	filetype := ".gguf" // Replace with your desired file extension

	files, err := os.ReadDir(s.ModelPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), filetype) {
			models = append(models, file.Name())
		}
	}

	json.NewEncoder(w).Encode(models)

}

func (s *Server) saveModelHandler(w http.ResponseWriter, r *http.Request) {
	// Access the model value from the path
	vars := mux.Vars(r)
	modelname := vars["model"]

	if _, ok := s.LoadedModels[modelname]; !ok {
		http.Error(w, "Model not loaded", http.StatusBadRequest)
		return
	}

	err := s.LoadedModels[modelname].Config.Save(s.ModelPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Model saved"))

}

func (s *Server) loadModelFromFileHandler(w http.ResponseWriter, r *http.Request) {
	// Access the model value from the path
	vars := mux.Vars(r)
	konfigname := vars["konfig"]

	config := types.NewModelRequestWithDefaults()
	err := config.Load(s.ModelPath, konfigname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	modelname := config.ModelName
	if _, ok := s.LoadedModels[modelname]; ok {
		http.Error(w, "Model already loaded", http.StatusBadRequest)
		return
	}

	res, err := s.LoadModellFromFile(modelname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//return the config
	json.NewEncoder(w).Encode(res.Config)
}

func (s *Server) LoadModellFromFile(modelname string) (*Runner, error) {

	req := types.NewModelRequestWithDefaults()
	err := req.Load(s.ModelPath, modelname)
	if err != nil {
		return nil, err
	}
	req.ModelName = modelname

	if _, ok := s.usedPorts[req.Port]; ok {
		return nil, fmt.Errorf("Port %d already in use", req.Port)
	}

	ctx, Cancel := context.WithCancel(context.Background())
	newRunner := NewRunner(ctx, Cancel, s.PathToLLama, s.ModelPath, req)
	//Load model
	err = newRunner.Run()
	if err != nil {
		return nil, err
	}
	// Logging
	go func() {
		for msg := range newRunner.LogChan {
			logger.Info(msg)
		}
	}()

	// Error handling
	go func() {
		for err := range newRunner.ErrorChan {
			logger.Error(err)
		}
		Cancel()
		delete(s.usedPorts, s.LoadedModels[modelname].Config.Port)
		delete(s.LoadedModels, modelname)
	}()

	s.usedPorts[req.Port] = true
	s.LoadedModels[modelname] = newRunner
	return newRunner, nil

}
func (s *Server) getAvailableKonfigsHandler(w http.ResponseWriter, r *http.Request) {
	konfigs := []string{}

	filetype := ".json" // Replace with your desired file extension

	files, err := os.ReadDir(s.ModelPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), filetype) {
			//remove extension
			fn := strings.TrimSuffix(file.Name(), filetype)
			konfigs = append(konfigs, fn)
		}
	}

	json.NewEncoder(w).Encode(konfigs)
}
func (s *Server) GetKonfigHandler(w http.ResponseWriter, r *http.Request) {
	// Access the model value from the path
	vars := mux.Vars(r)
	konfigname := vars["konfig"]

	config := types.NewModelRequestWithDefaults()
	err := config.Load(s.ModelPath, konfigname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(config)
}

func (s *Server) SaveKonfigHandler(w http.ResponseWriter, r *http.Request) {
	// Access the model value from the path
	vars := mux.Vars(r)
	konfigname := vars["konfig"]

	config := types.NewModelRequestWithDefaults()
	err := json.NewDecoder(r.Body).Decode(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config.ModelName = konfigname
	err = config.Save(s.ModelPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Konfig saved"))
}

func (s *Server) DeleteKonfigHandler(w http.ResponseWriter, r *http.Request) {
	// Access the model value from the path
	vars := mux.Vars(r)
	konfigname := vars["konfig"]

	config := types.NewModelRequestWithDefaults()
	err := config.Load(s.ModelPath, konfigname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	modelname := config.ModelName
	if _, ok := s.LoadedModels[modelname]; ok {
		http.Error(w, "Model is loaded", http.StatusBadRequest)
		return
	}

	err = config.Delete(s.ModelPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Konfig deleted"))
}

func (s *Server) completionProxy(w http.ResponseWriter, r *http.Request) {
	// Access the model value from the path
	vars := mux.Vars(r)
	modelname := vars["model"]

	//check if model is loaded
	if _, ok := s.LoadedModels[modelname]; !ok {
		http.Error(w, "Model not loaded", http.StatusBadRequest)
		return
	}

	s.genericProxy(w, r, "/completion", s.LoadedModels[modelname].Config.Port)
}

func (s *Server) infillProxy(w http.ResponseWriter, r *http.Request) {
	// Access the model value from the path
	vars := mux.Vars(r)
	modelname := vars["model"]

	//check if model is loaded
	if _, ok := s.LoadedModels[modelname]; !ok {
		http.Error(w, "Model not loaded", http.StatusBadRequest)
		return
	}
	s.genericProxy(w, r, "/infill", s.LoadedModels[modelname].Config.Port)

}

func (s *Server) genericProxy(w http.ResponseWriter, r *http.Request, path string, port int) {
	//proxy request to model
	newRequest := &http.Request{
		URL: &url.URL{
			Scheme: r.URL.Scheme,
			Host:   fmt.Sprintf("localhost:%d", port),
			Path:   path,
		},
		Method: r.Method,
		Header: r.Header,
		Body:   r.Body,
	}
	// Send the proxy request
	resp, err := http.DefaultClient.Do(newRequest)
	if err != nil {
		http.Error(w, "Failed to proxy request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy headers and status code
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)

	// Copy body
	io.Copy(w, resp.Body)
}

var logger = logrus.New()

func (s *Server) AddRoutes() {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/{model}/completion", s.completionProxy).Methods("POST")
	r.HandleFunc("/api/v1/{model}/infill", s.infillProxy).Methods("POST")

	r.HandleFunc("/api/v1/{model}/load", s.loadModelHandler).Methods("POST")
	r.HandleFunc("/api/v1/{model}/unload", s.unloadModelHandler).Methods("GET")
	r.HandleFunc("/api/v1/{model}/savetofile", s.saveModelHandler).Methods("GET")

	r.HandleFunc("/api/v1/models/loaded", s.getLoadedModelsHandler).Methods("GET")
	r.HandleFunc("/api/v1/models/available", s.getAvailableModelsHandler).Methods("GET")
	r.HandleFunc("/api/v1/models/download/{path}", s.downloadModelHandler).Methods("GET")
	r.HandleFunc("/api/v1/konfigs/available", s.getAvailableKonfigsHandler).Methods("GET")
	r.HandleFunc("/api/v1/{konfig}", s.GetKonfigHandler).Methods("GET")
	r.HandleFunc("/api/v1/{konfig}", s.SaveKonfigHandler).Methods("POST")
	r.HandleFunc("/api/v1/{konfig}", s.DeleteKonfigHandler).Methods("DELETE")
	r.HandleFunc("/api/v1/{konfig}/load", s.loadModelFromFileHandler).Methods("GET")
	s.Router = r
}

func GetServer(ModelPath, PathToLLama, Addr string) *Server {
	s := &Server{
		Router:       mux.NewRouter(),
		ModelPath:    ModelPath,
		PathToLLama:  PathToLLama,
		LoadedModels: map[string]*Runner{},
		usedPorts:    map[int]bool{8080: true},
	}
	s.AddRoutes()
	s.Server = &http.Server{
		Addr:    Addr,
		Handler: s.Router,
	}
	return s
}

func init() {
	// Initialize logger
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetLevel(logrus.InfoLevel)
}
