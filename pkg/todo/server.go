package todo

import (
	"github.com/wei840222/todo/proto"

	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func NewServer(db *gorm.DB) (proto.TodoServer, error) {
	if err := db.AutoMigrate(&Todo{}); err != nil {
		return nil, err
	}
	return &server{
		db: db,
	}, nil
}

type server struct {
	proto.UnimplementedTodoServer
	db *gorm.DB
}

func (s server) ListTodos(_ context.Context, _ *emptypb.Empty) (*proto.ListTodoResponse, error) {
	var todos []*Todo
	if err := s.db.Find(&todos).Error; err != nil {
		return nil, err
	}

	var res proto.ListTodoResponse
	for _, t := range todos {
		res.Todos = append(res.Todos, &proto.TodoResponse{
			Id:   t.UUID.String(),
			Text: t.Text,
			Done: t.Done,
		})
	}

	return &res, nil
}

func (s server) CreateTodo(_ context.Context, req *proto.CreateTodoRequest) (*proto.CreateTodoResponse, error) {
	id := uuid.New()

	if err := s.db.Create(&Todo{UUID: id, Text: req.GetText()}).Error; err != nil {
		return nil, err
	}

	return &proto.CreateTodoResponse{
		Id: id.String(),
	}, nil
}
