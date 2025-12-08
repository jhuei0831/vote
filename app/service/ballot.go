package service

import (
	"vote/app/database"
	"vote/app/model"
	"vote/app/repository"
)

type BallotService struct {
}

func NewBallotService() BallotService {
	return BallotService{}
}

// CreateBallots 建立投票
func (b BallotService) CreateBallots(voter uint64, ballotSelections model.BallotCreate) ([]*model.Ballot, error) {
	repository := repository.NewBallotRepository()
	err := repository.CreateBallots(voter, ballotSelections)
	if err != nil {
		return nil, err
	}

	var ballots []*model.Ballot
	err = database.SqlSession.Where("password_id = ?", voter).Preload("BallotSelects").Find(&ballots).Error

	if err != nil {
		return nil, err
	}

	return ballots, nil
}

// CheckIfVoterHasVoted 檢查投票者是否已經投票
func (b BallotService) CheckIfVoterHasVoted(voterId uint64) (bool, error) {
	var count int64
	err := database.SqlSession.Model(&model.Ballot{}).
		Where("password_id = ?", voterId).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetBallotByVoterId 根據投票者ID獲取選票
func (b BallotService) GetBallotByVoterId(voterId uint64) ([][]string) {
	// database.SqlSession.Model(&model.Ballot{}).Where("vote")

	return make([][]string, 1)
}
