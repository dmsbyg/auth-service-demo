package auth

type User struct {
	ID       string
	Email    string
	Password []byte
}
