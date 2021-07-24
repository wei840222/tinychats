package user

import (
	"context"
	"errors"

	"github.com/wei840222/tinychats/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func NewServer(db *gorm.DB) (proto.UserServer, error) {
	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}
	return &server{
		db: db,
	}, nil
}

type server struct {
	proto.UnimplementedUserServer
	db *gorm.DB
}

func (s server) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	var user User
	if err := s.db.WithContext(ctx).
		Where("user_id = ?", req.GetId()).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	return &proto.GetUserResponse{
		User: &proto.UserResponse{
			Id:        user.UserID,
			Name:      user.Name,
			AvatarUrl: user.AvatarURL,
		},
	}, nil
}

func (s server) MutiGetUsers(ctx context.Context, req *proto.MutiGetUsersRequest) (*proto.MutiGetUsersResponse, error) {
	var res proto.MutiGetUsersResponse
	if len(req.GetIds()) <= 0 {
		return &res, nil
	}
	var users []*User
	if err := s.db.WithContext(ctx).
		Where("user_id IN (?)", req.GetIds()).
		Find(&users).Error; err != nil {
		return nil, err
	}
	for _, user := range users {
		res.Users = append(res.Users, &proto.UserResponse{
			Id:        user.UserID,
			Name:      user.Name,
			AvatarUrl: user.AvatarURL,
		})
	}
	return &res, nil
}

func (s server) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*emptypb.Empty, error) {
	if err := s.db.WithContext(ctx).
		Create(&User{
			UserID:    req.GetId(),
			Name:      req.GetName(),
			AvatarURL: req.GetAvatarUrl(),
		}).Error; err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s server) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*emptypb.Empty, error) {
	if err := s.db.WithContext(ctx).
		Where("user_id = ?", req.GetId()).
		Updates(&User{
			Name:      req.GetName(),
			AvatarURL: req.GetAvatarUrl(),
		}).Error; err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
