package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strconv"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/google/uuid"
	"github.com/wei840222/tinychats/graph/generated"
	"github.com/wei840222/tinychats/graph/model"
	"github.com/wei840222/tinychats/pkg"
	"github.com/wei840222/tinychats/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *messageResolver) User(ctx context.Context, obj *model.Message) (*model.User, error) {
	res, err := r.dataloadersFormContext(ctx).userLoader.Load(obj.User.ID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *mutationResolver) CreateMessage(ctx context.Context, input model.NewMessage) (*model.Message, error) {
	user, err := pkg.GetLINELoginUserFormContext(ctx)
	if err != nil {
		return nil, err
	}
	res, err := r.messageClient.CreateMessage(ctx, &proto.CreateMessageRequest{
		UserId: user.UserID,
		Text:   input.Text,
	})
	if err != nil {
		return nil, err
	}
	newMessage := model.Message{
		ID:        strconv.Itoa(int(res.GetMessage().GetId())),
		Text:      res.GetMessage().GetText(),
		CreatedAt: res.GetMessage().GetCreatedAt().AsTime().Format(time.RFC3339),
		User: &model.User{
			ID: res.GetMessage().GetUserId(),
		},
	}

	for _, messageCreatedChan := range r.messageCreatedChans {
		messageCreatedChan <- &newMessage
	}
	return &newMessage, nil
}

func (r *queryResolver) CurrentUser(ctx context.Context) (*model.User, error) {
	user, err := pkg.GetLINELoginUserFormContext(ctx)
	if err != nil {
		return nil, err
	}
	res, err := r.userClient.GetUser(ctx, &proto.GetUserRequest{Id: user.UserID})
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return nil, err
		}
		if _, err := r.userClient.CreateUser(ctx, &proto.CreateUserRequest{
			Id:        user.UserID,
			Name:      user.DisplayName,
			AvatarUrl: user.PictureURL,
		}); err != nil {
			return nil, err
		}
	} else {
		if _, err := r.userClient.UpdateUser(ctx, &proto.UpdateUserRequest{
			Id:        res.GetUser().GetId(),
			Name:      user.DisplayName,
			AvatarUrl: user.PictureURL,
		}); err != nil {
			return nil, err
		}
	}
	return &model.User{
		ID:        user.UserID,
		Name:      user.DisplayName,
		AvatarURL: pointer.ToStringOrNil(user.PictureURL),
	}, nil
}

func (r *queryResolver) Messages(ctx context.Context) ([]*model.Message, error) {
	res, err := r.messageClient.ListMessages(ctx, &proto.ListMessagesRequest{})
	if err != nil {
		return nil, err
	}
	var messages []*model.Message
	for _, messageResponse := range res.GetMessages() {
		messages = append(messages, &model.Message{
			ID:        strconv.Itoa(int(messageResponse.GetId())),
			Text:      messageResponse.GetText(),
			CreatedAt: messageResponse.GetCreatedAt().AsTime().Format(time.RFC3339),
			User: &model.User{
				ID: messageResponse.GetUserId(),
			},
		})
	}
	return messages, nil
}

func (r *subscriptionResolver) MessageCreated(ctx context.Context) (<-chan *model.Message, error) {
	uuid := uuid.New().String()
	r.Lock()
	r.messageCreatedChans[uuid] = make(chan *model.Message, 1)
	r.Unlock()

	go func() {
		<-ctx.Done()
		r.Lock()
		delete(r.messageCreatedChans, uuid)
		r.Unlock()
	}()

	return r.messageCreatedChans[uuid], nil
}

// Message returns generated.MessageResolver implementation.
func (r *Resolver) Message() generated.MessageResolver { return &messageResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type messageResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
