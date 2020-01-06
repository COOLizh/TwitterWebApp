package model

import "fmt"

//User model
type User struct {
	ID           uint    `json:"id"`
	UserName     string  `json:"username"`
	Email        string  `json:"email"`
	PasswordHash string  `json:"password"`
	Following    []uint  `json:"subscriptions"`
	UserTweets   *Tweets `json:"tweets"`
	TweetsFeed   []uint  `json:"feed"`
}

// JwtToken ...
type JwtToken struct {
	Token string `json:"token"`
}

func (u User) String() string {
	return fmt.Sprintf("User:\n\tID - %d\n\tUserName - %q\n\tEmail - %q\n\tPasswordHash - %q\n", u.ID, u.UserName, u.Email, u.PasswordHash)
}
