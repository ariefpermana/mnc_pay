package model

import "time"

type PaymentRequest struct {
	UserId      int    `json:"user_id"`
	AccountNo   string `json:"account_no"`
	AccountName string `json:"account_name"`
	Amount      string `json:"amount"`
	Merchant    string `json:"merchant"`
	TrxId       string `json:"trx_id"`
	Status      string `json:"status"`
}

type PaymentResp struct {
	TrxId string
}

type Payment struct {
	Id          int        `gorm:"primaryKey;column:id;type:integer"`
	UserId      int        `gorm:"column:user_id;type:integer"`
	AccountNo   string     `gorm:"column:account_no;type:varchar(100)"`
	AccountName string     `gorm:"column:account_name;type:varchar(255)"`
	Amount      float64    `gorm:"column:amount;type:double"`
	Merchant    string     `gorm:"column:merchant;type:varchar(255)"`
	TrxId       string     `gorm:"column:trx_id;type:varchar(25)"`
	Status      string     `gorm:"column:status;type:varchar(25)"`
	CreatedAt   *time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt   *time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (Payment) TableName() string {
	return "payment"
}
