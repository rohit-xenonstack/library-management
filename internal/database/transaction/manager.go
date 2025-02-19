package transaction

import (
	"context"
	"sync"

	"gorm.io/gorm"
)

type TxManager struct {
	db *gorm.DB
	mu sync.Mutex
}

func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{
		db: db,
	}
}

func (tm *TxManager) ExecuteInTx(ctx context.Context, fn func(*gorm.DB) error) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tx := tm.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // re-throw panic after rollback
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
