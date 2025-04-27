package repository

import (
	"context"
	"errors"
	"mnc/internal/model"

	"golang.org/x/crypto/bcrypt"
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

func (r *userRepository) FindByID(ctx context.Context, request model.UserRequest) (model.User, bool, error) {
	var user model.User

	// Find the user by username
	err := r.db.WithContext(ctx).
		Table("user").
		Where("username = ?", request.Username).
		First(&user).Error

	if err != nil {
		// Return false if no user is found
		return model.User{}, false, err
	}

	// Compare the password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		// Return false if password doesn't match
		return model.User{}, false, nil
	}

	return user, true, nil
}

func (r *userRepository) FindUserExist(ctx context.Context, username string) (exists bool, err error) {
	var user model.User
	err = r.db.WithContext(ctx).
		Table("user").
		Where("username = ?", username).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // tidak ada
		}
		return false, err // error lain
	}

	return true, nil // ditemukan
}

func (r *userRepository) ClearUserSession(ctx context.Context, userID int) error {
	return r.db.WithContext(ctx).
		Table("user").
		Where("id = ?", userID).
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
