package service

import (
	"vote/app/database"
	"vote/app/model"
	"vote/app/repository"
)

type CandidateService struct {
}

func NewCandidateService() CandidateService {
	return CandidateService{}
}

// GetCandidate 根據提供的 ID 取得候選人。
func (c CandidateService) GetCandidate(id uint64, isAdmin bool, userId uint64) (*model.Candidate, error) {
	candidateOne := &model.Candidate{}
	query := database.SqlSession.
		Where("candidates.id = ?", id)
	
	if !isAdmin {
		query = query.
			Joins("JOIN questions ON candidates.question_id = questions.id").
			Joins("JOIN votes ON questions.vote_id = votes.uuid").
			Where("votes.user_id = ?", userId)
	}
		
	err	:= query.First(&candidateOne).Error
	if err != nil {
		return nil, err
	}
	return candidateOne, nil
}

// SelectAllCandidates 檢索所有候選人。
func (c CandidateService) GetCandidates(candidateQuery *model.CandidateQuery) ([]*model.CandidateConnection, error) {
	candidates, total, err := repository.NewCandidateRepository().GetCandidates(candidateQuery)
	if err != nil {
		return nil, err
	}

	paginationRepository := repository.NewPaginationRepository[*model.CandidateQuery, model.Candidate]()
	candidates, hasPreviousPage, hasNextPage := paginationRepository.HasPreviousNextPage(candidates, candidateQuery)

	paginationService := NewPaginationService[model.Candidate, model.CandidateEdge, *model.CandidateConnection]()
	connection := paginationService.BuildConnection(candidates, total, hasPreviousPage, hasNextPage,
		func(question model.Candidate) uint64 {
			return question.ID
		},
	)

	return []*model.CandidateConnection{connection}, nil
}

// CreateOneCandidate 創建新的候選人。
func (c CandidateService) CreateCandidate(userInfo model.UserInfo, form model.CandidateCreate) (*model.Candidate, error) {
	candidate, err := repository.NewCandidateRepository().CreateCandidate(userInfo, form)
	if err != nil {
		return nil, err
	}

	return candidate, nil
}

func (c CandidateService) UpdateCandidate(userInfo model.UserInfo, id uint64, form model.CandidateUpdate) (*model.Candidate, error) {
	candidate, err := repository.NewCandidateRepository().UpdateCandidate(userInfo, id, form)
	if err != nil {
		return nil, err
	}

	return candidate, nil
}

// DeleteOneCandidate 刪除候選人。
func (c CandidateService) DeleteCandidate(ids []uint64, userInfo model.UserInfo) ([]*model.Candidate, error) {
	return repository.NewCandidateRepository().DeleteCandidates(ids, userInfo)
}