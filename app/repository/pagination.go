package repository

import (
	"fmt"
	"vote/app/utils"

	"gorm.io/gorm"
)

// PaginationQuery 定義分頁查詢所需的方法
type PaginationQuery interface {
	GetFirst() int
	GetAfter() string
	GetLast() int
	GetBefore() string
}

type PaginationRepository[T PaginationQuery, S any] struct {
}

func NewPaginationRepository[T PaginationQuery, S any]() PaginationRepository[T, S] {
	return PaginationRepository[T, S]{}
}

// Handler 處理分頁邏輯，根據提供的分頁參數修改 GORM 查詢。
func (p *PaginationRepository[T, S]) Handler(db *gorm.DB, query T) (*gorm.DB, error) {
	// 檢查分頁參數
	err := p.CheckPaginationParams(query)
	if err != nil {
		return nil, err
	}

	// 處理 Forward Pagination
	if query.GetFirst() > 0 && query.GetAfter() == "" && query.GetBefore() == "" {
		db = db.Order("created_at DESC").Limit(query.GetFirst() + 1)
	}

	if query.GetFirst() > 0 && query.GetAfter() != "" {
		after, _ := (&utils.Password{}).Decrypt(query.GetAfter())
		db = db.Where("id > ?", after)
		db = db.Limit(query.GetFirst() + 1)
	}

	// 處理 Backward Pagination
	if query.GetLast() > 0 {
		if query.GetBefore() != "" {
			before, _ := (&utils.Password{}).Decrypt(query.GetBefore())
			db = db.Where("id < ?", before)
		}
		db = db.Order("created_at asc").Limit(query.GetLast() + 1)
	}

	return db, nil
}

// HasPreviousNextPage 判斷是否有上一頁或下一頁
func (p *PaginationRepository[T, S]) HasPreviousNextPage(items []S, query T) ([]S, bool, bool) {
	hasPreviousPage := false
	hasNextPage := false

	if query.GetFirst() > 0 && query.GetAfter() == "" && query.GetBefore() == "" {
		hasPreviousPage = false
		hasNextPage = len(items) > query.GetFirst()
		if hasNextPage {
			items = items[:len(items)-1]
		}
	}

	if query.GetFirst() > 0 && query.GetAfter() != "" {
		hasPreviousPage = true
		hasNextPage = len(items) > query.GetFirst()
		if hasNextPage {
			items = items[:len(items)-1]
		}
	}

	if query.GetLast() > 0 && query.GetAfter() == "" && query.GetBefore() == "" {
		hasNextPage = false
		hasPreviousPage = len(items) > query.GetLast()
		if hasPreviousPage {
			items = items[1:]
		}
	}

	if query.GetLast() > 0 && query.GetBefore() != "" {
		hasNextPage = true
		hasPreviousPage = len(items) > query.GetLast()
		if hasPreviousPage {
			items = items[1:]
		}
	}

	return items, hasPreviousPage, hasNextPage
}

// CheckPaginationParams 檢查分頁參數的有效性
func (p *PaginationRepository[T, S]) CheckPaginationParams(query T) error {
	if query.GetFirst() <= 0 && query.GetLast() <= 0 {
		return fmt.Errorf("must provide either first or last parameter")
	}

	if query.GetFirst() > 0 && query.GetLast() > 0 {
		return fmt.Errorf("cannot provide both first and last parameters")
	}

	if query.GetAfter() != "" && query.GetBefore() != "" {
		return fmt.Errorf("cannot provide both after and before parameters")
	}
	return nil
}
