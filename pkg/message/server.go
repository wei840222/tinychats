package message

import (
	"context"

	"github.com/wei840222/tinychats/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func NewServer(db *gorm.DB) (proto.MessageServer, error) {
	if err := db.AutoMigrate(&Message{}); err != nil {
		return nil, err
	}
	return &server{
		db: db,
	}, nil
}

type server struct {
	proto.UnimplementedMessageServer
	db *gorm.DB
}

func (s server) CreateMessage(ctx context.Context, req *proto.CreateMessageRequest) (*proto.CreateMessageResponse, error) {
	newMessage := Message{
		UserID: req.GetUserId(),
		Text:   req.GetText(),
	}
	if err := s.db.WithContext(ctx).
		Create(&newMessage).Error; err != nil {
		return nil, err
	}
	return &proto.CreateMessageResponse{
		Id:        int64(newMessage.ID),
		CreatedAt: timestamppb.New(newMessage.CreatedAt),
	}, nil
}

func (s server) ListMessages(ctx context.Context, _ *proto.ListMessagesRequest) (*proto.ListMessagesResponse, error) {
	var messages []*Message
	if err := s.db.WithContext(ctx).Find(&messages).Error; err != nil {
		return nil, err
	}
	var res proto.ListMessagesResponse
	for _, message := range messages {
		res.Messages = append(res.Messages, &proto.MessageResponse{
			Id:        int64(message.ID),
			UserId:    message.UserID,
			Text:      message.Text,
			CreatedAt: timestamppb.New(message.CreatedAt),
		})
	}
	return &res, nil
}
