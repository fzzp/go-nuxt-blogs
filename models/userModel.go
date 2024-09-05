package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID        int64   `json:"id"`
	Email     string  `json:"email"`
	Password  string  `json:"-"`
	Username  string  `json:"username"`
	Avatar    string  `json:"avatar"`
	Role      int64   `json:"role"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	DeletedAt *string `json:"-"`
}

// Hash 加密密码，并修改自身 password
func (user *User) Hash() (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return "", err
	}
	user.Password = string(hashedPassword)
	return user.Password, nil
}

// Matches 匹配密码是否正确
func (user *User) Matches(plaintext string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plaintext))
}
