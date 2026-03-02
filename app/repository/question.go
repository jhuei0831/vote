package repository

import (
	"vote/app/database"
	"vote/app/model"

	"gorm.io/gorm/clause"
	"github.com/google/uuid"
)

type QuestionRepository struct {
}

func NewQuestionRepository() QuestionRepository {
	return QuestionRepository{}
}

// GetQuestionByID 根據提供的 ID 檢查問題是否存在，並根據需要預加載候選人。
func (q QuestionRepository) GetQuestion(id uint64, isAdmin bool, userId uint64, preloadCandidates bool) (*model.Question, error) {
	var question *model.Question

	query := database.SqlSession.Where("questions.id = ?", id)

	// 如果需要預加載候選人，則將其添加到查詢中。
	if preloadCandidates {
		query = query.Preload("Candidates")
	}

	// 如果用戶不是管理員，則添加用戶 ID 條件。
	if !isAdmin {
		query = query.
			Joins("JOIN votes ON questions.vote_id = votes.uuid").
			Where("votes.user_id = ?", userId)
	}

	err := query.First(&question).Error
	if err != nil {
		return nil, err
	}

	return question, nil
}

// GetQuestionsByVoteID 根據提供的 VoteID 獲取相關問題。
func (q QuestionRepository) GetQuestionsByVoteID(voteID uuid.UUID) ([]model.Question, error) {
	var questions []model.Question

	err := database.SqlSession.
		Where("vote_id = ?", voteID).
		Find(&questions).Error
	if err != nil {
		return nil, err
	}

	return questions, nil
}

// GetQuestions 根據條件取得所有問題。
func (q QuestionRepository) GetQuestions(questionQuery *model.QuestionQuery) ([]model.Question, int64, error) {
	var questions []model.Question
	var total int64

	query := database.SqlSession.Model(&model.Question{}).Where("vote_id = ?", questionQuery.VoteID)

	// 標題模糊查詢
	if questionQuery.Title != "" {
		query = query.Where("questions.title LIKE ?", "%"+questionQuery.Title+"%")
	}

	// 計算總筆數
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 使用分頁服務處理分頁
	paginationRepository := NewPaginationRepository[*model.QuestionQuery, model.Question]()
	query, err = paginationRepository.Handler(query, questionQuery)
	if err != nil {
		return nil, 0, err
	}

	if questionQuery.Candidates {
		if err := query.Preload("Candidates").Find(&questions).Error; err != nil {
			return nil, 0, err
		}
	} else {
		if err := query.Find(&questions).Error; err != nil {
			return nil, 0, err
		}
	}

	// 查詢資料
	err = query.Find(&questions).Error

	return questions, total, err
}

// CreateQuestion 創建新的問題。
func (q QuestionRepository) CreateQuestion(form model.QuestionCreate) (*model.Question, error) {
	question := model.Question{
		VoteID:      form.VoteID,
		Title:       form.Title,
		Description: form.Description,
	}

	insertErr := database.SqlSession.Model(&model.Question{}).Create(&question).Error

	return &question, insertErr
}

// UpdateQuestion 更新問題。
func (q QuestionRepository) UpdateQuestion(id uint64, form model.QuestionUpdate) (*model.Question, error) {
	var question model.Question

	updateError := database.SqlSession.
		Model(&question).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Omit("vote_id").
		Updates(&form).Error

	return &question, updateError
}

// DeleteQuestions 刪除問題。
func (q QuestionRepository) DeleteQuestions(ids []uint64, userInfo model.UserInfo) ([]*model.Question, error) {
	var questions []*model.Question

	query := database.SqlSession.
		Model(&model.Question{}).Where("id IN ?", ids)

	// 非管理員需檢查所屬 user
	if !userInfo.IsAdmin {
		query = query.
			Where("vote_id IN (SELECT uuid FROM votes WHERE user_id = ?)", userInfo.UserID)
	}

	// 使用 Returning 子句刪除並返回刪除的記錄
	deleteErr := query.Clauses(clause.Returning{}).Delete(&questions).Error

	return questions, deleteErr
}