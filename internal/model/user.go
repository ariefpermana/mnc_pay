package model

import "time"

type User struct {
	Id        int       `gorm:"primaryKey;autoIncrement;column:id;type:integer"`
	Username  string    `gorm:"column:username;type:varchar(100)"`
	Password  string    `gorm:"column:password;type:varchar(255)"`
	Token     string    `gorm:"column:token;type:varchar(255)"`
	Role      int       `gorm:"column:role;type:integer"`
	IsLogin   int       `gorm:"column:is_login;type:integer"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (User) TableName() string {
	return "user"
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int    `json:"role"`
	IsLogin  int    `json:"is_login"`
}

type UserLoginResp struct {
	Id        int
	Username  string
	AuthToken string
}

type UserCreateResp struct {
	Id       int
	Username string
}

type UserLogoutResp struct {
	Id       int
	Username string
	IsLogin  int
}
