package domain

import "time"

type Tweet struct {
	UserID         string            `json:"user_id"`
	UserScreenName string            `json:"user_screen_name"`
	UserName       string            `json:"user_name"`
	TweetID        string            `json:"tweet_id"`
	Text           string            `json:"text"`
	QuoteCount     float64           `json:"quote_count"`
	FavoriteCount  float64           `json:"favorite_count"`
	RetweetCount   float64           `json:"retweet_count"`
	ReplyCount     float64           `json:"reply_count"`
	CreatedAt      string            `json:"created_at"`
	NestedURL      []*TweetNestedURL `json:"nested_url"`
}

type TweetNestedURL struct {
	CanonicalURL string `json:"canonical_url"`
	Domain       string `json:"domain"`
}

type URL struct {
	URL         string `json:"canonical_url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TweetMedia struct {
	TweetID       string  `json:"tweet_id"`
	MediaType     float64 `json:"media_type"`
	FavoriteCount float64 `json:"favorite_count"`
	RetweetCount  float64 `json:"retweet_count"`
}

type Media struct {
	TweetID  string `json:"tweet_id"`
	MediaURL string `json:"media_url"`
}

type TweetTransition struct {
	UserID        uint64 `json:"user_id"`
	FollowerCount uint64 `json:"follower_count"`
	FriendCount   uint64 `json:"friend_count"`
	ListedCount   uint64 `json:"listed_count"`
	FavoriteCount uint64 `json:"favorite_count"`
	StatusCount   uint64 `json:"status_count"`
	CreatedAt     string `json:"created_at"`
}

type TweetRepository interface {
	Get() ([]*Tweet, error)
	GetByUser(userID uint64, startDate, endDate string, count int, orderBy string) ([]*Tweet, int, error)
	GetByUsers(userIDs []uint64, startDate, endDate time.Time, count int, orderBy string) ([]*Tweet, int, error)
	GetByDomain(userID uint64, startDate, endDate string, count int, orderBy string, domainName string) ([]*Tweet, int, []*URL, error)
	GetByMediaType(userID uint64, startDate, endDate string, count int, orderBy string, mediaType int) ([]*TweetMedia, int, []*Media, error)
}

type TransitionRepository interface {
	GetTransitionByUser(userID uint64, startDate, endDate string, count int) ([]*TweetTransition, error)
}
