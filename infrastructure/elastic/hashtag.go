package elastic

import (
	"bytes"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"sns-api/domain"
	"sns-api/logger"
	"strings"
	"time"
)

type hashtagRepository struct {
	l  logger.Logging
	es *elasticsearch.Client
}

func NewHashtagRepository(logger logger.Logging, conn *elasticsearch.Client) *hashtagRepository {
	return &hashtagRepository{
		l:  logger,
		es: conn,
	}
}

func (t *hashtagRepository) Get(keyword string, hashtag []string, tweetType []int, retweetMin, retweetMax, quoteMin, quoteMax, favoriteMin, favoriteMax int, userInclude, userExclude, hashtagInclude, hashtagExclude []string, userFollowerMin, userFollowerMax, userStatusMin, userStatusMax, count int, startDate, endDate time.Time) ([]*domain.Hashtag, int, error) {
	var buf bytes.Buffer
	var hashtags []*domain.Hashtag
	ctx := context.Background()
	var mustQuery []map[string]interface{}
	var mustNotQuery []map[string]interface{}

	buildQuery(&mustQuery, keyword, "match_phrase", "tweet")
	buildQuery(&mustQuery, hashtag, "wildcard", "hashtag")
	buildQuery(&mustQuery, tweetType, "terms", "tweet_type")
	buildQuery(&mustQuery, retweetMin, "retweet_count", "gte")
	buildQuery(&mustQuery, retweetMax, "retweet_count", "lte")
	buildQuery(&mustQuery, quoteMin, "quote_count", "gte")
	buildQuery(&mustQuery, quoteMax, "quote_count", "lte")
	buildQuery(&mustQuery, favoriteMin, "favorite_count", "gte")
	buildQuery(&mustQuery, favoriteMax, "favorite_count", "lte")
	buildQuery(&mustQuery, userInclude, "terms", "user_screen_name")
	buildQuery(&mustQuery, hashtagInclude, "match", "hashtag")
	buildQuery(&mustQuery, userFollowerMin, "user_followers_count", "gte")
	buildQuery(&mustQuery, userFollowerMax, "user_followers_count", "lte")
	buildQuery(&mustQuery, userStatusMin, "user_statuses_count", "gte")
	buildQuery(&mustQuery, userStatusMax, "user_statuses_count", "lte")
	buildQuery(&mustQuery, startDate, "created_at", "gte")
	buildQuery(&mustQuery, endDate, "created_at", "lte")
	buildQuery(&mustNotQuery, userExclude, "terms", "user_screen_name")
	buildQuery(&mustNotQuery, hashtagExclude, "match", "hashtag")

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":     mustQuery,
				"must_not": mustNotQuery,
			},
		},
		"aggs": map[string]interface{}{
			"distinct_hashtag_count": map[string]interface{}{
				"cardinality": map[string]interface{}{
					"field":               "hashtag",
					"precision_threshold": count,
				},
			},
			"group_by_hashtag": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "hashtag",
					"order": map[string]interface{}{
						"_count": "desc",
					},
					"size": count,
				},
				"aggs": map[string]interface{}{
					"retweet_avg": map[string]interface{}{
						"avg": map[string]interface{}{
							"field": "retweet_count",
						},
					},
					"retweet_sum": map[string]interface{}{
						"sum": map[string]interface{}{
							"field": "retweet_count",
						},
					},
					"favorite_avg": map[string]interface{}{
						"avg": map[string]interface{}{
							"field": "favorite_count",
						},
					},
					"favorite_sum": map[string]interface{}{
						"sum": map[string]interface{}{
							"field": "favorite_count",
						},
					},
					"quote_avg": map[string]interface{}{
						"avg": map[string]interface{}{
							"field": "quote_count",
						},
					},
					"quote_sum": map[string]interface{}{
						"sum": map[string]interface{}{
							"field": "quote_count",
						},
					},
					"reply_avg": map[string]interface{}{
						"avg": map[string]interface{}{
							"field": "reply_count",
						},
					},
					"reply_sum": map[string]interface{}{
						"sum": map[string]interface{}{
							"field": "reply_count",
						},
					},
				},
			},
		},
	}

	if err := encodeQuery(&buf, query); err != nil {
		t.l.Errorf(fmt.Sprintf("failed to encode query: %s", err))
		return nil, 0, err
	}

	mDiff := monthDiff(startDate, endDate)
	monthList := buildIndexByTimeAdd(tweetIndex, startDate, mDiff)
	r, err := search(ctx, t.l, t.es, strings.Join(monthList, ","), &buf, 0)
	if err != nil {
		t.l.Errorf(fmt.Sprintf("failed to search: %s", err))
		return nil, 0, err
	}

	hits := int(r["aggregations"].(map[string]interface{})["distinct_hashtag_count"].(map[string]interface{})["value"].(float64))
	for _, agg := range r["aggregations"].(map[string]interface{})["group_by_hashtag"].(map[string]interface{})["buckets"].([]interface{}) {
		hashtag := domain.Hashtag{
			Hashtag:       agg.(map[string]interface{})["key"].(string),
			StatusCount:   uint64(agg.(map[string]interface{})["doc_count"].(float64)),
			RetweetAvg:    agg.(map[string]interface{})["retweet_avg"].(map[string]interface{})["value"].(float64),
			RetweetCount:  uint64(agg.(map[string]interface{})["retweet_sum"].(map[string]interface{})["value"].(float64)),
			FavoriteAvg:   agg.(map[string]interface{})["favorite_avg"].(map[string]interface{})["value"].(float64),
			FavoriteCount: uint64(agg.(map[string]interface{})["favorite_sum"].(map[string]interface{})["value"].(float64)),
			ReplyAvg:      agg.(map[string]interface{})["reply_avg"].(map[string]interface{})["value"].(float64),
			ReplyCount:    uint64(agg.(map[string]interface{})["reply_sum"].(map[string]interface{})["value"].(float64)),
			QuoteAvg:      agg.(map[string]interface{})["quote_avg"].(map[string]interface{})["value"].(float64),
			QuoteCount:    uint64(agg.(map[string]interface{})["quote_sum"].(map[string]interface{})["value"].(float64)),
		}
		hashtags = append(hashtags, &hashtag)
	}
	return hashtags, hits, nil
}

func (t *hashtagRepository) Search(hashtag string, startDate, endDate time.Time, count int) ([]*domain.HashtagBySearch, int, error) {
	var buf bytes.Buffer
	var hashtags []*domain.HashtagBySearch
	ctx := context.Background()
	var mustQuery []map[string]interface{}

	buildQuery(&mustQuery, hashtag, "wildcard", "hashtag")
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":     mustQuery,
			},
		},
		"aggs": map[string]interface{}{
			"distinct_hashtag_count": map[string]interface{}{
				"cardinality": map[string]interface{}{
					"field":               "hashtag",
					"precision_threshold": count,
				},
			},
			"group_by_hashtag": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "hashtag",
					"order": map[string]interface{}{
						"_count": "desc",
					},
					"size": count,
				},
			},
		},
	}

	if err := encodeQuery(&buf, query); err != nil {
		t.l.Errorf(fmt.Sprintf("failed to encode query: %s", err))
		return nil, 0, err
	}

	mDiff := monthDiff(startDate, endDate)
	monthList := buildIndexByTimeAdd(tweetIndex, startDate, mDiff)
	r, err := search(ctx, t.l, t.es, strings.Join(monthList, ","), &buf, 0)
	if err != nil {
		t.l.Errorf(fmt.Sprintf("failed to search: %s", err))
		return nil, 0, err
	}

	hits := int(r["aggregations"].(map[string]interface{})["distinct_hashtag_count"].(map[string]interface{})["value"].(float64))
	for _, agg := range r["aggregations"].(map[string]interface{})["group_by_hashtag"].(map[string]interface{})["buckets"].([]interface{}) {
		hashtag := domain.HashtagBySearch{
			Hashtag:       agg.(map[string]interface{})["key"].(string),
			StatusCount:   uint64(agg.(map[string]interface{})["doc_count"].(float64)),
		}
		hashtags = append(hashtags, &hashtag)
	}
	return hashtags, hits, nil
}
