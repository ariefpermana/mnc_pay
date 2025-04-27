package repository

import (
	"context"
	"mnc/internal/model"

	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func InitPaymentRepo(db *gorm.DB) PaymentRepository {
	_ = db.AutoMigrate(&model.Payment{})
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, req model.PaymentRequest) (model.PaymentResp, error) {
	err := r.db.WithContext(ctx).
		Table("payment").
		Create(&req).Error

	if err != nil {
		return model.PaymentResp{}, err
	}
	// Konversi data yang baru di-insert ke dalam format response
	paymentResp := model.PaymentResp{
		TrxId: req.TrxId, // Ambil ID yang terisi dari database
	}

	// Kembalikan response
	return paymentResp, nil
}
