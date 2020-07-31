package model

import (
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Email string `gorm:"type:varchar(100);unique_index"`
	Name string `gorm:"type:varchar(20);"`
	Password string
}

func (u User) EncPassword() {
	bytes := []byte(u.Password)
	hash, err := bcrypt.GenerateFromPassword(bytes, bcrypt.MinCost)
	if err != nil {
		glog.Error(err)
	}

	u.Password = string(hash)
}

func (u User) GetgRPCModel() *grpc_author.UserRes {
	return &grpc_author.UserRes{
		Id: uint32(u.ID),
		Email: u.Email,
		Name: u.Name,
	}
}