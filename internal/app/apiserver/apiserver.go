package apiserver

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/COOLizh/TwitterWebApp/internal/app/model"
	"github.com/COOLizh/TwitterWebApp/pkg/repository"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// APIserver ...
type APIserver struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	rep    *repository.UsersRepositoryMongo
}

// New ...
func New(config *Config) *APIserver {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.DbConnectionString))
	return &APIserver{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
		rep:    repository.NewUsersRepositoryMongo(client.Database(config.DbName)),
	}
}

// Start ...
func (s *APIserver) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.configureRouter()
	s.logger.Info("router configured")

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
		user, err := s.rep.Save(user)
		if err != nil {
			s.logger.Info(err)
		} else {
			s.logger.Info("added user " + user.String())
		}
		json.NewEncoder(w).Encode(user)
	}
}

func (s *APIserver) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Login")
	}
}
