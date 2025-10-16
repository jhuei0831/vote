package service

import (
	"vote/app/database"
	"vote/app/model"
)

type BallotService struct {
}

func NewBallotService() BallotService {
	return BallotService{}
}

// CreateBallots 建立投票
func (b BallotService) CreateBallots(voter uint64, selectedCandidates map[uint64]map[uint64]bool) error {

	transaction := database.SqlSession.Begin()
	for questionId, candidates := range selectedCandidates {
		ballot := model.Ballot{
			PasswordID: voter,
			QuestionID: questionId,
		}
		err := transaction.Create(&ballot).Error
		if err != nil {
			transaction.Rollback()
			return err
		}

		for cid := range candidates {
			ballotSelect := model.BallotSelect{
				BallotID:    ballot.ID,
				CandidateID: cid,
			}
			err = transaction.Create(&ballotSelect).Error
			if err != nil {
				transaction.Rollback()
				return err
			}
		}
	}

	err := transaction.Commit().Error

	if err != nil {
		return err
	}

	return nil
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
