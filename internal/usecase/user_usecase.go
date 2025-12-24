package usecase

import (
	"context"
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
	"pura-agung-kertajaya-backend/internal/model/converter"
	"pura-agung-kertajaya-backend/internal/repository"
	"pura-agung-kertajaya-backend/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
	TokenUtil      *util.TokenUtil
	RecaptchaUtil  *util.RecaptchaUtil
}

type UserUseCase interface {
	Login(ctx context.Context, req *model.LoginUserRequest, fiberCtx *fiber.Ctx) (*model.UserResponse, error)
	Current(ctx context.Context, userID int) (*model.UserResponse, error)
	UpdateProfile(ctx context.Context, userID int, req *model.UpdateUserRequest) (*model.UserResponse, error)
	Logout(ctx context.Context, fiberCtx *fiber.Ctx) (bool, error)
}

func NewUserUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate, userRepository *repository.UserRepository, tokenUtil *util.TokenUtil, recaptchaUtil *util.RecaptchaUtil) UserUseCase {
	return &userUseCase{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		UserRepository: userRepository,
		TokenUtil:      tokenUtil,
		RecaptchaUtil:  recaptchaUtil,
	}
}

func (c *userUseCase) Login(ctx context.Context, req *model.LoginUserRequest, fiberCtx *fiber.Ctx) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(req); err != nil {
		return nil, fiber.ErrBadRequest
	}

	if !c.RecaptchaUtil.Verify(ctx, req.RecaptchaToken) {
		c.Log.Warn("reCAPTCHA verification failed")
		return nil, fiber.ErrForbidden
	}

	var user entity.User
	if err := c.UserRepository.FindByEmail(tx, &user, req.Email); err != nil {
		return nil, fiber.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fiber.ErrUnauthorized
	}

	token, jti, err := c.TokenUtil.CreateToken(ctx, &model.Auth{ID: user.ID})
	if err != nil {
		c.Log.Errorf("Failed to create token: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Simpan ke cookie
	fiberCtx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		HTTPOnly: true,
		SameSite: "None",
		Secure:   true,
		Path:     "/",
		Domain:   "admin.puraagungkertajaya.my.id",
		MaxAge:   86400,
	})

	c.Log.Infof("User %s logged in (JTI=%s)", user.Email, jti)

	if err := tx.Commit().Error; err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(&user), nil
}

func (c *userUseCase) Logout(ctx context.Context, fiberCtx *fiber.Ctx) (bool, error) {
	jwtToken := fiberCtx.Cookies("access_token")
	_, jti, err := c.TokenUtil.ParseToken(ctx, jwtToken)
	if err != nil {
		return false, fiber.ErrUnauthorized
	}

	if err := c.TokenUtil.RevokeToken(ctx, jti); err != nil {
		return false, fiber.ErrInternalServerError
	}

	fiberCtx.ClearCookie("access_token")
	return true, nil
}

func (c *userUseCase) Current(ctx context.Context, userID int) (*model.UserResponse, error) {
	var user entity.User
	if err := c.UserRepository.FindById(c.DB, &user, userID); err != nil {
		return nil, fiber.ErrNotFound
	}
	return converter.UserToResponse(&user), nil
}

func (c *userUseCase) UpdateProfile(ctx context.Context, userID int, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(req); err != nil {
		c.Log.Warnf("Invalid update profile request: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	var user entity.User
	if err := c.UserRepository.FindById(tx, &user, userID); err != nil {
		c.Log.Warnf("User not found: %+v", err)
		return nil, fiber.ErrNotFound
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Log.Warnf("Failed to hash password: %+v", err)
			return nil, fiber.ErrInternalServerError
		}
		user.Password = string(hashed)
	}

	if err := c.UserRepository.Update(tx, &user); err != nil {
		c.Log.Warnf("Failed to update user: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(&user), nil
}
