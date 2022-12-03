package main

import (
	"context"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type repository struct {
	db *gorm.DB
}

func NewGormRepository(ctx context.Context, connString string) (Repository, error) {
	var r Repository
	{
		db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			log.Fatalf("Failed to open db: %v", err)
		}

		maxOpenConn := 100

		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}
		sqlDB.SetMaxIdleConns(maxOpenConn)
		sqlDB.SetMaxOpenConns(maxOpenConn)
		sqlDB.SetConnMaxLifetime(time.Hour)

		err = db.AutoMigrate(&Category{})
		if err != nil {
			log.Fatalf("Failed to migrate db for type Category: %v", err)
		}

		r = &repository{db: db}

		// Preheat the DB connections by pinging
		// in parallel up to the maxOpenConn count
		for i := 0; i < maxOpenConn; i++ {
			go sqlDB.Ping()
		}
	}

	return r, nil
}

func (p *repository) GetCategories(ctx context.Context) (*[]Category, error) {
	var v []Category
	tx := p.db.WithContext(ctx).Find(&v)
	if tx.Error == gorm.ErrRecordNotFound {
		return nil, gorm.ErrRecordNotFound
	}
	return &v, tx.Error
}
