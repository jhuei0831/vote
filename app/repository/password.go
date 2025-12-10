package repository

import (
	"vote/app/database"
	"vote/app/model"

	"github.com/google/uuid"
)

type PasswordRepository struct {
}

func NewPasswordRepository() PasswordRepository {
	return PasswordRepository{}
}

// GetVoteIdByVoterId 根據投票者ID獲取選票
func (p PasswordRepository) GetVoteIdByVoterId(voterId uint64) (uuid.UUID, error) {
	var voteId uuid.UUID
	err := database.SqlSession.Model(&model.Password{}).
		Select("vore_id").
		Where("id = ?", voterId).
		Scan(&voteId).Error

	if err != nil {
		return uuid.Nil, err
	}

	return voteId, nil
}
