package user

type User struct {
	Id        uint64
	Telephone string
	Password  string
	RealName  string
	Nickname  string
	Avatar    string
	Birthday  string
	Gender    uint8
	Status    uint8
}
