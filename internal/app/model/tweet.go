package model

import "time"

//Tweet ...
type Tweet struct {
	ID      uint
	Message string
	Date    time.Time
}
