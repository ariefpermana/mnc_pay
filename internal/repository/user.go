package repository

import (
	"context"
	"errors"
	"fmt"
	"mnc/internal/model"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func InitUserRepo(db *gorm.DB) UserRepository {
	_ = db.AutoMigrate(&model.User{})
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, req model.UserRequest) (model.UserCreateResp, error) {
	user := model.User{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
		IsLogin:  0, // default
	}
	// Simpan ke database
	err := r.db.WithContext(ctx).
		Table("user").
		Create(&user).Error

	if err != nil {
		return model.UserCreateResp{}, err
	}

	// Konversi data yang baru di-insert ke dalam format response
	userCreateResp := model.UserCreateResp{
		Id:       user.Id,       // Ambil ID yang terisi dari database
		Username: user.Username, // Ambil Username dari request
	}

	// Kembalikan response
	return userCreateResp, nil
}

func (r *userRepository) FindByID(ctx context.Context, userId int) (model.User, bool, error) {
	var user model.User
	// Find the user by id
	err := r.db.WithContext(ctx).
		Table("user").
		Where("id = ?", userId).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, false, nil
		}
		return model.User{}, false, err
	}

	return user, true, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, request model.UserRequest) (model.User, bool, error) {
	var user model.User

	// Find the user by username
	err := r.db.WithContext(ctx).
		Table("user").
		Where("username = ?", request.Username).
		First(&user).Error

	// Check if the error is related to a non-existing user
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return false and nil if user is not found
			return model.User{}, false, nil
		}
		// Return the error if it's something other than 'not found'
		return model.User{}, false, err
	}

	// If no error, return the user and true
	return user, true, nil
}

func (r *userRepository) FindUserExist(ctx context.Context, username string) (exists bool, err error) {
	var user model.User
	err = r.db.WithContext(ctx).
		Table("user").
		Where("username = ?", username).
		First(&user).Error

	if err != nil {
		fmt.Println("DEBUG: DB Error =", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("DEBUG: Record not found")
			return false, nil // tidak ada
		}
		fmt.Println("DEBUG: Other error")

		return false, err // error lain
	}

	return true, nil // ditemukan
}

func (r *userRepository) ClearUserSession(ctx context.Context, username string) error {
	return r.db.WithContext(ctx).
		Table("user").
		Where("username = ?", username).
		Updates(map[string]interface{}{
			"token":    nil,
			"is_login": 0,
		}).Error
}

func (r *userRepository) Login(ctx context.Context, userID int, token string) error {
	err := r.db.WithContext(ctx).Table("user").Where("id = ?", userID).Updates(map[string]interface{}{
		"token":    token,
		"is_login": 1,
	}).Error
	if err != nil {
		return err
	}

	return err
}
