package repo

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"gitlab.com/promptech1/infuser-author/model"
)

type UserRepository interface {
	Create(u *model.User) (*model.User, error)
	FindOne(id uint) *model.User
	FindOneByEmail(email string) *model.User
}

type userRepositoryDB struct {
	DB *gorm.DB
}

func (r *userRepositoryDB) FindOne(id uint) *model.User {
	u := &model.User{}
	if r.DB.First(&u, id).RecordNotFound() {
		return nil
	}
	return u
}

func (r *userRepositoryDB) FindOneByEmail(email string) *model.User {
	u := &model.User{}
	if r.DB.Where("email = ?", email).First(&u).RecordNotFound() {
		return nil
	}
	return u
}

func (r *userRepositoryDB) Create(u *model.User) (*model.User, error) {
	u.EncPassword()
	if dbc := r.DB.Create(u); dbc.Error != nil {
		glog.Errorf("Error in user repository create: %v", dbc.Error)
		return nil, fmt.Errorf("Error in user repository create: %v", dbc.Error)
	}

	return u, nil
}

func NewUserRepository(db *gorm.DB)  UserRepository {
	return &userRepositoryDB{
		DB: db,
	}
}