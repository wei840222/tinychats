package todo

import (
	"github.com/wei840222/todo/proto"

	"context"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type Server struct {
	proto.UnimplementedTodoServer
	sync.RWMutex
	todos []*proto.TodoResponse
	DB    *gorm.DB
}

func (s *Server) ListTodos(_ context.Context, _ *emptypb.Empty) (*proto.ListTodoResponse, error) {
	s.RLock()
	defer s.RUnlock()
	return &proto.ListTodoResponse{
		Todos: s.todos,
	}, nil
}

func (s *Server) CreateTodo(_ context.Context, req *proto.CreateTodoRequest) (*proto.CreateTodoResponse, error) {
	s.Lock()
	defer s.Unlock()
	id := uuid.New()
	s.todos = append(s.todos, &proto.TodoResponse{
		Id:   id.String(),
		Text: req.GetText(),
	})
	return &proto.CreateTodoResponse{
		Id: id.String(),
	}, nil
}
