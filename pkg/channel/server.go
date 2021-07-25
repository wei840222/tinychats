package channel

import (
	"context"

	"github.com/wei840222/tinychats/proto"
	"gorm.io/gorm"
)

func NewServer(db *gorm.DB) (proto.ChannelServer, error) {
	if err := db.AutoMigrate(&Channel{}); err != nil {
		return nil, err
	}
	return &server{
		db: db,
	}, nil
}

type server struct {
	proto.UnimplementedChannelServer
	db *gorm.DB
}

func (s server) CreateChannel(ctx context.Context, req *proto.CreateChannelRequest) (*proto.CreateChannelResponse, error) {
	newChannel := Channel{
		Name: req.GetName(),
	}
	if err := s.db.WithContext(ctx).
		Create(&newChannel).Error; err != nil {
		return nil, err
	}
	return &proto.CreateChannelResponse{
		Channel: &proto.ChannelResponse{
			Id:   int64(newChannel.ID),
			Name: newChannel.Name,
		},
	}, nil
}

func (s server) ListChannels(ctx context.Context, req *proto.ListChannelsRequest) (*proto.ListChannelsResponse, error) {
	var channels []*Channel
	if err := s.db.WithContext(ctx).
		Find(&channels).Error; err != nil {
		return nil, err
	}
	var res proto.ListChannelsResponse
	for _, channel := range channels {
		res.Channels = append(res.Channels, &proto.ChannelResponse{
			Id:   int64(channel.ID),
			Name: channel.Name,
		})
	}
	return &res, nil
}
