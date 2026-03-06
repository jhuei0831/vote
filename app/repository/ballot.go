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

// GetBallotByVoterId gets ballots by voter ID
func (b BallotRepository) GetBallotByVoterId(voterId uint64) []model.Ballot {
	var ballots []model.Ballot
	database.SqlSession.Where("invitation_id = ?", voterId).Preload("BallotSelects").Find(&ballots)
	return ballots
}

// GetBallots gets ballots based on query parameters
func (b BallotRepository) GetBallots(ballotQuery model.BallotQuery) ([]model.Ballot, int64, error) {
	var ballots []model.Ballot
	var total int64

	query := database.SqlSession.Model(&model.Ballot{}).Preload("BallotSelects")

	if ballotQuery.SessionID != uuid.Nil {
		query = query.Joins("JOIN polls ON ballots.poll_id = polls.id").
			Where("polls.vote_id = ?", ballotQuery.SessionID)
	}

	if ballotQuery.PollID != 0 {
		query = query.Where("ballots.poll_id = ?", ballotQuery.PollID)
	}

	if ballotQuery.VoterID != 0 {
		query = query.Where("ballots.invitation_id = ?", ballotQuery.VoterID)
	}

	// Count total records
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get data
	// Use pagination service to handle pagination
	paginationRepository := NewPaginationRepository[*model.BallotQuery, model.Ballot]()
	query, err = paginationRepository.Handler(query, &ballotQuery)
	if err != nil {
		return nil, 0, err
	}

	// Query data
	err = query.Find(&ballots).Error
	if err != nil {
		return nil, 0, err
	}

	return ballots, total, nil
}

// CreateBallots creates ballots for a voter
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

	for _, poll := range ballotSelections.Selections {
		// Collect selected poll options
		selectedPollOptions := make([]uint64, 0)
		for _, pollOption := range poll.PollOptions {
			if pollOption.IsSelected {
				selectedPollOptions = append(selectedPollOptions, pollOption.PollOptionID)
			}
		}

		// Skip this poll if no poll options are selected
		if len(selectedPollOptions) == 0 {
			continue
		}

		// Create ballot record
		ballot := model.Ballot{
			InvitationID: voter,
			PollID:       poll.PollID,
		}
		if err := tx.Create(&ballot).Error; err != nil {
			tx.Rollback()
			return err
		}

		// Batch create ballot selection records
		ballotSelects := make([]model.BallotSelect, len(selectedPollOptions))
		for i, cid := range selectedPollOptions {
			ballotSelects[i] = model.BallotSelect{
				BallotID:     ballot.ID,
				PollOptionID: cid,
			}
		}

		if err := tx.Create(&ballotSelects).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// CheckIfVoterHasVoted checks if voter has already voted
func (b BallotRepository) CheckIfVoterHasVoted(voterId uint64) (bool, error) {
	var count int64
	err := database.SqlSession.Model(&model.Ballot{}).
		Where("invitation_id = ?", voterId).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// DeleteBallot deletes ballot by voter ID
func (b BallotRepository) DeleteBallot(voterID uint64) error {
	err := database.SqlSession.Where("invitation_id = ?", voterID).Delete(&model.Ballot{}).Error
	if err != nil {
		return err
	}

	return nil
}
