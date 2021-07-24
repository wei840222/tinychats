package graph

import (
	"github.com/wei840222/tinychats/graph/model"
	"github.com/wei840222/tinychats/proto"
	"go.uber.org/zap"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Logger        *zap.Logger
	UserClient    proto.UserClient
	MessageClient proto.MessageClient

	MessageCreatedChan chan *model.Message
}
