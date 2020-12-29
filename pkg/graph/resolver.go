package graph

import (
	"context"

	"github.com/sebastianvera/ghreviews"
	"github.com/sebastianvera/ghreviews/pkg/graph/generated"
	"github.com/sirupsen/logrus"
)

type Resolver struct {
	reviewService ghreviews.ReviewService
	logger        *logrus.Logger
	hub           *hub
}

func NewResolver(logger *logrus.Logger, reviewService ghreviews.ReviewService) *Resolver {
	return &Resolver{
		reviewService: reviewService,
		logger:        logger,
		hub:           NewHub(),
	}
}

func (r *mutationResolver) CreateReview(ctx context.Context, reviewInput generated.CreateReviewInput) (*ghreviews.Review, error) {
	review, err := r.reviewService.CreateReview(
		reviewInput.Username,
		reviewInput.AvatarURL,
		reviewInput.Content,
	)
	if err != nil {
		return nil, err
	}

	total, err := r.reviewService.CountReviews()
	if err != nil {
		return nil, err
	}

	go func() {
		msg := &generated.NewReviewEvent{Total: total, NewReview: review}
		r.hub.BroadcastNewReviewEvent(msg)
	}()

	return review, nil
}

func (r *queryResolver) GetMostRecentReviews(ctx context.Context, limit *int) ([]*ghreviews.Review, error) {
	l := 10
	if limit != nil {
		l = *limit
	}
	return r.reviewService.GetMostRecentReviews(l)
}

func (r *queryResolver) GetReviewsByUsername(ctx context.Context, username string) ([]*ghreviews.Review, error) {
	return r.reviewService.GetMostRecentReviewsByUsername(username)
}

func (r *queryResolver) CountReviewsByUsername(ctx context.Context, username string) (int, error) {
	return r.reviewService.CountReviewsByUsername(username)
}

func (r *subscriptionResolver) Feed(ctx context.Context) (<-chan *generated.NewReviewEvent, error) {
	cr := make(chan *generated.NewReviewEvent, 1)
	r.hub.Add(cr)

	r.logger.Debugf("[connected]: %d users\n", r.hub.Size())
	go func() {
		<-ctx.Done()
		r.hub.Remove(cr)
		r.logger.Debugf("[disconnected]: %d users\n", r.hub.Size())
	}()

	return cr, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Subscription returns MutationResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
