package models

type User struct {
	Model
	FirstName string `gorm:"column:first_name" json:"first_name"`
	LastName  string `gorm:"column:last_name" json:"last_name"`
	Address   string `gorm:"column:address" json:"address"`
	Email     string `gorm:"column:email" json:"email"`
}

var user User

// NewUser creates a new User
func NewUser() *User {
	return &User{}
}

func (u *User) Get(id uint64) (*User, error) {
	user.ID = id
	// err := env.Config.DB.Where("id = ?", id). // Do the query
	// 	First(&user). // Make it scalar
	// 	Error // retrieve error or null
	// return &user, err
	return &user, nil
}
