package service

import (
	"bytes"
	"context"
	"errors"
	"text/template"
	"time"

	"github.com/faruqputraaa/go-ticket/config"
	"github.com/faruqputraaa/go-ticket/internal/entity"
	"github.com/faruqputraaa/go-ticket/internal/http/dto"
	"github.com/faruqputraaa/go-ticket/internal/repository"
	"github.com/faruqputraaa/go-ticket/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

type UserService interface {
	Login(ctx context.Context, username string, password string) (*entity.JWTCustomClaims, error)
	Register(ctx context.Context, req dto.UserRegisterRequest) error
	GetAll(ctx context.Context) ([]entity.User, error)
	GetByID(ctx context.Context, id int) (*entity.User, error)
	Create(ctx context.Context, req dto.CreateUserRequest) error
	Update(ctx context.Context, req dto.UpdateUserRequest) error
	Delete(ctx context.Context, user *entity.User) error
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error
	RequestResetPassword(ctx context.Context, email string) error
	VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) error
}

type userService struct {
	cfg            *config.Config
	userRepository repository.UserRepository
}

func NewUserService(
	cfg *config.Config,
	userRepository repository.UserRepository,
) UserService {
	return &userService{cfg, userRepository}
}

func (s *userService) Login(ctx context.Context, username string, password string) (*entity.JWTCustomClaims, error) {
	user, err := s.userRepository.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("username atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("username atau password salah")
	}

	expiredTime := time.Now().Add(time.Hour * 10)

	claims := &entity.JWTCustomClaims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-ticket",
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}

	return claims, nil
}

func (s *userService) Register(ctx context.Context, req dto.UserRegisterRequest) error {
	user := new(entity.User)
	user.Username = req.Username
	user.Email = req.Email
	user.Role = "BUYER"
	user.ResetPasswordToken = utils.RandomString(10)
	user.VerifyEmailToken = utils.RandomString(10)
	user.IsVerified = 0

	exist, err := s.userRepository.GetByUsername(ctx, req.Username)
	if err == nil && exist != nil {
		return errors.New("username sudah digunakan")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	err = s.userRepository.Create(ctx, user)
	if err != nil {
		return err
	}
	templatePath := "./templates/email/verify-email.html"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var replacerEmail = struct {
		Token string
	}{
		Token: user.VerifyEmailToken,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, replacerEmail); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.SMTPConfig.Email)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Verify Email Request")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(
		s.cfg.SMTPConfig.Host,
		s.cfg.SMTPConfig.Port,
		s.cfg.SMTPConfig.Email,
		s.cfg.SMTPConfig.Password,
	)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (s *userService) GetAll(ctx context.Context) ([]entity.User, error) {
	return s.userRepository.GetAll(ctx)
}

func (s *userService) GetByID(ctx context.Context, id int) (*entity.User, error) {
	return s.userRepository.GetByID(ctx, id)
}

func (s *userService) Create(ctx context.Context, req dto.CreateUserRequest) error {
	exist, err := s.userRepository.GetByUsername(ctx, req.Username)
	if err == nil && exist != nil {
		return errors.New("username sudah digunakan")
	}

	user := &entity.User{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.userRepository.Create(ctx, user)
}

func (s *userService) Update(ctx context.Context, req dto.UpdateUserRequest) error {
	user, err := s.userRepository.GetByID(ctx, req.IDUser)
	if err != nil {
		return err
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Password != "" {
		user.Password = req.Password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)

		return s.userRepository.Create(ctx, user)
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	return s.userRepository.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, user *entity.User) error {
	return s.userRepository.Delete(ctx, user)
}

func (s *userService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	// check user berdasarkan token reset password
	user, err := s.userRepository.GetByResetPasswordToken(ctx, req.Token)
	if err != nil {
		return errors.New("token reset password salah")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepository.Update(ctx, user)
}

func (s *userService) RequestResetPassword(ctx context.Context, email string) error {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return errors.New("email tidak ditemukan")
	}

	templatePath := "./templates/email/reset-password.html"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var replacerEmail = struct {
		Token string
	}{
		Token: user.ResetPasswordToken,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, replacerEmail); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.SMTPConfig.Email)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Reset Password Request")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(
		s.cfg.SMTPConfig.Host,
		s.cfg.SMTPConfig.Port,
		s.cfg.SMTPConfig.Email,
		s.cfg.SMTPConfig.Password,
	)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (s *userService) VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) error {
	user, err := s.userRepository.GetByVerifyEmailToken(ctx, req.Token)
	if err != nil {
		return errors.New("token verify email salah")
	}

	user.IsVerified = 1
	return s.userRepository.Update(ctx, user)
}
