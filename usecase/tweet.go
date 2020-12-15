package usecase

import (
	"fmt"
	"sns-api/domain"
	"sns-api/logger"
	"time"
)

type TweetUseCase interface {
	Get() ([]*domain.Tweet, error)
	GetByUser(userID uint64, startDate, endDate string, count int, orderBy string) ([]*domain.Tweet, int, error)
	GetByUsers(userIDs []uint64, startDate, endDate time.Time, count int, orderBy string) ([]*domain.Tweet, int, error)
	GetByDomain(userID uint64, startDate, endDate string, count int, orderBy string, domainName string) ([]*domain.Tweet, int, []*domain.URL, error)
	GetByMediaType(userID uint64, startDate, endDate string, count int, orderBy string, mediaType int) ([]*domain.TweetMedia, int, []*domain.Media, error)
	GetTransitionByUser(userID uint64, startDate, endDate string, count int) ([]*domain.TweetTransition, error)
}

type tweetUseCase struct {
	l                    logger.Logging
	tweetRepository      domain.TweetRepository
	transitionRepository domain.TransitionRepository
}

func NewTweetUseCase(l logger.Logging, tr domain.TweetRepository, tts domain.TransitionRepository) TweetUseCase {
	return &tweetUseCase{
		l:                    l,
		tweetRepository:      tr,
		transitionRepository: tts,
	}
}

func (t *tweetUseCase) Get() ([]*domain.Tweet, error) {
	tweets, err := t.tweetRepository.Get()
	if err != nil {
		t.l.Errorf(fmt.Sprintf("failed to Get: %v", err))
		return nil, err
	}
	return tweets, nil
}

func (t *tweetUseCase) GetByUser(userID uint64, startDate, endDate string, count int, orderBy string) ([]*domain.Tweet, int, error) {
	tweets, hits, err := t.tweetRepository.GetByUser(userID, startDate, endDate, count, orderBy)
	if err != nil {
		t.l.Errorf(fmt.Sprintf("failed to GetByUser: %v", err))
		return nil, 0, err
	}
	return tweets, hits, nil
}

func (t *tweetUseCase) GetByUsers(userIDs []uint64, startDate, endDate time.Time, count int, orderBy string) ([]*domain.Tweet, int, error) {
	tweets, hits, err := t.tweetRepository.GetByUsers(userIDs, startDate, endDate, count, orderBy)
	if err != nil {
		t.l.Errorf(fmt.Sprintf("failed to GetByUsers: %v", err))
		return nil, 0, err
	}
	return tweets, hits, nil
}

func (t *tweetUseCase) GetByDomain(userID uint64, startDate, endDate string, count int, orderBy string, domainName string) ([]*domain.Tweet, int, []*domain.URL, error) {
	tweets, hits, urlInfo, err := t.tweetRepository.GetByDomain(userID, startDate, endDate, count, orderBy, domainName)
	if err != nil {
		t.l.Errorf(fmt.Sprintf("failed to GetByDomain: %v", err))
		return nil, 0, nil, err
	}
	return tweets, hits, urlInfo, nil
}

func (t *tweetUseCase) GetByMediaType(userID uint64, startDate, endDate string, count int, orderBy string, mediaType int) ([]*domain.TweetMedia, int, []*domain.Media, error) {
	tweets, hits, media, err := t.tweetRepository.GetByMediaType(userID, startDate, endDate, count, orderBy, mediaType)
	t.l.Info("function usecase.GetByMedia done")
	if err != nil {
		t.l.Errorf(fmt.Sprintf("failed to GetByMediaType: %v", err))
		return nil, 0, nil, err
	}
	return tweets, hits, media, nil
}

func (t *tweetUseCase) GetTransitionByUser(userID uint64, startDate, endDate string, count int) ([]*domain.TweetTransition, error) {
	tts, err := t.transitionRepository.GetTransitionByUser(userID, startDate, endDate, count)
	if err != nil {
		t.l.Errorf(fmt.Sprintf("failed to GetTransitionByUser: %v", err))
		return nil, err
	}
	return tts, nil
}
