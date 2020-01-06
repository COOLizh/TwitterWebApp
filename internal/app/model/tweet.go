package model

import "time"

//Tweet ...
type Tweet struct {
	ID       uint      `json:"id"`
	AuthorID uint      `json:"author"`
	Message  string    `json:"message"`
	Date     time.Time `json:"date"`
}
