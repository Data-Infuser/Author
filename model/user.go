package model

import (
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"
	"golang.org/x/crypto/bcrypt"
)

// User : 사용자 관리 모델
type User struct {
	gorm.Model
	Email    string `gorm:"type:varchar(100);unique_index"`
	Name     string `gorm:"type:varchar(20);"`
	Password string `gorm:"size:191;"`
}

// EncPassword : 사용자 비밀번호 암호화 처리
func (u *User) EncPassword() {
	bytes := []byte(u.Password)
	hash, err := bcrypt.GenerateFromPassword(bytes, bcrypt.MinCost)
	if err != nil {
		glog.Error(err)
	}

	glog.Infof("Hashed Password: %s", string(hash))
	u.Password = string(hash)
}

// GetgRPCModel : gRPC Message 변환
func (u *User) GetgRPCModel() *grpc_author.UserRes {
	return &grpc_author.UserRes{
		Id:    uint32(u.ID),
		Email: u.Email,
		Name:  u.Name,
	}
}
