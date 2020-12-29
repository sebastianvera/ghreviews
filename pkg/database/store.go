package database

import (
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sebastianvera/ghreviews"
)

var _ ghreviews.ReviewService = &Store{}

// Store represents a database interface
type Store struct {
	mu      sync.RWMutex
	reviews []*ghreviews.Review
}

// NewStore constructs an in memory store
func NewStore() *Store {
	return &Store{}
}

var id uint64 = 0

func getNextID() string {
	return strconv.FormatUint(atomic.AddUint64(&id, 1), 10)
}

// CreateReview inserts a review on the database
func (s *Store) CreateReview(username, avatarUrl, content string) (*ghreviews.Review, error) {
	r := &ghreviews.Review{
		ID:        getNextID(),
		Username:  username,
		AvatarURL: avatarUrl,
		Content:   content,
		CreatedAt: toMilliseconds(time.Now()),
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// add to the beginning
	s.reviews = append([]*ghreviews.Review{r}, s.reviews...)

	return r, nil
}

// GetMostRecentReviews returns N most recent reviews
func (s *Store) GetMostRecentReviews(limit int) ([]*ghreviews.Review, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	last := min(limit, len(s.reviews))
	rr := s.reviews[:last]
	return rr, nil
}

// GetMostRecentReviewsByUsername returns N most recent reviews for a given user
func (s *Store) GetMostRecentReviewsByUsername(username string) ([]*ghreviews.Review, error) {
	rr := []*ghreviews.Review{}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, r := range s.reviews {
		if r.Username == username {
			rr = append(rr, r)
		}
	}

	return rr, nil
}

// CountReviews returns total amount of reviews for a given user
func (s *Store) CountReviews() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.reviews), nil
}

// CountReviewsByUsername returns N most recent reviews for a given user
func (s *Store) CountReviewsByUsername(username string) (int, error) {
	count := 0

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, r := range s.reviews {
		if r.Username == username {
			count++
		}
	}

	return count, nil
}

func toMilliseconds(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
