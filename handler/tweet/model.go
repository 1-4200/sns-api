package tweet

import "time"

type Response struct {
	Hits int         `json:"hits"`
	Res  interface{} `json:"res"`
}

type ResponseDomain struct {
	Hits    int         `json:"hits"`
	Tweets  interface{} `json:"tweets"`
	UrlInfo interface{} `json:"url_info"`
}

type ResponseMedia struct {
	Hits   int         `json:"hits"`
	Tweets interface{} `json:"tweets"`
	Media  interface{} `json:"media"`
}

type ResponseTransition struct {
	Hits        int         `json:"hits"`
	Transitions interface{} `json:"transitions"`
}

type UserForm struct {
	UserID    uint64    `json:"user_id" form:"user_id" binding:"required"`
	StartDate time.Time `json:"start_date" form:"start_date" binding:"required" time_format:"2006-01-02 15:04"`
	EndDate   time.Time `json:"end_date" form:"end_date" binding:"required,gtefield=StartDate" time_format:"2006-01-02 15:04"`
	OrderBy   string    `json:"order_by" form:"order_by" binding:"omitempty,oneof=retweet_count quote_count favorite_count created_at inserted_at"`
	Count     int       `json:"count" form:"count" binding:"omitempty,min=1,max=10000"`
}

type UsersForm struct {
	UserIDs   []uint64  `json:"user_ids" form:"user_ids" binding:"required,max=10000"`
	StartDate time.Time `json:"start_date" form:"start_date" binding:"required" time_format:"2006-01-02 15:04"`
	EndDate   time.Time `json:"end_date" form:"end_date" binding:"required,gtefield=StartDate" time_format:"2006-01-02 15:04"`
	OrderBy   string    `json:"order_by" form:"order_by" binding:"omitempty,oneof=retweet_count quote_count favorite_count created_at inserted_at"`
	Count     int       `json:"count" form:"count" binding:"omitempty,min=1,max=10000"`
}

type URLForm struct {
	UserID    uint64    `json:"user_id" form:"user_id" binding:"required"`
	StartDate time.Time `json:"start_date" form:"start_date" binding:"required" time_format:"2006-01-02 15:04"`
	EndDate   time.Time `json:"end_date" form:"end_date" binding:"required,gtefield=StartDate" time_format:"2006-01-02 15:04"`
	OrderBy   string    `json:"order_by" form:"order_by" binding:"omitempty,oneof=retweet_count quote_count favorite_count created_at inserted_at"`
	Count     int       `json:"count" form:"count" binding:"omitempty,min=1,max=10000"`
	Domain    string    `json:"domain" form:"domain" binding:"required"`
}

type MediaForm struct {
	UserID    uint64    `json:"user_id" form:"user_id" binding:"required"`
	StartDate time.Time `json:"start_date" form:"start_date" binding:"required" time_format:"2006-01-02 15:04"`
	EndDate   time.Time `json:"end_date" form:"end_date" binding:"required,gtefield=StartDate" time_format:"2006-01-02 15:04"`
	OrderBy   string    `json:"order_by" form:"order_by" binding:"omitempty,oneof=retweet_count quote_count favorite_count created_at inserted_at"`
	Count     int       `json:"count" form:"count" binding:"omitempty,min=1,max=10000"`
	MediaType int       `json:"media_type" form:"media_type" binding:"required,oneof=-1 2 3 4"`
}

type TransitionForm struct {
	UserID    uint64    `json:"user_id" form:"user_id" binding:"required"`
	StartDate time.Time `json:"start_date" form:"start_date" binding:"required" time_format:"2006-01-02"`
	EndDate   time.Time `json:"end_date" form:"end_date" binding:"required,gtefield=StartDate" time_format:"2006-01-02"`
	Count     int       `json:"count" form:"count" binding:"omitempty,min=1,max=100000"`
}
