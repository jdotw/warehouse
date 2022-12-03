package main

import (
	"context"
	"errors"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type repository struct {
	db *gorm.DB
}

func NewGormRepository(connString string) (Repository, error) {
	var r Repository
	{
		db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		ok(err)
		maxOpenConn := 100

		sqlDB, err := db.DB()
		ok(err)

		sqlDB.SetMaxIdleConns(maxOpenConn)
		sqlDB.SetMaxOpenConns(maxOpenConn)
		sqlDB.SetConnMaxLifetime(time.Hour)

		err = db.AutoMigrate(&Category{})
		ok(err)

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

func (p *repository) CreateCategory(ctx context.Context, category *Category) (*Category, error) {
	return nil, errors.New("not implemented")
}
func (p *repository) DeleteCategory(ctx context.Context, categoryID string) error {
	return errors.New("not implemented")
}
func (p *repository) GetCategory(ctx context.Context, categoryID string) (*Category, error) {
	return nil, errors.New("not implemented")
}
func (p *repository) UpdateCategory(ctx context.Context, category *Category) (*Category, error) {
	return nil, errors.New("not implemented")
}
