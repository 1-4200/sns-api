package usecase

import (
	"fmt"
	"sns-api/domain"
	"sns-api/logger"
	"time"
)

type HashtagUseCase interface {
	Get(keyword string, hashtag []string, tweetType []int, retweetMin, retweetMax, quoteMin, quoteMax, favoriteMin, favoriteMax int, userInclude, userExclude, hashtagInclude, hashtagExclude []string, userFollowerMin, userFollowerMax, userStatusMin, userStatusMax, count int, startDate, endDate time.Time) ([]*domain.Hashtag, int, error)
	Search(hashtag string, startDate, endDate time.Time, count int) ([]*domain.HashtagBySearch, int, error)
}

type hashtagUseCase struct {
	l                 logger.Logging
	hashtagRepository domain.HashtagRepository
}

func NewHashtagUseCase(l logger.Logging, hr domain.HashtagRepository) HashtagUseCase {
	return &hashtagUseCase{
		l:                 l,
		hashtagRepository: hr,
	}
}

func (h *hashtagUseCase) Get(keyword string, hashtag []string, tweetType []int, retweetMin, retweetMax, quoteMin, quoteMax, favoriteMin, favoriteMax int, userInclude, userExclude, hashtagInclude, hashtagExclude []string, userFollowerMin, userFollowerMax, userStatusMin, userStatusMax, count int, startDate, endDate time.Time) ([]*domain.Hashtag, int, error) {
	hashtags, hits, err := h.hashtagRepository.Get(keyword, hashtag, tweetType, retweetMin, retweetMax, quoteMin, quoteMax, favoriteMin, favoriteMax, userInclude, userExclude, hashtagInclude, hashtagExclude, userFollowerMin, userFollowerMax, userStatusMin, userStatusMax, count, startDate, endDate)
	if err != nil {
		h.l.Errorf(fmt.Sprintf("failed to Get: %v", err))
		return nil, 0, err
	}
	return hashtags, hits, nil
}

func (h *hashtagUseCase) Search(hashtag string, startDate, endDate time.Time, count int) ([]*domain.HashtagBySearch, int, error) {
	hashtags, hits, err := h.hashtagRepository.Search(hashtag, startDate, endDate, count)
	if err != nil {
		h.l.Errorf(fmt.Sprintf("failed to Search: %v", err))
		return nil, 0, err
	}
	return hashtags, hits, nil
}
