package user

import "time"

type User struct {
	Id        uint64
	Telephone string
	Password  string
	RealName  string
	Nickname  string
	Avatar    string
	Birthday  time.Time
	Gender    uint8
	Status    uint8
}
