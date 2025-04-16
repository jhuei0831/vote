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

// SelectQuestions 處理所有問題查詢的共用邏輯，根據 needCandidates 決定是否預載 Candidates。
func (q QuestionService) SelectAllQuestions(voteId uuid.UUID, isAdmin bool, userId uint64, questionQuery model.QuestionQuery) ([]model.Question, int64, error) {
	var questions []model.Question
	var total int64

	query := database.SqlSession.Model(&model.Question{}).Where("vote_id = ?", voteId)

	// 非管理員需檢查所屬 user
	if !isAdmin {
		query = query.Joins("JOIN votes ON questions.vote_id = votes.id").Where("votes.user_id = ?", userId)
	}

	// 標題模糊查詢
	if questionQuery.Title != "" {
		query = query.Where("questions.title LIKE ?", "%"+questionQuery.Title+"%")
	}

	// 計算總筆數
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分頁
	page := questionQuery.Page
	size := questionQuery.Size
	if page > 0 && size > 0 {
		offset := (page - 1) * size
		query = query.Offset(offset).Limit(size)
	}

	// 查詢資料，根據 needCandidates 決定是否預載 Candidates
	if questionQuery.Candidates {
		if err := query.Preload("Candidates").Find(&questions).Error; err != nil {
			return nil, 0, err
		}
	} else {
		if err := query.Find(&questions).Error; err != nil {
			return nil, 0, err
		}
	}
	return questions, total, nil
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