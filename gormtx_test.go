package gormtx_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/zhu168/gormtx"
	"gorm.io/gorm"

	"gorm.io/driver/sqlite"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func init() {
	_ = os.Remove("test.db")
}

func TestTrueTx(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Product{})

	gormtx := gormtx.New(db)
	gormtx.Begin()
	gormtx.Exec(func() error { return gormtx.TX.Create(&Product{Code: "001", Price: 100}).Error })
	gormtx.Exec(func() error { return gormtx.TX.Create(&Product{Code: "002", Price: 100}).Error })
	gormtx.Commit()
	products := []Product{}
	db.Find(&products)
	if len(products) != 2 {
		t.Errorf("test fail %d", len(products))
	}
	gormtx.Begin()
	gormtx.Exec(func() error { return gormtx.TX.Create(&Product{Code: "003", Price: 100}).Error })
	gormtx.Exec(func() error { return gormtx.TX.Create(&Product{Code: "004", Price: 100}).Error })
	gormtx.Rollback()
	gormtx.Commit()
	if len(products) != 2 {
		t.Errorf("test fail %d", len(products))
	}
}
func TestFailTx(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Product{})
	gormtx := gormtx.New(db)
	gormtx.Begin()
	gormtx.Exec(func() error { return gormtx.TX.Create(&Product{Code: "101", Price: 100}).Error })
	gormtx.Exec(func() error { return fmt.Errorf("test error") })
	gormtx.Exec(func() error { return gormtx.TX.Create(&Product{Code: "102", Price: 100}).Error })
	gormtx.Commit()
	products := []Product{}
	db.Where("Code like '10%'").Find(&products)
	if len(products) != 0 {
		t.Errorf("test fail %d", len(products))
	}
}
