package service

import (
	"strconv"
	"vote/app/database"
	"vote/app/model"
	"vote/app/utils"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type VoteService struct {
}

func NewVoteService() VoteService {
	return VoteService{}
}

// SelectOneVote 根據提供的 ID 檢查投票是否存在。
func (v VoteService) SelectOneVote(id uuid.UUID) (*model.Vote, error) {
	voteOne := &model.Vote{}
	err := database.SqlSession.
		Select([]string{"id", "title", "description", "user_id", "start_time", "end_time"}).
		Where("id=?", id).
		First(&voteOne).Error
	if err != nil {
		return nil, err
	} else {
		return voteOne, nil
	}
}

// SelectAllVotes 檢索所有投票。
func (v VoteService) SelectAllVotes(isAdmin bool, userId uint64, voteQuery *model.VoteQuery) ([]*model.VoteConnection, error) {
	var votes []model.Vote
	var total int64
	query := database.SqlSession.Model(&model.Vote{})

	if !isAdmin {
		query = query.Where("user_id = ?", userId)
	}

	// 計算總筆數
	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 使用分頁服務處理分頁
	paginationService := NewPaginationService[*model.VoteQuery, model.Vote]()
	query, err = paginationService.Handler(query, voteQuery)
	if err != nil {
		return nil, err
	}

	// 查詢資料
	err = query.Find(&votes).Error
	votes, hasPreviousPage, hasNextPage := paginationService.HasPreviousNextPage(votes, voteQuery)

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
	vote := model.Vote{
		Title:       form.Title,
		Description: form.Description,
		UserID:      form.UserID,
		StartTime:   form.StartTime,
		EndTime:     form.EndTime,
	}

	insertErr := database.SqlSession.Model(&model.Vote{}).Create(&vote).Error
	return &vote, insertErr
}

// UpdateOneVote 更新投票。
func (v VoteService) UpdateVote(uuid uuid.UUID, form model.VoteUpdate) (*model.Vote, error) {
    var vote model.Vote
    // 更新投票並掃描返回的結果
    update := database.SqlSession.Model(&vote).
        Clauses(clause.Returning{}).
        Where("uuid=?", uuid).
        Updates(&form)

    return &vote, update.Error
}

// DeleteOneVote 刪除投票。
func (v VoteService) DeleteVote(voteUuids []uuid.UUID, isAdmin bool, userId uint64) ([]*model.Vote, error) {
	var votes []*model.Vote
	query := database.SqlSession.
		Clauses(clause.Returning{}).
		Where("uuid IN (?)", voteUuids)

	if !isAdmin {
		query = query.Where("user_id = ?", userId)
	}

	err := query.Delete(&votes).Error
	return votes, err
}
