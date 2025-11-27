package service

import (
	"fmt"
	"vote/app/model"
	"vote/app/repository"
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

// GetQuestions 處理所有問題查詢的共用邏輯,根據 needCandidates 決定是否預載 Candidates。
func (q QuestionService) GetQuestions(questionQuery *model.QuestionQuery) ([]*model.QuestionConnection, error) {
	questions, total, err := repository.NewQuestionRepository().GetQuestions(questionQuery)
	if err != nil {
		return nil, err
	}

	paginationRepository := repository.NewPaginationRepository[*model.QuestionQuery, model.Question]()
	questions, hasPreviousPage, hasNextPage := paginationRepository.HasPreviousNextPage(questions, questionQuery)

	paginationService := NewPaginationService[model.Question, model.QuestionEdge, *model.QuestionConnection]()
	connection := paginationService.BuildConnection(questions, total, hasPreviousPage, hasNextPage,
		func(question model.Question) uint64 {
			return question.ID
		},
	)

	return []*model.QuestionConnection{connection}, nil
}

// CreateOneQuestion 創建新的問題。
func (q QuestionService) CreateQuestion(form model.QuestionCreate) (*model.Question, error) {
	// check vote exists
	_, err := NewVoteService().GetVote(form.VoteID)
	if err != nil {
		return nil, fmt.Errorf("vote not found")
	}

	return repository.NewQuestionRepository().CreateQuestion(form)
}

// UpdateOneQuestion 更新問題。
func (q QuestionService) UpdateQuestion(id uint64, form model.QuestionUpdate) (*model.Question, error) {
	return repository.NewQuestionRepository().UpdateQuestion(id, form)
}

// DeleteOneQuestion 刪除問題。
func (q QuestionService) DeleteQuestion(ids []uint64, isAdmin bool, userId uint64) ([]*model.Question, error) {
	return repository.NewQuestionRepository().DeleteQuestions(ids, isAdmin, userId)
}
