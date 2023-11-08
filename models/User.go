package models

type User struct {
	ID       uint   `gorm:"unique" json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Number   string `gorm:"unique" json:"number"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Created  string `json:"created"`
	Updated  string `json:"updated"`
	Deleted  bool   `json:"deleted"`
}

type CreatedUser struct {
	ID      int    `json:"ID"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Number  string `json:"number"`
	Role    string `json:"role"`
	Created string `json:"Created"`
	Updated string `json:"Updated"`
}

type ResetPassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
