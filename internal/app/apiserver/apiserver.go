package apiserver

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/COOLizh/TwitterWebApp/internal/app/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// APIserver ...
type APIserver struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
}

// New ...
func New(config *Config) *APIserver {
	return &APIserver{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// Start ...
func (s *APIserver) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.configureRouter()
	s.logger.Info("starting api server")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIserver) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *APIserver) configureRouter() {
	s.router.HandleFunc("/register", s.handleRegistration()).Methods("POST")
	s.router.HandleFunc("/login", s.handleLogin())
}

/*
	handleRegistration provides registration
	Request example :
	{
		"username":"dddd",
		"email":"aaa",
		"password":"bbb"
	}
*/
func (s *APIserver) handleRegistration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var user model.User
		_ = json.NewDecoder(r.Body).Decode(&user)
		user.ID = 1
		json.NewEncoder(w).Encode(user)
	}
}

func (s *APIserver) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Login")
	}
}
