package entities

type UserEntity struct {
	Id       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}
