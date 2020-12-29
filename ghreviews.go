package ghreviews

// Review is the object that is sent to the client
type Review struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatarUrl"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"createdAt"`
}

// ReviewService represents the interface of the data storage
type ReviewService interface {
	CreateReview(username, avatarURL, content string) (*Review, error)
	GetMostRecentReviews(int) ([]*Review, error)
	GetMostRecentReviewsByUsername(username string) ([]*Review, error)
	CountReviews() (int, error)
	CountReviewsByUsername(username string) (int, error)
}
