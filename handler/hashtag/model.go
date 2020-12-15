package hashtag

import "time"

type Response struct {
	Hits int         `json:"hits"`
	Res  interface{} `json:"res"`
}

type Form struct {
	Keyword         string    `json:"keyword" form:"keyword" binding:"required_without=Hashtag"`
	Hashtag         []string  `json:"hashtag" form:"hashtag" binding:"required_without=Keyword"`
	TweetType       []int     `json:"tweet_type" form:"tweet_type" binding:"omitempty"`
	RetweetMin      int       `json:"retweet_min" form:"retweet_min" binding:"omitempty"`
	RetweetMax      int       `json:"retweet_max" form:"retweet_max" binding:"omitempty,gtefield=RetweetMin"`
	QuoteMin        int       `json:"quote_min" form:"quote_min" binding:"omitempty"`
	QuoteMax        int       `json:"quote_max" form:"quote_max" binding:"omitempty,gtefield=QuoteMin"`
	FavoriteMin     int       `json:"favorite_min" form:"favorite_min" binding:"omitempty"`
	FavoriteMax     int       `json:"favorite_max" form:"favorite_max" binding:"omitempty,gtefield=FavoriteMin"`
	UserInclude     []string  `json:"user_include" form:"user_include" binding:"omitempty"`
	UserExclude     []string  `json:"user_exclude" form:"user_exclude" binding:"omitempty"`
	HashtagInclude  []string  `json:"hashtag_include" form:"hashtag_include" binding:"omitempty"`
	HashtagExclude  []string  `json:"hashtag_exclude" form:"hashtag_exclude" binding:"omitempty"`
	UserFollowerMin int       `json:"user_follower_min" form:"user_follower_min" binding:"omitempty"`
	UserFollowerMax int       `json:"user_follower_max" form:"user_follower_max" binding:"omitempty,gtefield=UserFollowerMin"`
	UserStatusMin   int       `json:"user_status_min" form:"user_status_min" binding:"omitempty"`
	UserStatusMax   int       `json:"user_status_max" form:"user_status_max" binding:"omitempty,gtefield=UserStatusMin"`
	Count           int       `json:"count" form:"count" binding:"omitempty,min=1,max=10000"`
	StartDate       time.Time `json:"start_date" form:"start_date" binding:"required" time_format:"2006-01-02 15:04"`
	EndDate         time.Time `json:"end_date" form:"end_date" binding:"required,gtefield=StartDate" time_format:"2006-01-02 15:04"`
}

type SearchForm struct {
	Hashtag   string    `json:"hashtag" form:"hashtag" binding:"required"`
	StartDate time.Time `json:"start_date" form:"start_date" binding:"required" time_format:"2006-01-02 15:04"`
	EndDate   time.Time `json:"end_date" form:"end_date" binding:"required,gtefield=StartDate" time_format:"2006-01-02 15:04"`
	Count     int       `json:"count" form:"count" binding:"min=1,max=10000"`
}
