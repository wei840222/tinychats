package graph

import (
	"sync"

	"github.com/wei840222/tinychats/graph/model"
	"github.com/wei840222/tinychats/proto"
	"go.uber.org/zap"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	sync.Mutex

	Logger        *zap.Logger
	UserClient    proto.UserClient
	MessageClient proto.MessageClient

	MessageCreatedChans map[string]chan *model.Message
}
