package server

import (
	"context"

	"gitlab.com/promptech1/infuser-author/handler"
	grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"
	"gitlab.com/promptech1/infuser-author/model"
)

type userServer struct {
	handler *handler.UserHandler
}

func newUserServer(handler *handler.UserHandler) grpc_author.UserServiceServer {
	return &userServer{handler: handler}
}

func (s userServer) Signup(ctx context.Context, req *grpc_author.UserReq) (*grpc_author.UserRes, error) {
	if req.Password != req.PasswordConfirmation {
		return &grpc_author.UserRes{
			Code: grpc_author.UserRes_PASSWORD_NOT_MATCHED,
		}, nil
	}

	if has, err := model.CheckLoginId(s.handler.Ctx.Orm, req.LoginId); err != nil {
		return nil, err
	} else if has {
		return &grpc_author.UserRes{
			Code: grpc_author.UserRes_DUPLICATE_LOGIN_ID,
		}, nil
	}

	if has, err := model.CheckEmail(s.handler.Ctx.Orm, req.Email); err != nil {
		return nil, err
	} else if has {
		return &grpc_author.UserRes{
			Code: grpc_author.UserRes_DUPLICATE_EMAIL,
		}, nil
	}

	user := &model.User{
		LoginId:    req.LoginId,
		Email:      req.Email,
		Name:       req.Name,
		LoginCount: 0,
	}

	enc, err := model.EncryptPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user.Password = enc

	if _, err := s.handler.Ctx.Orm.Insert(user); err != nil {
		return nil, err
	}

	return &grpc_author.UserRes{
		Code:    grpc_author.UserRes_VALID,
		Id:      uint32(user.Id),
		LoginId: user.LoginId,
		Email:   user.Email,
		Name:    user.Name,
	}, nil
}
