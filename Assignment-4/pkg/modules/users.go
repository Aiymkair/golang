package modules

import "time"

type User struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Age       *int      `db:"age" json:"age,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
