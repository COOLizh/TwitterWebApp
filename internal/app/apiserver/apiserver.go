package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/COOLizh/TwitterWebApp/internal/app/model"
	"github.com/COOLizh/TwitterWebApp/pkg/repository"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
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
	s.router.HandleFunc("/login", s.handleLogin()).Methods("POST")
	s.router.HandleFunc("/subscribe", s.handleSubscribe()).Methods("POST")
	s.router.HandleFunc("/tweets", s.handlePostTweet()).Methods("POST")
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
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			s.logger.Error(err)
			return
		}
		user, err = s.rep.Save(user)
		if err != nil {
			s.logger.Error(err)
			return
		}
		s.logger.Info("added user " + user.String())
		json.NewEncoder(w).Encode(user)
	}
}

/*
	handleLogin : checks the entered data by the user
		if the check is successful, it creates a jwt and stores it in cookies
*/
func (s *APIserver) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var usr model.User
		err := json.NewDecoder(r.Body).Decode(&usr)
		if err != nil {
			s.logger.Error(err)
			return
		}
		user, err := s.rep.FindByEmailAndPassword(usr)
		if err != nil {
			s.logger.Error(err)
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": user.Email,
			"id":    user.ID,
		})
		tokenString, err := token.SignedString([]byte(s.config.JwtSecret))
		if err != nil {
			s.logger.Error(err)
			return
		}
		s.logger.Info("jwt for " + user.Email + " was succesfully created")
		expire := time.Now().AddDate(0, 0, 1)
		cookie := http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expire,
		}
		http.SetCookie(w, &cookie)
		s.logger.Info("jwt stored in cookies")
		json.NewEncoder(w).Encode(model.JwtToken{Token: tokenString})

	}
}

/*
	isLoggedIn : checks if the user is authorized using a cookie jwt
*/
func (s *APIserver) isLoggedIn(w http.ResponseWriter, r *http.Request) (model.User, error) {
	c, err := r.Cookie("token")
	if err != nil {
		return model.User{}, err
	}
	token, _ := jwt.Parse(c.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return []byte(s.config.JwtSecret), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user model.User
		mapstructure.Decode(claims, &user)
		return user, nil
	}
	return model.User{}, fmt.Errorf("Invalid authorization token")
}

func (s *APIserver) handleSubscribe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		user, err := s.isLoggedIn(w, r)
		if err != nil {
			s.logger.Error(err)
		} else {
			s.logger.Info("user " + user.Email + " has been logged in. verification was performed using jwt.")
			json.NewEncoder(w).Encode(user)
		}
	}
}

// handlePostTweet find user in db and update his slice of tweets
func (s *APIserver) handlePostTweet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		user, err := s.isLoggedIn(w, r)
		if err != nil {
			s.logger.Error(err)
			return
		}
		var tweet model.Tweet
		err = json.NewDecoder(r.Body).Decode(&tweet)
		if err != nil {
			s.logger.Error(err)
			return
		}
		tweet, err = s.rep.UpdateUserTweets(user, tweet)
		if err != nil {
			s.logger.Error(err)
			return
		}
		s.logger.Info("tweet for user " + user.Email + " posted")
		json.NewEncoder(w).Encode(tweet)
	}
}
