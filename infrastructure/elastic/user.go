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

type userRepository struct {
	l  logger.Logging
	es *elasticsearch.Client
}

func NewUserRepository(logger logger.Logging, conn *elasticsearch.Client) *userRepository {
	return &userRepository{
		l:  logger,
		es: conn,
	}
}

func (u *userRepository) Search(name, description, language string, followerMin, followerMax, statusMin, statusMax, favoriteMin, favoriteMax, followMin, followMax, listMin, listMax int, srScoreMin, srScoreMax float64, startDate, endDate time.Time, count int, orderBy string) ([]*domain.User, int, error) {
	var buf bytes.Buffer
	var users []*domain.User
	var filterQuery []map[string]interface{}
	var shouldQuery []map[string]interface{}

	buildQuery(&shouldQuery, name, "match_phrase", "screen_name")
	buildQuery(&shouldQuery, name, "match_phrase", "name")
	buildQuery(&shouldQuery, description, "match", "description")
	buildQuery(&filterQuery, language, "term", "language")
	buildQuery(&filterQuery, followerMin, "followers_count", "gte")
	buildQuery(&filterQuery, followerMax, "followers_count", "lte")
	buildQuery(&filterQuery, statusMin, "statuses_count", "gte")
	buildQuery(&filterQuery, statusMax, "statuses_count", "lte")
	buildQuery(&filterQuery, favoriteMin, "favourites_count", "gte")
	buildQuery(&filterQuery, favoriteMax, "favourites_count", "lte")
	buildQuery(&filterQuery, followMin, "friends_count", "gte")
	buildQuery(&filterQuery, followMax, "friends_count", "lte")
	buildQuery(&filterQuery, listMin, "listed_count", "gte")
	buildQuery(&filterQuery, listMax, "listed_count", "lte")
	buildQuery(&filterQuery, srScoreMin, "sr_score", "gte")
	buildQuery(&filterQuery, srScoreMax, "sr_score", "lte")

	ctx := context.Background()

	query := map[string]interface{}{
		"collapse": map[string]interface{}{
			"field": "id",
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should":               shouldQuery,
				"filter":               filterQuery,
				"minimum_should_match": 1,
			},
		},
		"sort": []map[string]interface{}{
			{
				orderBy: "desc",
			},
		},
	}

	if err := encodeQuery(&buf, query); err != nil {
		u.l.Errorf(fmt.Sprintf("failed to encode query: %s", err))
		return nil, 0, err
	}

	mDiff := monthDiff(startDate, endDate)
	monthList := buildIndexByTimeAdd(userIndex, startDate, mDiff)

	r, err := search(ctx, u.l, u.es, strings.Join(monthList, ","), &buf, count)
	if err != nil {
		u.l.Errorf(fmt.Sprintf("failed to search: %s", err))
		return nil, 0, err
	}

	hits := int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		description := ""
		score := 0.0
		if _, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})["description"].(string); ok {
			description = hit.(map[string]interface{})["_source"].(map[string]interface{})["description"].(string)
		}
		if _, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})["sr_score"].(float64); ok {
			score = hit.(map[string]interface{})["_source"].(map[string]interface{})["sr_score"].(float64)
		}
		user := domain.User{
			UserID:           hit.(map[string]interface{})["_source"].(map[string]interface{})["id"].(string),
			UserScreenName:   hit.(map[string]interface{})["_source"].(map[string]interface{})["screen_name"].(string),
			UserName:         hit.(map[string]interface{})["_source"].(map[string]interface{})["name"].(string),
			UserDescription:  description,
			UserImageProfile: hit.(map[string]interface{})["_source"].(map[string]interface{})["profile_image_url_https"].(string),
			Verified:         hit.(map[string]interface{})["_source"].(map[string]interface{})["verified"].(bool),
			FollowerCount:    hit.(map[string]interface{})["_source"].(map[string]interface{})["followers_count"].(float64),
			StatusCount:      hit.(map[string]interface{})["_source"].(map[string]interface{})["statuses_count"].(float64),
			FavoriteCount:    hit.(map[string]interface{})["_source"].(map[string]interface{})["favourites_count"].(float64),
			FollowCount:      hit.(map[string]interface{})["_source"].(map[string]interface{})["friends_count"].(float64),
			ListCount:        hit.(map[string]interface{})["_source"].(map[string]interface{})["listed_count"].(float64),
			SrScore:          score,
			CreatedAt:        hit.(map[string]interface{})["_source"].(map[string]interface{})["created_at"].(string),
		}
		users = append(users, &user)
	}
	return users, hits, nil
}

func (u *userRepository) GetById(userID uint64, startDate, endDate time.Time) (*domain.User, int, error) {
	var buf bytes.Buffer
	var user *domain.User
	ctx := context.Background()
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": []map[string]interface{}{
					{
						"match_phrase": map[string]interface{}{
							"id": userID,
						},
					},
				},
			},
		},
		"sort": []map[string]interface{}{
			{
				"inserted_at": "desc",
			},
		},
	}

	if err := encodeQuery(&buf, query); err != nil {
		u.l.Errorf(fmt.Sprintf("failed to encode query: %s", err))
		return nil, 0, err
	}

	mDiff := monthDiff(startDate, endDate)
	monthList := buildIndexByTimeAdd(userIndex, startDate, mDiff)

	r, err := search(ctx, u.l, u.es, strings.Join(monthList, ","), &buf, 1)
	if err != nil {
		u.l.Errorf(fmt.Sprintf("failed to search: %s", err))
		return nil, 0, err
	}

	hits := int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		description := ""
		score := 0.0
		if _, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})["description"].(string); ok {
			description = hit.(map[string]interface{})["_source"].(map[string]interface{})["description"].(string)
		}
		if _, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})["sr_score"].(float64); ok {
			score = hit.(map[string]interface{})["_source"].(map[string]interface{})["sr_score"].(float64)
		}
		user = &domain.User{
			UserID:           hit.(map[string]interface{})["_source"].(map[string]interface{})["id"].(string),
			UserScreenName:   hit.(map[string]interface{})["_source"].(map[string]interface{})["screen_name"].(string),
			UserName:         hit.(map[string]interface{})["_source"].(map[string]interface{})["name"].(string),
			UserDescription:  description,
			UserImageProfile: hit.(map[string]interface{})["_source"].(map[string]interface{})["profile_image_url_https"].(string),
			Verified:         hit.(map[string]interface{})["_source"].(map[string]interface{})["verified"].(bool),
			FollowerCount:    hit.(map[string]interface{})["_source"].(map[string]interface{})["followers_count"].(float64),
			StatusCount:      hit.(map[string]interface{})["_source"].(map[string]interface{})["statuses_count"].(float64),
			FavoriteCount:    hit.(map[string]interface{})["_source"].(map[string]interface{})["favourites_count"].(float64),
			FollowCount:      hit.(map[string]interface{})["_source"].(map[string]interface{})["friends_count"].(float64),
			ListCount:        hit.(map[string]interface{})["_source"].(map[string]interface{})["listed_count"].(float64),
			SrScore:          score,
			CreatedAt:        hit.(map[string]interface{})["_source"].(map[string]interface{})["created_at"].(string),
		}
	}
	return user, hits, nil
}

func (u *userRepository) GetByIds(userIDs []uint64, startDate, endDate time.Time) ([]*domain.User, int, error) {
	var buf bytes.Buffer
	var users []*domain.User
	var userQuery = make([]map[string]interface{}, 0, len(userIDs))
	ctx := context.Background()

	for _, id := range userIDs {
		q := map[string]interface{}{
			"match_phrase": map[string]interface{}{
				"id": id,
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
				"should": userQuery,
			},
		},
		"sort": []map[string]interface{}{
			{
				"inserted_at": "desc",
			},
		},
	}
	if err := encodeQuery(&buf, query); err != nil {
		u.l.Errorf(fmt.Sprintf("failed to encode query: %s", err))
		return nil, 0, err
	}

	mDiff := monthDiff(startDate, endDate)
	monthList := buildIndexByTimeAdd(userIndex, startDate, mDiff)

	r, err := search(ctx, u.l, u.es, strings.Join(monthList, ","), &buf, len(userIDs))
	if err != nil {
		u.l.Errorf(fmt.Sprintf("failed to search: %s", err))
		return nil, 0, err
	}

	hits := int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		description := ""
		score := 0.0
		if _, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})["description"].(string); ok {
			description = hit.(map[string]interface{})["_source"].(map[string]interface{})["description"].(string)
		}
		if _, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})["sr_score"].(float64); ok {
			score = hit.(map[string]interface{})["_source"].(map[string]interface{})["sr_score"].(float64)
		}
		user := domain.User{
			UserID:           hit.(map[string]interface{})["_source"].(map[string]interface{})["id"].(string),
			UserScreenName:   hit.(map[string]interface{})["_source"].(map[string]interface{})["screen_name"].(string),
			UserName:         hit.(map[string]interface{})["_source"].(map[string]interface{})["name"].(string),
			UserDescription:  description,
			UserImageProfile: hit.(map[string]interface{})["_source"].(map[string]interface{})["profile_image_url_https"].(string),
			Verified:         hit.(map[string]interface{})["_source"].(map[string]interface{})["verified"].(bool),
			FollowerCount:    hit.(map[string]interface{})["_source"].(map[string]interface{})["followers_count"].(float64),
			StatusCount:      hit.(map[string]interface{})["_source"].(map[string]interface{})["statuses_count"].(float64),
			FavoriteCount:    hit.(map[string]interface{})["_source"].(map[string]interface{})["favourites_count"].(float64),
			FollowCount:      hit.(map[string]interface{})["_source"].(map[string]interface{})["friends_count"].(float64),
			ListCount:        hit.(map[string]interface{})["_source"].(map[string]interface{})["listed_count"].(float64),
			SrScore:          score,
			CreatedAt:        hit.(map[string]interface{})["_source"].(map[string]interface{})["created_at"].(string),
		}
		users = append(users, &user)
	}
	return users, hits, nil
}
