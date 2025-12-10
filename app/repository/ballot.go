package repository

import (
	"fmt"
	"vote/app/database"
	"vote/app/model"

	"github.com/google/uuid"
)

type BallotRepository struct {
}

func NewBallotRepository() BallotRepository {
	return BallotRepository{}
}

// GetBallotByVoterId 根據投票者ID獲取選票
func (b BallotRepository) GetBallotByVoterId(voterId uint64) ([]model.Ballot) {
	var ballots []model.Ballot
	database.SqlSession.Where("password_id = ?", voterId).Preload("BallotSelects").Find(&ballots)
	return ballots
}

// GetBallotByVoteId 根據投票ID獲取選票
func (b BallotRepository) GetBallots(ballotQuery model.BallotQuery) ([]model.Ballot, int64, error) {
	var ballots []model.Ballot
	var total int64

	query := database.SqlSession.Model(&model.Ballot{}).Preload("BallotSelects")

	if ballotQuery.VoteID != uuid.Nil {
		query = query.Joins("JOIN questions ON ballots.question_id = questions.id").
			Where("questions.vote_id = ?", ballotQuery.VoteID)
	}

	if ballotQuery.QuestionID != 0 {
		query = query.Where("ballots.question_id = ?", ballotQuery.QuestionID)
	}

	if ballotQuery.VoterID != 0 {
		query = query.Where("ballots.password_id = ?", ballotQuery.VoterID)
	}

	// 計算總筆數
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 取得資料
	// 使用分頁服務處理分頁
	paginationRepository := NewPaginationRepository[*model.BallotQuery, model.Ballot]()
	query, err = paginationRepository.Handler(query, &ballotQuery)
	if err != nil {
		return nil, 0, err
	}

	// 查詢資料
	err = query.Find(&ballots).Error
	if err != nil {
		return nil, 0, err
	}

	return ballots, total, nil
}

// CreateBallots 建立投票
func (b BallotRepository) CreateBallots(voter uint64, ballotSelections model.BallotCreate) error {
	// check if voter has voted
	isVoted, err := b.CheckIfVoterHasVoted(voter)

	if err != nil {
		return err
	}

	if isVoted {
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
				selectedCandidates = append(selectedCandidates, candidate.CandidateID)
			}
		}

		// 如果沒有選中任何候選人，跳過此問題
		if len(selectedCandidates) == 0 {
			continue
		}

		// 建立投票記錄
		ballot := model.Ballot{
			PasswordID: voter,
			QuestionID: question.QuestionID,
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

// CheckIfVoterHasVoted 檢查投票者是否已經投票
func (b BallotRepository) CheckIfVoterHasVoted(voterId uint64) (bool, error) {
	var count int64
	err := database.SqlSession.Model(&model.Ballot{}).
		Where("password_id = ?", voterId).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// DeleteBallot 刪除選票
func (b BallotRepository) DeleteBallot(voterID uint64) error {
	err := database.SqlSession.Where("password_id = ?", voterID).Delete(&model.Ballot{}).Error
	if err != nil {
		return err
	}
	
	return nil
}