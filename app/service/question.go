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
func (q QuestionService) SelectOneQuestion(id uint64, isAdmin bool, userId uint64) (*model.Question, error) {
	return q.selectQuestion(id, isAdmin, userId, false)
}

// SelectQuestionWithCandidates 檢索問題及其候選人。
func (q QuestionService) SelectQuestionWithCandidates(id uint64, isAdmin bool, userId uint64) (*model.Question, error) {
	return q.selectQuestion(id, isAdmin, userId, true)
}

// selectQuestion 根據提供的 ID 檢查問題是否存在，並根據需要預加載候選人。
func (q QuestionService) selectQuestion(id uint64, isAdmin bool, userId uint64, preloadCandidates bool) (*model.Question, error) {
	questionOne := &model.Question{}
	query := database.SqlSession.
		Where("questions.id = ?", id).
		Joins("JOIN votes ON questions.vote_id = votes.id")

	// 如果需要預加載候選人，則將其添加到查詢中。
	if preloadCandidates {
		query = query.Preload("Candidates")
	}

	// 如果用戶不是管理員，則添加用戶 ID 條件。
	if !isAdmin {
		query = query.Where("votes.user_id = ?", userId)
	}

	err := query.First(&questionOne).Error
	if err != nil {
		return nil, err
	}

	return questionOne, nil
}

// SelectAllQuestions 檢索所有問題。
func (q QuestionService) SelectAllQuestions(voteId uuid.UUID, isAdmin bool, userId uint64) ([]model.Question, error) {
	var questions []model.Question
	query := database.SqlSession.Where("vote_id = ?", voteId)

	if !isAdmin {
		query = query.Joins("JOIN votes ON questions.vote_id = votes.id").Where("votes.user_id = ?", userId)
	}

	err := query.Find(&questions).Error
	if err != nil {
		return nil, err
	}
	return questions, nil
}

// CreateOneQuestion 創建新的問題。
func (q QuestionService) CreateQuestion(form model.QuestionCreate) (model.Question, error) {
	question := model.Question{
		VoteID:      form.VoteID,
		Title:       form.Title,
		Description: form.Description,
	}

	insertErr := database.SqlSession.Model(&model.Question{}).Create(&question).Error
	return question, insertErr
}