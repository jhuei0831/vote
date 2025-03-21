package service

import (
	"vote/app/database"
	"vote/app/model"

	"github.com/google/uuid"
)

type QuestionService struct {
}

func NewQuestionService() QuestionService {
	return QuestionService{}
}

// SelectOneQuestion 根據提供的 ID 檢查問題是否存在。
// TODO: 若非ADMIN，檢查問題的投票場次是該使用者的。
func (q QuestionService) SelectOneQuestion(voteId uuid.UUID, id uint64) (*model.Question, error) {
	questionOne := &model.Question{}
	err := database.SqlSession.Select([]string{"id", "vote_id", "title", "description"}).
		Where("vote_id=?", voteId).
		Where("id=?", id).
		First(&questionOne).Error

	if err != nil {
		return nil, err
	} else {
		return questionOne, nil
	}
}

// SelectAllQuestions 檢索所有問題。
// TODO: 若非ADMIN，檢查問題的投票場次是該使用者的。
func (q QuestionService) SelectAllQuestions(voteId uuid.UUID, isAdmin bool, userId uint64) ([]model.Question, error) {
	var questions []model.Question
	err := database.SqlSession.Where("vote_id = ?", voteId).Find(&questions).Error

	return questions, err
}

// CreateOneQuestion 創建新的問題。
// TODO: 若非ADMIN，檢查問題的投票場次是該使用者的。
func (q QuestionService) CreateQuestion(form model.QuestionCreate) (model.Question, error) {
	question := model.Question{
		VoteID:      form.VoteID,
		Title:       form.Title,
		Description: form.Description,
	}

	insertErr := database.SqlSession.Model(&model.Question{}).Create(&question).Error
	return question, insertErr
}
