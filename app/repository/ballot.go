package repository

import (
	"fmt"
	"vote/app/database"
	"vote/app/model"
)

type BallotRepository struct {
}

func NewBallotRepository() BallotRepository {
	return BallotRepository{}
}

// CreateBallots 建立投票
func (b BallotRepository) CreateBallots(voter uint64, ballotSelections model.BallotCreate) error {
	// check if voter has voted
	var existingCount int64
	err := database.SqlSession.Model(&model.Ballot{}).
		Where("password_id = ?", voter).
		Count(&existingCount).Error

	if err != nil {
		return err
	}

	if existingCount > 0 {
		return fmt.Errorf("voter %d has already voted", voter)
	}

	tx := database.SqlSession.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, question := range ballotSelections.Selections {
		// 收集選中的候選人
		selectedCandidates := make([]uint64, 0)
		for _, candidate := range question.Candidates {
			if candidate.IsSelected {
				selectedCandidates = append(selectedCandidates, candidate.ID)
			}
		}

		// 如果沒有選中任何候選人，跳過此問題
		if len(selectedCandidates) == 0 {
			continue
		}

		// 建立投票記錄
		ballot := model.Ballot{
			PasswordID: voter,
			QuestionID: question.ID,
		}
		if err := tx.Create(&ballot).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 批量建立選擇記錄
		ballotSelects := make([]model.BallotSelect, len(selectedCandidates))
		for i, cid := range selectedCandidates {
			ballotSelects[i] = model.BallotSelect{
				BallotID:    ballot.ID,
				CandidateID: cid,
			}
		}

		if err := tx.Create(&ballotSelects).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}