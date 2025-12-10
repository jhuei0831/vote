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

// GetBallotByVoterId 根據投票者ID獲取選票
func (b BallotService) GetBallotByVoterId(voterId uint64) ([]model.Ballot) {
	repository := repository.NewBallotRepository()
	ballots := repository.GetBallotByVoterId(voterId)

	return ballots
}

// GetBallotByVoteId 根據投票ID獲取選票
func (b BallotService) GetBallotByVoteId(ballotQuery model.BallotQuery) ([]*model.BallotConnection, error) {
	ballots, total, err := repository.NewBallotRepository().GetBallots(ballotQuery)
	if err != nil {
		return nil, err
	}

	paginationRepository := repository.NewPaginationRepository[*model.BallotQuery, model.Ballot]()
	ballotsPage, hasPreviousPage, hasNextPage := paginationRepository.HasPreviousNextPage(ballots, &ballotQuery)

	paginationService := NewPaginationService[model.Ballot, model.BallotEdge, *model.BallotConnection]()
	connection := paginationService.BuildConnection(ballotsPage, total, hasPreviousPage, hasNextPage,
		func(ballot model.Ballot) uint64 {
			return ballot.ID
		},
	)

	return []*model.BallotConnection{connection}, nil
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

// DeleteBallot 刪除選票
func (b BallotService) DeleteBallot(voterID uint64) error {
	repository := repository.NewBallotRepository()
	err := repository.DeleteBallot(voterID)
	if err != nil {
		return err
	}

	return nil
}