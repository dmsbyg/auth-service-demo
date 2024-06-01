package auth

import "time"

type User struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Password  []byte    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}
