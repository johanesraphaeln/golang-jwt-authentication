package models

type User struct {
	// ID       int    `json:"id" gorm:"primary_key"` //mysql driver
	ID       int    `json:"id" gorm:"primaryKey"` //postgres driver
	Username string `json:"username"`
	Password string `json:"-"`
}
