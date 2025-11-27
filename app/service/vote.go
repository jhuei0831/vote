package service

import (
	"vote/app/model"
	"vote/app/repository"

	"github.com/google/uuid"
)

type VoteService struct {
}

func NewVoteService() VoteService {
	return VoteService{}
}

// GetVote 根據提供的 ID 檢查投票是否存在。
func (v VoteService) GetVote(uuid uuid.UUID) (*model.Vote, error) {
	vote, err := repository.NewVoteRepository().GetVoteByUUID(uuid)

	if err != nil {
		return nil, err
	} else {
		return vote, nil
	}
}

// GetVotes 檢索所有投票。
func (v VoteService) GetVotes(isAdmin bool, userId uint64, voteQuery *model.VoteQuery) ([]*model.VoteConnection, error) {
	votes, total, err := repository.NewVoteRepository().GetVotes(isAdmin, userId, voteQuery)
	if err != nil {
		return nil, err
	}

	paginationRepository := repository.NewPaginationRepository[*model.VoteQuery, model.Vote]()
	votes, hasPreviousPage, hasNextPage := paginationRepository.HasPreviousNextPage(votes, voteQuery)

	paginationService := NewPaginationService[model.Vote, model.VoteEdge, *model.VoteConnection]()
	connection := paginationService.BuildConnection(votes, total, hasPreviousPage, hasNextPage,
		func(vote model.Vote) uint64 {
			return vote.ID
		},
	)

	return []*model.VoteConnection{connection}, nil
}

// CreateOneVote 創建新的投票。
func (v VoteService) CreateVote(form model.VoteCreate) (*model.Vote, error) {
	return repository.NewVoteRepository().CreateVote(form)
}

// UpdateOneVote 更新投票。
func (v VoteService) UpdateVote(uuid uuid.UUID, form model.VoteUpdate) (*model.Vote, error) {
	// 更新投票並掃描返回的結果
	return repository.NewVoteRepository().UpdateVote(uuid, form)
}

// DeleteOneVote 刪除投票。
func (v VoteService) DeleteVote(voteUuids []uuid.UUID) ([]*model.Vote, error) {
	return repository.NewVoteRepository().DeleteVotes(voteUuids)
}
