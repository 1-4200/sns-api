package domain

import "time"

type Hashtag struct {
	Hashtag       string  `json:"hashtag"`
	StatusCount   uint64  `json:"status_count"`
	RetweetAvg    float64 `json:"retweet_avg"`
	RetweetCount  uint64  `json:"retweet_count"`
	FavoriteAvg   float64 `json:"favorite_avg"`
	FavoriteCount uint64  `json:"favorite_count"`
	ReplyAvg      float64 `json:"reply_avg"`
	ReplyCount    uint64  `json:"reply_count"`
	QuoteAvg      float64 `json:"quote_avg"`
	QuoteCount    uint64  `json:"quote_count"`
}

type HashtagBySearch struct {
	Hashtag       string  `json:"hashtag"`
	StatusCount   uint64  `json:"status_count"`
}

type HashtagRepository interface {
	Get(keyword string, hashtag []string, tweetType []int, retweetMin, retweetMax, quoteMin, quoteMax, favoriteMin, favoriteMax int, userInclude, userExclude, hashtagInclude, hashtagExclude []string, userFollowerMin, userFollowerMax, userStatusMin, userStatusMax, count int, startDate, endDate time.Time) ([]*Hashtag, int, error)
	Search(hashtag string, startDate, endDate time.Time, count int) ([]*HashtagBySearch, int, error)
}
