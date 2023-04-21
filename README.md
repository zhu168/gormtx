# gormTx
gorm transaction simplification tool,simplify transaction rollback handling。

## How to use it

```go
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
```