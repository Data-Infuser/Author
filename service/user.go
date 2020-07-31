package service

import (
	grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"
	"gitlab.com/promptech1/infuser-author/model"
	repo "gitlab.com/promptech1/infuser-author/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService interface {
	Create(userReq *grpc_author.UserReq) (*grpc_author.UserRes, error)
}

type userService struct {
	userRepo repo.UserRepository
}

func NewUserService(repo repo.UserRepository) UserService {
	return &userService{userRepo: repo}
}

func (s userService) Create(userReq *grpc_author.UserReq) (*grpc_author.UserRes, error) {
	user := &model.User{
		Email: userReq.Email,
		Name: userReq.Name,
		Password: userReq.Password,
	}

	u := s.userRepo.FindOneByEmail(user.Email)
	if u != nil {
		return nil, status.Errorf(codes.AlreadyExists,
			"이미 사용중인 이메일 주소입니다.")
	}

	user, err := s.userRepo.Create(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return user.GetgRPCModel(), nil
}