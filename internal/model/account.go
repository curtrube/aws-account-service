package model

import "time"

type Account struct {
	ID              string    `db:"id"`
	Name            string    `db:"name"`
	Email           string    `db:"email"`
	JoinedMethod    string    `db:"joined_method"`
	JoinedTimestamp time.Time `db:"joined_timestamp"`
	State           string    `db:"state"`
}
