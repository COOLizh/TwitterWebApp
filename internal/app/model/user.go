package model

//User ...
type User struct {
	ID          uint
	UserName    string
	Email       string
	PaswordHash string
	Followers   []uint
	Following   []uint
	TweetsFeed  []uint
}
