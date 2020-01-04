package model

//User ...
type User struct {
	ID          uint    `json:"id"`
	UserName    string  `json:"username"`
	Email       string  `json:"email"`
	PaswordHash string  `json:"password"`
	Following   []uint  `json:"subscriptions"`
	UserTweets  *Tweets `json:"tweets"`
	TweetsFeed  []uint  `json:"feed"`
}
