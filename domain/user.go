package domain

import "time"

type User struct {
	UserID           string  `json:"user_id"`
	UserScreenName   string  `json:"user_screen_name"`
	UserName         string  `json:"user_name"`
	UserDescription  string  `json:"user_description"`
	UserImageProfile string  `json:"user_image_profile"`
	Verified         bool    `json:"verified"`
	FollowerCount    float64 `json:"follower_count"`
	StatusCount      float64 `json:"status_count"`
	FavoriteCount    float64 `json:"favorite_count"`
	FollowCount      float64 `json:"follow_count"`
	ListCount        float64 `json:"list_count"`
	SrScore          float64 `json:"sr_score"`
	CreatedAt        string  `json:"created_at"`
}

type UserRepository interface {
	Search(name, description, language string, followerMin, followerMax, statusMin, statusMax, favoriteMin, favoriteMax, followMin, followMax, listMin, listMax int, srScoreMin, srScoreMax float64, startDate, endDate time.Time, count int, orderBy string) ([]*User, int, error)
	GetById(userID uint64, startDate, endDate time.Time) (*User, int, error)
	GetByIds(userIDs []uint64, startDate, endDate time.Time) ([]*User, int, error)
}
