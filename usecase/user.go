package usecase

import (
	"fmt"
	"sns-api/domain"
	"sns-api/logger"
	"time"
)

type UserUseCase interface {
	Search(name, description, language string, followerMin, followerMax, statusMin, statusMax, favoriteMin, favoriteMax, followMin, followMax, listMin, listMax int, srScoreMin, srScoreMax float64, startDate, endDate time.Time, count int, orderBy string) ([]*domain.User, int, error)
	GetById(userID uint64, startDate, endDate time.Time) (*domain.User, int, error)
	GetByIds(userIDs []uint64, startDate, endDate time.Time) ([]*domain.User, int, error)
}

type userUseCase struct {
	l              logger.Logging
	userRepository domain.UserRepository
}

func NewUserUseCase(l logger.Logging, ur domain.UserRepository) UserUseCase {
	return &userUseCase{
		l:              l,
		userRepository: ur,
	}
}

func (uu *userUseCase) Search(name, description, language string, followerMin, followerMax, statusMin, statusMax, favoriteMin, favoriteMax, followMin, followMax, listMin, listMax int, srScoreMin, srScoreMax float64, startDate, endDate time.Time, count int, orderBy string) ([]*domain.User, int, error) {
	users, hits, err := uu.userRepository.Search(name, description, language, followerMin, followerMax, statusMin, statusMax, favoriteMin, favoriteMax, followMin, followMax, listMin, listMax, srScoreMin, srScoreMax, startDate, endDate, count, orderBy)
	if err != nil {
		uu.l.Errorf(fmt.Sprintf("failed to GetById: %v", err))
		return nil, 0, err
	}
	return users, hits, nil
}

func (uu *userUseCase) GetById(userID uint64, startDate, endDate time.Time) (*domain.User, int, error) {
	user, hits, err := uu.userRepository.GetById(userID, startDate, endDate)
	if err != nil {
		uu.l.Errorf(fmt.Sprintf("failed to GetById: %v", err))
		return nil, 0, err
	}
	return user, hits, nil
}

func (uu *userUseCase) GetByIds(userIDs []uint64, startDate, endDate time.Time) ([]*domain.User, int, error) {
	users, hits, err := uu.userRepository.GetByIds(userIDs, startDate, endDate)
	if err != nil {
		uu.l.Errorf(fmt.Sprintf("failed to GetById: %v", err))
		return nil, 0, err
	}
	return users, hits, nil
}
