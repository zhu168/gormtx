package gormtx

import (
	"gorm.io/gorm"
)

type GORMTX struct {
	DB           *gorm.DB
	TX           *gorm.DB
	Error        error
	AutoRollback bool
	Rollbacked   bool
}

// New gormtx
func New(db *gorm.DB) (t *GORMTX) {
	t = &GORMTX{}
	t.AutoRollback = true
	t.DB = db
	return
}

// Begin transaction
func (my *GORMTX) Begin() error {
	my.TX = my.DB.Begin()
	if my.Error != nil {
		my.Error = my.TX.Error
		if my.AutoRollback {
			my.Rollbacked = true
		}
		return my.Error
	}
	return nil
}

// Execute transaction code, if the transaction has been rolled back, f will not be executed
func (my *GORMTX) Exec(f func() error) error {
	if my.AutoRollback && my.Rollbacked {
		return nil
	}
	my.Error = f()
	if my.Error != nil {
		if my.AutoRollback {
			my.TX.Rollback()
			my.Rollbacked = true
		}
		return my.Error
	}
	return nil
}

// Commit the transaction, if the transaction has rolled, the commit will not be executed
func (my *GORMTX) Commit() error {
	if my.Rollbacked {
		return nil
	}
	my.TX.Commit()
	if my.Error != nil {
		my.Error = my.TX.Error
		my.Rollback()
		return my.Error
	}
	return nil
}

// Rollback the transaction, if the transaction has been rolled, Rollback will not be executed
func (my *GORMTX) Rollback() error {
	if my.Rollbacked {
		return nil
	}
	my.TX.Rollback()
	if my.Error != nil {
		my.Error = my.TX.Error
		my.Rollbacked = true
		return my.Error
	}
	return nil
}
