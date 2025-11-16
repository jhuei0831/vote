package service

import (
	"fmt"
	"strconv"
	"vote/app/database"
	"vote/app/model"
	"vote/app/repository"
	"vote/app/utils"
)

type QuestionService struct {
}

func NewQuestionService() QuestionService {
	return QuestionService{}
}

// GetQuestion 根據提供的 ID 檢查問題是否存在。
func (q QuestionService) GetQuestion(id uint64, isAdmin bool, userId uint64) (*model.Question, error) {
	return repository.NewQuestionRepository().GetQuestion(id, isAdmin, userId, false)
}

// SelectQuestionWithCandidates 檢索問題及其候選人。
func (q QuestionService) SelectQuestionWithCandidates(id uint64, isAdmin bool, userId uint64) (*model.Question, error) {
	return repository.NewQuestionRepository().GetQuestion(id, isAdmin, userId, true)
}

// SelectQuestions 處理所有問題查詢的共用邏輯，根據 needCandidates 決定是否預載 Candidates。
func (q QuestionService) GetQuestions(isAdmin bool, userId uint64, questionQuery *model.QuestionQuery) ([]*model.QuestionConnection, error) {
	questions, total, err := repository.NewQuestionRepository().GetQuestions(isAdmin, userId, questionQuery)
	if err != nil {
		return nil, err
	}

	paginationRepository := repository.NewPaginationRepository[*model.QuestionQuery, model.Question]()
	questions, hasPreviousPage, hasNextPage := paginationRepository.HasPreviousNextPage(questions, questionQuery)

	var edges []model.QuestionEdge
	for _, question := range questions {
		cursor, _ := (&utils.Password{}).Encrypt(strconv.FormatUint(question.ID, 10))
		edges = append(edges, model.QuestionEdge{
			Node:   question,
			Cursor: cursor,
		})
	}

	questionConnection := &model.QuestionConnection{
		Edges: edges,
		PageInfo: model.PageInfo{
			StartCursor:     edges[0].Cursor,
			EndCursor:       edges[len(edges)-1].Cursor,
			HasNextPage:     hasNextPage,
			HasPreviousPage: hasPreviousPage,
		},
		TotalCount: total,
	}

	var result []*model.QuestionConnection
	result = append(result, questionConnection)

	return result, err
}

// CreateOneQuestion 創建新的問題。
func (q QuestionService) CreateQuestion(form model.QuestionCreate) (*model.Question, error) {
	// check vote exists
	_, err := NewVoteService().GetVote(form.VoteID)
	if err != nil {
		return nil, fmt.Errorf("vote not found")
	}

	question := model.Question{
		VoteID:      form.VoteID,
		Title:       form.Title,
		Description: form.Description,
	}

	insertErr := database.SqlSession.Model(&model.Question{}).Create(&question).Error
	return &question, insertErr
}
