package repository

import (
	"context"
	"fmt"
	"strings"
	//"time"

	"github.com/COOLizh/TwitterWebApp/internal/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// UsersRepositoryMongo gets a pointer to the database
type UsersRepositoryMongo struct {
	db *mongo.Database
}

// NewUsersRepositoryMongo initialize UsersRepositoryMongo
func NewUsersRepositoryMongo(database *mongo.Database) *UsersRepositoryMongo {
	return &UsersRepositoryMongo{db: database}
}

// Save : adding user to database if such user doesn't exists
func (r *UsersRepositoryMongo) Save(user model.User) (model.User, error) {
	if strings.Contains(user.UserName, " ") || strings.Contains(user.Email, " ") || strings.Contains(user.PasswordHash, " ") {
		return model.User{}, fmt.Errorf("input must not contain a space")
	}
	if user.UserName == "" || user.Email == "" || user.PasswordHash == "" {
		return model.User{}, fmt.Errorf("empty input")
	}
	collection := r.db.Collection("Users")
	var foundResult bson.M
	collection.FindOne(context.TODO(), bson.D{{Key: "email", Value: user.Email}}).Decode(&foundResult)
	if foundResult["email"] == user.Email {
		return model.User{}, fmt.Errorf("user with email %q already exists", user.Email)
	}
	collection.FindOne(context.TODO(), bson.D{{Key: "username", Value: user.UserName}}).Decode(&foundResult)
	if foundResult["username"] == user.UserName {
		return model.User{}, fmt.Errorf("user with username %q already exists", user.UserName)
	}
	id, err := collection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		return model.User{}, err
	}
	id++
	user.ID = uint(id)
	user.PasswordHash, err = hashPassword(user.PasswordHash)
	if err != nil {
		return model.User{}, err
	}
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

//hashPassword get hash of password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//checkPasswordHash compare password with hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
