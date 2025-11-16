package service

import (
	"strconv"
	"vote/app/model"
	"vote/app/repository"
	"vote/app/utils"

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
	// 查詢資料
	votes, total, err := repository.NewVoteRepository().GetVotes(isAdmin, userId, voteQuery)
	if err != nil {
		return nil, err
	}
	
	paginationRepository := repository.NewPaginationRepository[*model.VoteQuery, model.Vote]()
	votes, hasPreviousPage, hasNextPage := paginationRepository.HasPreviousNextPage(votes, voteQuery)

	var edges []model.VoteEdge
	for _, vote := range votes {
		cursor, _ := (&utils.Password{}).Encrypt(strconv.FormatUint(vote.ID, 10))
		edges = append(edges, model.VoteEdge{
			Node:   vote,
			Cursor: cursor,
		})
	}

	voteConnection := &model.VoteConnection{
		Edges: edges,
		PageInfo: model.PageInfo{
			StartCursor:     edges[0].Cursor,
			EndCursor:       edges[len(edges)-1].Cursor,
			HasNextPage:     hasNextPage,
			HasPreviousPage: hasPreviousPage,
		},
		TotalCount: total,
	}

	var result []*model.VoteConnection
	result = append(result, voteConnection)

	return result, err
}

// CreateOneVote 創建新的投票。
func (v VoteService) CreateVote(form model.VoteCreate) (*model.Vote, error) {
	vote, insertErr := repository.NewVoteRepository().CreateVote(form)
	
	return vote, insertErr
}

// UpdateOneVote 更新投票。
func (v VoteService) UpdateVote(uuid uuid.UUID, form model.VoteUpdate) (*model.Vote, error) {
	// 更新投票並掃描返回的結果
	vote, updateErr := repository.NewVoteRepository().UpdateVote(uuid, form)

	return vote, updateErr
}

// DeleteOneVote 刪除投票。
func (v VoteService) DeleteVote(voteUuids []uuid.UUID, isAdmin bool, userId uint64) ([]*model.Vote, error) {
	votes, err := repository.NewVoteRepository().DeleteVotes(voteUuids, isAdmin, userId)

	return votes, err
}
