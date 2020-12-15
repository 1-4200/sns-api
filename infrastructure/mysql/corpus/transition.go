package corpus

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sns-api/domain"
	"sns-api/logger"
)

type tweetRepository struct {
	l  logger.Logging
	db *sql.DB
}

func NewTweetRepository(logger logger.Logging, db *sql.DB) *tweetRepository {
	return &tweetRepository{
		l:  logger,
		db: db,
	}
}

func (t *tweetRepository) GetTransitionByUser(userID uint64, startDate, endDate string, count int) ([]*domain.TweetTransition, error) {
	sql := `SELECT user_id, followers_count, friends_count, listed_count, favourites_count, statuses_count, created_at
			FROM tw_fullarchive_user_data
			WHERE user_id = ?
			  AND created_at BETWEEN ? AND ?
			ORDER BY created_at DESC
			LIMIT ?`
	rows, err := t.db.Query(sql, userID, startDate, endDate, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tts []*domain.TweetTransition
	for rows.Next() {
		tt := &domain.TweetTransition{}
		if err = rows.Scan(&tt.UserID, &tt.FollowerCount, &tt.FriendCount, &tt.ListedCount, &tt.FavoriteCount, &tt.StatusCount, &tt.CreatedAt); err != nil {
			return nil, err
		}
		tts = append(tts, tt)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tts, nil
}
