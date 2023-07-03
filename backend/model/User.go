package model

type User struct {
	ID        int64  `json:"id" form:"ID" gorm:"primaryKey;autoIncrement;column:ID"`
	Username  string `json:"username" gorm:"primaryKey;column:username"`
	Password  string `json:"password" gorm:"column:password"`
	Password2 string `json:"password2" gorm:"-"`
	Nickname  string `json:"nickname" gorm:"column:nickName"`
	Email     string `json:"email" gorm:"column:email"`
}
