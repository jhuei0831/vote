package service

import (
	"vote/app/database"
	"vote/app/model"

	"gorm.io/gorm"
)

type CandidateService struct {
}

func NewCandidateService() CandidateService {
	return CandidateService{}
}

// SelectOneCandidate 根據提供的 ID 檢查候選人是否存在。
func (c CandidateService) SelectOneCandidate(id uint64, isAdmin bool, userId uint64) (*model.Candidate, error) {
	candidateOne := &model.Candidate{}
	query := database.SqlSession.
		Where("candidates.id = ?", id)
	
	if !isAdmin {
		query = query.
			Joins("JOIN questions ON candidates.question_id = questions.id").
			Joins("JOIN votes ON questions.vote_id = votes.id").
			Where("votes.user_id = ?", userId)
	}
		
	err	:= query.First(&candidateOne).Error
	if err != nil {
		return nil, err
	}
	return candidateOne, nil
}

// SelectAllCandidates 檢索所有候選人。
func (c CandidateService) SelectAllCandidates(questionId uint64, isAdmin bool, userId uint64) ([]model.Candidate, error) {
	var candidates []model.Candidate
	query := database.SqlSession.
		Where("candidates.question_id = ?", questionId)
		
	if !isAdmin {
		query = query.
			Joins("JOIN questions ON candidates.question_id = questions.id").
			Joins("JOIN votes ON questions.vote_id = votes.id").
			Where("votes.user_id = ?", userId)
	}

	err := query.Find(&candidates).Error
	
	if err != nil {
		return nil, err
	}

	if len(candidates) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return candidates, nil
}

// CreateOneCandidate 創建新的候選人。
func (c CandidateService) CreateCandidate(form model.CandidateCreate) (model.Candidate, error) {
	candidate := model.Candidate{
		QuestionID: form.QuestionID,
		Name:       form.Name,
	}
	
	insertErr := database.SqlSession.Model(&model.Candidate{}).Create(&candidate).Error
	return candidate, insertErr
}