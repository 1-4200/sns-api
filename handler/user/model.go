package user

import "time"

type Response struct {
	Hits int         `json:"hits"`
	Res  interface{} `json:"res"`
}

type SearchForm struct {
	Name        string    `json:"name" form:"name" binding:"required_without=Description"`
	Description string    `json:"description" form:"description" binding:"required_without=Name"`
	Language    string    `json:"language" form:"language" binding:"omitempty"`
	FollowerMin int       `json:"follower_min" form:"follower_min" binding:"omitempty"`
	FollowerMax int       `json:"follower_max" form:"follower_max" binding:"omitempty,gtefield=FollowerMin"`
	StatusMin   int       `json:"status_min" form:"status_min" binding:"omitempty"`
	StatusMax   int       `json:"status_max" form:"status_max" binding:"omitempty,gtefield=StatusMin"`
	FavoriteMin int       `json:"favorite_min" form:"favorite_min" binding:"omitempty"`
	FavoriteMax int       `json:"favorite_max" form:"favorite_max" binding:"omitempty,gtefield=FavoriteMin"`
	FollowMin   int       `json:"follow_min" form:"follow_min" binding:"omitempty"`
	FollowMax   int       `json:"follow_max" form:"follow_max" binding:"omitempty,gtefield=FollowMin"`
	ListMin     int       `json:"list_min" form:"list_min" binding:"omitempty"`
	ListMax     int       `json:"list_max" form:"list_max" binding:"omitempty,gtefield=ListMin"`
	SrScoreMin  float64   `json:"sr_score_min" form:"sr_score_min" binding:"omitempty"`
	SrScoreMax  float64   `json:"sr_score_max" form:"sr_score_max" binding:"omitempty,gtefield=SrScoreMin"`
	StartDate   time.Time `json:"start_date" form:"start_date" binding:"required" time_format:"2006-01-02 15:04"`
	EndDate     time.Time `json:"end_date" form:"end_date" binding:"required,gtefield=StartDate" time_format:"2006-01-02 15:04"`
	OrderBy     string    `json:"order_by" form:"order_by" binding:"omitempty,oneof=followers_count friends_count listed_count favourites_count statuses_count"`
	Count       int       `json:"count" form:"count" binding:"omitempty,min=1,max=10000"`
}

type IDForm struct {
	UserID    uint64    `json:"user_id" form:"user_id" binding:"required"`
	StartDate time.Time `json:"start_date" form:"start_date" binding:"required" time_format:"2006-01-02 15:04"`
	EndDate   time.Time `json:"end_date" form:"end_date" binding:"required,gtefield=StartDate" time_format:"2006-01-02 15:04"`
}

type IDsForm struct {
	UserIDs   []uint64  `json:"user_ids" form:"user_ids" binding:"required,max=10000"`
	StartDate time.Time `json:"start_date" form:"start_date" binding:"required" time_format:"2006-01-02 15:04"`
	EndDate   time.Time `json:"end_date" form:"end_date" binding:"required,gtefield=StartDate" time_format:"2006-01-02 15:04"`
}
