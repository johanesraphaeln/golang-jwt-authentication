package models

type User struct {
	// mysql driver
	// ID       int    `json:"id" gorm:"primary_key"`
	ID       int    `json:"id" gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"-"`
}
