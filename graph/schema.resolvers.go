package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/wei840222/todo/graph/generated"
	"github.com/wei840222/todo/graph/model"
	"github.com/wei840222/todo/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	res, err := r.TodoClient.CreateTodo(ctx, &proto.CreateTodoRequest{
		Text: input.Text,
	})
	if err != nil {
		return nil, err
	}
	return &model.Todo{
		ID:   res.GetId(),
		Text: input.Text,
		Done: false,
	}, nil
}

func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	res, err := r.TodoClient.ListTodos(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	var todos []*model.Todo
	for _, t := range res.GetTodos() {
		todos = append(todos, &model.Todo{
			ID:   t.GetId(),
			Text: t.GetText(),
			Done: t.GetDone(),
		})
	}
	return todos, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
