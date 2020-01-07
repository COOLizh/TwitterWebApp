package repository

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/COOLizh/TwitterWebApp/internal/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// UsersRepositoryMongo ...
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

// FindByEmailAndPassword : checks if the given user exists with such email and passwrod and returns this user
func (r *UsersRepositoryMongo) FindByEmailAndPassword(user model.User) (model.User, error) {
	if strings.Contains(user.Email, " ") || strings.Contains(user.PasswordHash, " ") {
		return model.User{}, fmt.Errorf("input must not contain a space")
	}
	if user.Email == "" || user.PasswordHash == "" {
		return model.User{}, fmt.Errorf("empty input")
	}
	collection := r.db.Collection("Users")
	var foundResult bson.M
	collection.FindOne(context.TODO(), bson.D{{Key: "email", Value: user.Email}}).Decode(&foundResult)
	if foundResult["email"] == user.Email && checkPasswordHash(user.PasswordHash, foundResult["passwordhash"].(string)) {
		return model.User{
			ID:           uint(foundResult["id"].(int64)),
			UserName:     foundResult["username"].(string),
			Email:        foundResult["email"].(string),
			PasswordHash: foundResult["passwordhash"].(string),
		}, nil
	}
	return model.User{}, fmt.Errorf("no surch registered user. password or email entered incorrectly")
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

// UpdateUserTweets find user and add his tweet
func (r *UsersRepositoryMongo) UpdateUserTweets(user model.User, tweet model.Tweet) (model.Tweet, error) {
	collection := r.db.Collection("Users")
	var foundResult model.User
	collection.FindOne(context.TODO(), bson.D{{Key: "id", Value: user.ID}}).Decode(&foundResult)
	id := len(foundResult.UserTweets) + 1
	tweet.ID = uint(id)
	tweet.AuthorID = foundResult.ID
	tweet.Date = time.Now()

	// in user's feed he also see his own tweets
	foundResult.UserTweets = append(foundResult.UserTweets, tweet)
	foundResult.TweetsFeed = append(foundResult.TweetsFeed, tweet)

	// adding this tweet in tweetsFeed of subscribers
	for _, follower := range foundResult.Followers {
		var tmpUser model.User
		collection.FindOne(context.TODO(), bson.D{{Key: "username", Value: follower}}).Decode(&tmpUser)
		tmpUser.TweetsFeed = append(tmpUser.TweetsFeed, tweet)
		_, err := collection.UpdateOne(context.Background(), bson.D{{Key: "username", Value: follower}}, bson.M{"$set": tmpUser})
		if err != nil {
			return model.Tweet{}, err
		}
	}

	sort.Slice(foundResult.TweetsFeed, func(i, j int) bool { return foundResult.TweetsFeed[i].Date.After(foundResult.TweetsFeed[j].Date) })
	_, err := collection.UpdateOne(context.Background(), bson.D{{Key: "id", Value: foundResult.ID}}, bson.M{"$set": foundResult})
	if err != nil {
		return model.Tweet{}, err
	}
	return tweet, nil
}

// AddToFollowing : add userToSubscribe username to subscriptions slice of user and add userToSubscribe tweets to TweetsFeed of user
func (r *UsersRepositoryMongo) AddToFollowing(user, userToSubscribe model.User) error {
	collection := r.db.Collection("Users")
	var _user, _userToSubscribe model.User
	collection.FindOne(context.TODO(), bson.D{{Key: "id", Value: user.ID}}).Decode(&_user)
	collection.FindOne(context.TODO(), bson.D{{Key: "username", Value: userToSubscribe.UserName}}).Decode(&_userToSubscribe)
	if _userToSubscribe.UserName != userToSubscribe.UserName {
		return fmt.Errorf("there is no user with such username")
	}
	if _userToSubscribe.UserName == _user.UserName {
		return fmt.Errorf("can't add to followers yourself")
	}
	_user.Following = append(_user.Following, _userToSubscribe.UserName)
	_userToSubscribe.Followers = append(_userToSubscribe.Followers, _user.UserName)
	for _, userToSubscribeTweet := range _userToSubscribe.UserTweets {
		_user.TweetsFeed = append(_user.TweetsFeed, userToSubscribeTweet)
	}
	sort.Slice(_user.TweetsFeed, func(i, j int) bool { return _user.TweetsFeed[i].Date.After(_user.TweetsFeed[j].Date) })
	_, err := collection.UpdateOne(context.Background(), bson.D{{Key: "id", Value: _user.ID}}, bson.M{"$set": _user})
	if err != nil {
		return err
	}
	_, err = collection.UpdateOne(context.Background(), bson.D{{Key: "id", Value: _userToSubscribe.ID}}, bson.M{"$set": _userToSubscribe})
	if err != nil {
		return err
	}
	return nil
}

// GetTweetsFeed return all user tweets sorted by date
func (r *UsersRepositoryMongo) GetTweetsFeed(user model.User) ([]model.Tweet, error) {
	collection := r.db.Collection("Users")
	var _user model.User
	collection.FindOne(context.TODO(), bson.D{{Key: "id", Value: user.ID}}).Decode(&_user)
	sort.Slice(_user.TweetsFeed, func(i, j int) bool { return _user.TweetsFeed[i].Date.After(_user.TweetsFeed[j].Date) })
	_, err := collection.UpdateOne(context.Background(), bson.D{{Key: "id", Value: _user.ID}}, bson.M{"$set": _user})
	if err != nil {
		return []model.Tweet{}, err
	}
	return _user.TweetsFeed, nil
}
