package repository

import (
	"context"
	"mnc/internal/model"

	"gorm.io/gorm"
)

type activityLogRepository struct {
	db *gorm.DB
}

func IniLogRepo(db *gorm.DB) ActivityLogRepository {
	_ = db.AutoMigrate(&model.ActivityLog{})
	return &activityLogRepository{db: db}
}

func (r *activityLogRepository) Create(ctx context.Context, req model.ActivityLog) error {
	err := r.db.WithContext(ctx).
		Table("activity_log").
		Create(&req).Error
	return err
}
