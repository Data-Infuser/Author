package model

import (
	"net/http"
	"time"

	errors "gitlab.com/promptech1/infuser-author/error"
	"golang.org/x/crypto/bcrypt"
	"xorm.io/xorm"
)

type User struct {
	Id                  uint   `xorm:"pk autoincr"`
	GroupId             uint   `xorm:"index"`
	LoginId             string `xorm:"unique"`
	Password            string
	Email               string `xorm:"unique"`
	Name                string
	LoginCount          uint
	LastLoginAt         time.Time
	ResetPasswordToken  string
	ResetPasswordSentAt time.Time
	CreatedAt           time.Time  `xorm:"created"`
	UpdatedAt           time.Time  `xorm:"updated"`
	DeletedAt           *time.Time `xorm:"deleted index"`
}

func (User) TableName() string {
	return "user"
}

func (u *User) Find(orm *xorm.Engine) error {
	found, err := orm.Get(u)

	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}

	if !found {
		return errors.NewWithCode(http.StatusNotFound, "user not found")
	}

	return nil
}

func EncryptPassword(password string) (string, error) {
	enc, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(enc), nil
}

func ComparePasswords(hashedPwd string, plainPwd string) (bool, error) {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPwd))
	if err != nil {
		return false, err
	}

	return true, nil
}

func CheckLoginId(orm *xorm.Engine, loginId string) (bool, error) {
	return orm.Get(&User{LoginId: loginId})
}

func CheckEmail(orm *xorm.Engine, email string) (bool, error) {
	return orm.Get(&User{Email: email})
}
