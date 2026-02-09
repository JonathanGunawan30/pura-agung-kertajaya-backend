package usecase

import (
	"context"
	"errors"
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/model/converter"
	"pura-agung-kertajaya-backend/internal/repository"
	"pura-agung-kertajaya-backend/internal/util"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase interface {
	Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserResponse, string, error)
	Current(ctx context.Context, userID string) (*model.UserResponse, error)
	UpdateProfile(ctx context.Context, userID string, req *model.UpdateUserRequest) (*model.UserResponse, error)
	Logout(ctx context.Context, tokenString string) error
}

type userUseCase struct {
	DB             *gorm.DB
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
	TokenUtil      *util.TokenUtil
	RecaptchaUtil  *util.RecaptchaUtil
}

func NewUserUseCase(
	db *gorm.DB,
	validate *validator.Validate,
	userRepository *repository.UserRepository,
	tokenUtil *util.TokenUtil,
	recaptchaUtil *util.RecaptchaUtil,
) UserUseCase {
	return &userUseCase{
		DB:             db,
		Validate:       validate,
		UserRepository: userRepository,
		TokenUtil:      tokenUtil,
		RecaptchaUtil:  recaptchaUtil,
	}
}

func (c *userUseCase) Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserResponse, string, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(req); err != nil {
		return nil, "", err
	}

	if !c.RecaptchaUtil.Verify(ctx, req.RecaptchaToken) {
		return nil, "", model.ErrForbidden("ReCAPTCHA verification failed")
	}

	var user entity.User
	if err := c.UserRepository.FindByEmail(tx, &user, req.Email); err != nil {
		return nil, "", model.ErrUnauthorized("Invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "", model.ErrUnauthorized("Invalid email or password")
	}

	token, _, err := c.TokenUtil.CreateToken(ctx, &model.Auth{
		ID:    user.ID,
		Role:  user.Role,
		Email: req.Email,
	})
	if err != nil {
		return nil, "", err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, "", err
	}

	return converter.UserToResponse(&user), token, nil
}

func (c *userUseCase) Logout(ctx context.Context, tokenString string) error {
	_, jti, err := c.TokenUtil.ParseToken(ctx, tokenString)
	if err != nil {
		return model.ErrUnauthorized("Invalid or expired token")
	}

	if err := c.TokenUtil.RevokeToken(ctx, jti); err != nil {
		return err
	}

	return nil
}

func (c *userUseCase) Current(ctx context.Context, userID string) (*model.UserResponse, error) {
	var user entity.User
	if err := c.UserRepository.FindById(c.DB, &user, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound("user not found")
		}
		return nil, err
	}
	return converter.UserToResponse(&user), nil
}

func (c *userUseCase) UpdateProfile(ctx context.Context, userID string, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(req); err != nil {
		return nil, err
	}

	var user entity.User
	if err := c.UserRepository.FindById(tx, &user, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound("user not found")
		}
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashed)
	}

	if err := c.UserRepository.Update(tx, &user); err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return converter.UserToResponse(&user), nil
}
