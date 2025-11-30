package repository

import (
	"fmt"
	"vote/app/database"
	"vote/app/model"

	"gorm.io/gorm/clause"
)

type CandidateRepository struct {
}

func NewCandidateRepository() CandidateRepository {
	return CandidateRepository{}
}

// GetCandidateByID 根據提供的 ID 檢查候選人是否存在。
func (c CandidateRepository) GetCandidateByID(id uint64, questionId uint64) (*model.Candidate, error) {
	candidate := &model.Candidate{}
	err := database.SqlSession.
		Where("id = ? AND question_id = ?", id, questionId).
		First(&candidate).Error
	if err != nil {
		return nil, err
	}
	return candidate, nil
}

// GetCandidates 根據條件取得所有候選人。
func (c CandidateRepository) GetCandidates(candidateQuery *model.CandidateQuery) ([]model.Candidate, int64, error) {
	var candidates []model.Candidate
	var total int64

	query := database.SqlSession.Model(&model.Candidate{})

	// QuestionID 精確查詢
	if candidateQuery.QuestionID != 0 {
		query = query.Where("question_id = ?", candidateQuery.QuestionID)
	}

	// Name 模糊查詢
	if candidateQuery.Name != "" {
		query = query.Where("name LIKE ?", "%"+candidateQuery.Name+"%")
	}

	// 計算總筆數
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 取得資料
	// 使用分頁服務處理分頁
	paginationRepository := NewPaginationRepository[*model.CandidateQuery, model.Question]()
	query, err = paginationRepository.Handler(query, candidateQuery)
	if err != nil {
		return nil, 0, err
	}

	if err := query.Find(&candidates).Error; err != nil {
		return nil, 0, err
	}

	// 查詢資料
	err = query.Find(&candidates).Error

	return candidates, total, err
}

// CreateCandidate 創建新的候選人。
func (c CandidateRepository) CreateCandidate(userInfo model.UserInfo, form model.CandidateCreate) (*model.Candidate, error) {
	// Verify that the questionID belongs to the voteID
	questionRepository := NewQuestionRepository()
	_, err := questionRepository.GetQuestion(form.QuestionID, userInfo.IsAdmin, userInfo.UserID, false)
	if err != nil {
		return nil, fmt.Errorf("%s", "question record not found")
	}

	candidate := model.Candidate{
		QuestionID: form.QuestionID,
		Name:       form.Name,
	}

	insertErr := database.SqlSession.Model(&model.Candidate{}).Create(&candidate).Error

	return &candidate, insertErr
}

func (c CandidateRepository) UpdateCandidate(userInfo model.UserInfo, id uint64, form model.CandidateUpdate) (*model.Candidate, error) {
	// Verify that the questionID and candidate id belongs to the candidate
	_, err := c.GetCandidateByID(id, form.QuestionID)
	if err != nil {
		return nil, fmt.Errorf("%s", "candidate record not found")
	}

	var candidate model.Candidate

	updateError := database.SqlSession.Model(&candidate).
		Clauses(clause.Returning{}).
		Where("id=?", id).
		Omit("question_id"). // 禁止更新 question_id 字段
		Updates(&form).Error

	return &candidate, updateError
}

// DeleteCandidates 刪除候選人。
func (c CandidateRepository) DeleteCandidates(ids []uint64, userInfo model.UserInfo) ([]*model.Candidate, error) {
	var candidates []*model.Candidate

	query := database.SqlSession.
		Model(&model.Candidate{}).
		Where("candidates.id IN ?", ids)

	// 非管理員需檢查所屬 user
	if !userInfo.IsAdmin {
		query = query.
			Joins("JOIN questions ON candidates.question_id = questions.id").
			Joins("JOIN votes ON questions.vote_id = votes.uuid").
			Where("votes.user_id = ?", userInfo.UserID)
	}

	// 取得候選人並直接刪除
	if err := query.Find(&candidates).Error; err != nil {
		return nil, err
	}

	if len(candidates) == 0 {
		return candidates, nil
	}

	// 提取 IDs 並進行刪除
	candidateIds := make([]uint64, len(candidates))
	for i, candidate := range candidates {
		candidateIds[i] = candidate.ID
	}

	deleteError := database.SqlSession.Delete(&model.Candidate{}, candidateIds).Error
	return candidates, deleteError
}
