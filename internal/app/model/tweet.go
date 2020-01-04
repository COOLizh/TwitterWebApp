package model

import "time"

//Tweet ...
type Tweet struct {
	ID      uint      `json:"id"`
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
}

//Tweets ...
type Tweets []Tweet
