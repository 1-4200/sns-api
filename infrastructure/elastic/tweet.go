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

const (
	tweetTypeNormal int = 1

	mediaTypeAll    int = -1
	mediaTypePhoto  int = 2
	mediaTypeVideo  int = 3
	mediaTypeGif    int = 4
)

type tweetRepository struct {
	l  logger.Logging
	es *elasticsearch.Client
}

func NewTweetRepository(logger logger.Logging, conn *elasticsearch.Client) *tweetRepository {
	return &tweetRepository{
		l:  logger,
		es: conn,
	}
}

func (t *tweetRepository) Get() ([]*domain.Tweet, error) {
	return []*domain.Tweet{}, nil
}

func (t *tweetRepository) GetByUser(userID uint64, startDate, endDate string, count int, orderBy string) ([]*domain.Tweet, int, error) {
	var r map[string]interface{}
	var buf bytes.Buffer
	var tweets []*domain.Tweet
	ctx := context.Background()

	query := map[string]interface{}{
		"collapse": map[string]interface{}{
			"field": "id",
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match_phrase": map[string]interface{}{
							"user_id": userID,
						},
					},
				},
				"filter": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"created_at": map[string]interface{}{
								"gte": fmt.Sprintf("%s:00", startDate),
								"lte": fmt.Sprintf("%s:59", endDate),
							},
						},
					},
					{
						"match": map[string]interface{}{
							"tweet_type": tweetTypeNormal,
						},
					},
				},
			},
		},
		"sort": []map[string]interface{}{
			{
				orderBy: "desc",
			},
		},
	}

	if err := encodeQuery(&buf, query); err != nil {
		t.l.Errorf(fmt.Sprintf("failed to encode query: %s", err))
		return nil, 0, err
	}

	r, err := search(ctx, t.l, t.es, fmt.Sprintf("%s-*", tweetIndex), &buf, count)
	if err != nil {
		t.l.Errorf(fmt.Sprintf("failed to search: %s", err))
		return nil, 0, err
	}

	hits := int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var tweetsUrls []*domain.TweetNestedURL
		createdAt, err := convertTime(hit.(map[string]interface{})["_source"].(map[string]interface{})["created_at"].(string))
		if err != nil {
			t.l.Error(fmt.Sprintf("failed to convert tweet time: %s", err))
		}
		if hit.(map[string]interface{})["_source"].(map[string]interface{})["nested_url"] != nil {
			for _, url := range hit.(map[string]interface{})["_source"].(map[string]interface{})["nested_url"].([]interface{}) {
				tweetsUrl := domain.TweetNestedURL{
					CanonicalURL: url.(map[string]interface{})["canonical_url"].(string),
					Domain:       url.(map[string]interface{})["domain"].(string),
				}
				tweetsUrls = append(tweetsUrls, &tweetsUrl)
			}
		}
		tweet := domain.Tweet{
			UserID:         hit.(map[string]interface{})["_source"].(map[string]interface{})["user_id"].(string),
			UserScreenName: hit.(map[string]interface{})["_source"].(map[string]interface{})["user_screen_name"].(string),
			UserName:       hit.(map[string]interface{})["_source"].(map[string]interface{})["user_name"].(string),
			TweetID:        hit.(map[string]interface{})["_id"].(string),
			Text:           hit.(map[string]interface{})["_source"].(map[string]interface{})["tweet"].(string),
			QuoteCount:     hit.(map[string]interface{})["_source"].(map[string]interface{})["quote_count"].(float64),
			FavoriteCount:  hit.(map[string]interface{})["_source"].(map[string]interface{})["favorite_count"].(float64),
			RetweetCount:   hit.(map[string]interface{})["_source"].(map[string]interface{})["retweet_count"].(float64),
			ReplyCount:     hit.(map[string]interface{})["_source"].(map[string]interface{})["reply_count"].(float64),
			CreatedAt:      createdAt,
			NestedURL:      tweetsUrls,
		}
		tweets = append(tweets, &tweet)
	}
	return tweets, hits, nil
}

func (t *tweetRepository) GetByUsers(userIDs []uint64, startDate, endDate time.Time, count int, orderBy string) ([]*domain.Tweet, int, error) {
	var r map[string]interface{}
	var buf bytes.Buffer
	var tweets []*domain.Tweet
	var userQuery = make([]map[string]interface{}, 0, len(userIDs))
	ctx := context.Background()

	for _, id := range userIDs {
		q := map[string]interface{}{
			"match_phrase": map[string]interface{}{
				"user_id": id,
			},
		}
		userQuery = append(userQuery, q)
	}

	query := map[string]interface{}{
		"collapse": map[string]interface{}{
			"field": "id",
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": []map[string]interface{}{
					{
						"bool": map[string]interface{}{
							"should": userQuery,
						},
					},
					{
						"match": map[string]interface{}{
							"tweet_type": tweetTypeNormal,
						},
					},
				},
			},
		},
		"sort": []map[string]interface{}{
			{
				orderBy: "desc",
			},
		},
	}

	if err := encodeQuery(&buf, query); err != nil {
		t.l.Errorf(fmt.Sprintf("failed to encode query: %s", err))
		return nil, 0, err
	}

	mDiff := monthDiff(startDate, endDate)
	monthList := buildIndexByTimeAdd(tweetIndex, startDate, mDiff)

	r, err := search(ctx, t.l, t.es, strings.Join(monthList, ","), &buf, count)
	if err != nil {
		t.l.Errorf(fmt.Sprintf("failed to search: %s", err))
		return nil, 0, err
	}

	hits := int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		createdAt, err := convertTime(hit.(map[string]interface{})["_source"].(map[string]interface{})["created_at"].(string))
		if err != nil {
			t.l.Error(fmt.Sprintf("failed to convert tweet time: %s", err))
		}
		tweet := domain.Tweet{
			UserID:         hit.(map[string]interface{})["_source"].(map[string]interface{})["user_id"].(string),
			UserScreenName: hit.(map[string]interface{})["_source"].(map[string]interface{})["user_screen_name"].(string),
			UserName:       hit.(map[string]interface{})["_source"].(map[string]interface{})["user_name"].(string),
			TweetID:        hit.(map[string]interface{})["_id"].(string),
			Text:           hit.(map[string]interface{})["_source"].(map[string]interface{})["tweet"].(string),
			QuoteCount:     hit.(map[string]interface{})["_source"].(map[string]interface{})["quote_count"].(float64),
			FavoriteCount:  hit.(map[string]interface{})["_source"].(map[string]interface{})["favorite_count"].(float64),
			RetweetCount:   hit.(map[string]interface{})["_source"].(map[string]interface{})["retweet_count"].(float64),
			ReplyCount:     hit.(map[string]interface{})["_source"].(map[string]interface{})["reply_count"].(float64),
			CreatedAt:      createdAt,
		}
		tweets = append(tweets, &tweet)
	}
	return tweets, hits, nil
}

func (t *tweetRepository) GetByDomain(userID uint64, startDate string, endDate string, count int, orderBy string, domainName string) ([]*domain.Tweet, int, []*domain.URL, error) {

	var esResultDomain map[string]interface{}
	var esResultURL map[string]interface{}
	var buf bytes.Buffer
	var tweets []*domain.Tweet
	var url_info []*domain.URL
	var tmpUrl map[string]interface{}
	var queryURLParamsUrl []map[string]interface{}

	ctx := context.Background()
	queryDomain := map[string]interface{}{
		"collapse": map[string]interface{}{
			"field": "id",
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match_phrase": map[string]interface{}{
							"user_id": userID,
						},
					},
					{
						"nested": map[string]interface{}{
							"path": "nested_url",
							"inner_hits": map[string]interface{}{},
							"query": map[string]interface{}{
								"match_phrase": map[string]interface{}{
									"nested_url.domain": domainName,
								},
							},
						},
					},
				},
				"filter": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"created_at": map[string]interface{}{
								"gte": fmt.Sprintf("%s:00", startDate),
								"lte": fmt.Sprintf("%s:59", endDate),
							},
						},
					},
					{
						"match": map[string]interface{}{
							"tweet_type": tweetTypeNormal,
						},
					},
				},
			},
		},
		"sort": []map[string]interface{}{
			{
				orderBy: "desc",
			},
		},
	}

	if errDomain := encodeQuery(&buf, queryDomain); errDomain != nil {
		t.l.Errorf(fmt.Sprintf("failed to encode query: %s", errDomain))
		return nil, 0, nil, errDomain
	}

	esResultDomain, errDomain := search(ctx, t.l, t.es, fmt.Sprintf("%s-*", tweetIndex), &buf, count)
	if errDomain != nil {
		t.l.Errorf(fmt.Sprintf("failed to search: %s", errDomain))
		return nil, 0, nil, errDomain
	}

	hits := int(esResultDomain["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	for _, hit := range esResultDomain["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var tweetsUrls []*domain.TweetNestedURL
		createdAt, errDomain := convertTime(hit.(map[string]interface{})["_source"].(map[string]interface{})["created_at"].(string))
		if errDomain != nil {
			t.l.Error(fmt.Sprintf("failed to convert tweet time: %s", errDomain))
		}
		for _, url := range hit.(map[string]interface{})["inner_hits"].(map[string]interface{})["nested_url"].(map[string]interface{})["hits"].(map[string]interface{})["hits"].([]interface{}) {
			tmpUrl = map[string]interface{}{
				"match_phrase": map[string]interface{}{
					"canonical_url": url.(map[string]interface{})["_source"].(map[string]interface{})["canonical_url"].(string),
				},
			}
			tweetsUrl := domain.TweetNestedURL{
				CanonicalURL: url.(map[string]interface{})["_source"].(map[string]interface{})["canonical_url"].(string),
				Domain:       url.(map[string]interface{})["_source"].(map[string]interface{})["domain"].(string),
			}
			tweetsUrls = append(tweetsUrls, &tweetsUrl)
			queryURLParamsUrl = append(queryURLParamsUrl, tmpUrl)
		}

		tweet := domain.Tweet{
			UserID:         hit.(map[string]interface{})["_source"].(map[string]interface{})["user_id"].(string),
			UserScreenName: hit.(map[string]interface{})["_source"].(map[string]interface{})["user_screen_name"].(string),
			UserName:       hit.(map[string]interface{})["_source"].(map[string]interface{})["user_name"].(string),
			TweetID:        hit.(map[string]interface{})["_id"].(string),
			Text:           hit.(map[string]interface{})["_source"].(map[string]interface{})["tweet"].(string),
			QuoteCount:     hit.(map[string]interface{})["_source"].(map[string]interface{})["quote_count"].(float64),
			FavoriteCount:  hit.(map[string]interface{})["_source"].(map[string]interface{})["favorite_count"].(float64),
			RetweetCount:   hit.(map[string]interface{})["_source"].(map[string]interface{})["retweet_count"].(float64),
			ReplyCount:     hit.(map[string]interface{})["_source"].(map[string]interface{})["reply_count"].(float64),
			CreatedAt:      createdAt,
			NestedURL:      tweetsUrls,
		}
		tweets = append(tweets, &tweet)
	}

	queryURL := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"bool": map[string]interface{}{
							"should": queryURLParamsUrl,
						},
					},
				},
			},
		},
	}

	if errURL := encodeQuery(&buf, queryURL); errURL != nil {
		t.l.Errorf(fmt.Sprintf("failed to encode query URL: %s", errURL))
		return nil, 0, nil, errURL
	}

	//上記Domainのクエリに対して複数のURLが想定されるのでMaxレコード数を増加
	esResultURL, errURL := search(ctx, t.l, t.es, fmt.Sprintf("%s-*", urlIndex), &buf, 10000)
	if errURL != nil {
		t.l.Errorf(fmt.Sprintf("failed to search URL: %s", errURL))
		return nil, 0, nil, errURL
	}

	hitsURL := int(esResultURL["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	if hitsURL > 0 {
		for _, hitURL := range esResultURL["hits"].(map[string]interface{})["hits"].([]interface{}) {
			var title = ""
			var description = ""
			if _, ok := hitURL.(map[string]interface{})["_source"].(map[string]interface{})["unwound"].(map[string]interface{}); ok {
				if hitURL.(map[string]interface{})["_source"].(map[string]interface{})["unwound"].(map[string]interface{})["title"] != nil {
					title = hitURL.(map[string]interface{})["_source"].(map[string]interface{})["unwound"].(map[string]interface{})["title"].(string)
				}
				if hitURL.(map[string]interface{})["_source"].(map[string]interface{})["unwound"].(map[string]interface{})["description"] != nil {
					description = hitURL.(map[string]interface{})["_source"].(map[string]interface{})["unwound"].(map[string]interface{})["description"].(string)
				}
			}
			tmpURL := domain.URL{
				URL:         hitURL.(map[string]interface{})["_source"].(map[string]interface{})["canonical_url"].(string),
				Title:       title,
				Description: description,
			}
			url_info = append(url_info, &tmpURL)
		}
	}

	t.l.Info("function elastic.GetByDomain done")
	return tweets, hits, url_info, nil
}

func (t *tweetRepository) GetByMediaType(userID uint64, startDate string, endDate string, count int, orderBy string, mediaType int) ([]*domain.TweetMedia, int, []*domain.Media, error) {

	var esResultTweet map[string]interface{}
	var esResultMedia map[string]interface{}
	var buf bytes.Buffer
	var tweets []*domain.TweetMedia
	var media []*domain.Media
	var queryMediaParamsTweetID []map[string]interface{}
	var queryTweetParamsMediaType []map[string]interface{}
	ctx := context.Background()
	if mediaType == mediaTypeAll {
		queryTweetParamsMediaType = []map[string]interface{}{
			{
				"match_phrase": map[string]interface{}{
					"media_type": mediaTypePhoto,
				},
			},
			{
				"match_phrase": map[string]interface{}{
					"media_type": mediaTypeVideo,
				},
			},
			{
				"match_phrase": map[string]interface{}{
					"media_type": mediaTypeGif,
				},
			},
		}
	} else {
		queryTweetParamsMediaType = []map[string]interface{}{
			{
				"match_phrase": map[string]interface{}{
					"media_type": mediaType,
				},
			},
		}
	}

	queryTweet := map[string]interface{}{
		"collapse": map[string]interface{}{
			"field": "id",
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match_phrase": map[string]interface{}{
							"user_id": userID,
						},
					},
					{
						"bool": map[string]interface{}{
							"should": queryTweetParamsMediaType,
						},
					},
				},
				"filter": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"created_at": map[string]interface{}{
								"gte": fmt.Sprintf("%s:00", startDate),
								"lte": fmt.Sprintf("%s:59", endDate),
							},
						},
					},
					{
						"match": map[string]interface{}{
							"tweet_type": tweetTypeNormal,
						},
					},
				},
			},
		},
		"sort": []map[string]interface{}{
			{
				orderBy: "desc",
			},
		},
	}

	if errTweet := encodeQuery(&buf, queryTweet); errTweet != nil {
		t.l.Errorf(fmt.Sprintf("failed to encode query: %s", errTweet))
		return nil, 0, nil, errTweet
	}

	esResultTweet, errTweet := search(ctx, t.l, t.es, fmt.Sprintf("%s-*", tweetIndex), &buf, count)
	if errTweet != nil {
		t.l.Errorf(fmt.Sprintf("failed to search: %s", errTweet))
		return nil, 0, nil, errTweet
	}

	hits := int(esResultTweet["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	for _, hit := range esResultTweet["hits"].(map[string]interface{})["hits"].([]interface{}) {

		tmptweetID := map[string]interface{}{
			"match_phrase": map[string]interface{}{
				"source_status_id": hit.(map[string]interface{})["_source"].(map[string]interface{})["id"].(string),
			},
		}

		Media := domain.TweetMedia{
			TweetID:       hit.(map[string]interface{})["_id"].(string),
			MediaType:     hit.(map[string]interface{})["_source"].(map[string]interface{})["media_type"].(float64),
			FavoriteCount: hit.(map[string]interface{})["_source"].(map[string]interface{})["favorite_count"].(float64),
			RetweetCount:  hit.(map[string]interface{})["_source"].(map[string]interface{})["retweet_count"].(float64),
		}
		tweets = append(tweets, &Media)
		queryMediaParamsTweetID = append(queryMediaParamsTweetID, tmptweetID)
	}

	queryMedia := map[string]interface{}{
		"collapse": map[string]interface{}{
			"field": "id",
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": queryMediaParamsTweetID,
			},
		},
	}

	if errMedia := encodeQuery(&buf, queryMedia); errMedia != nil {
		t.l.Errorf(fmt.Sprintf("failed to encode query URL: %s", errMedia))
		return nil, 0, nil, errMedia
	}

	esResultMedia, errMedia := search(ctx, t.l, t.es, fmt.Sprintf("%s-*", mediaIndex), &buf, count)
	if errMedia != nil {
		t.l.Errorf(fmt.Sprintf("failed to search URL: %s", errMedia))
		return nil, 0, nil, errMedia
	}

	for _, hitMedia := range esResultMedia["hits"].(map[string]interface{})["hits"].([]interface{}) {
		tmpMedia := domain.Media{
			TweetID:  hitMedia.(map[string]interface{})["_source"].(map[string]interface{})["source_status_id"].(string),
			MediaURL: hitMedia.(map[string]interface{})["_source"].(map[string]interface{})["media_url_https"].(string),
		}
		media = append(media, &tmpMedia)
	}
	t.l.Info("function elastic.GetByMedia done")
	return tweets, hits, media, nil
}
