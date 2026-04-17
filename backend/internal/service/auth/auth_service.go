package authservice

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	authdto "github.com/ahmadzakyarifin/school-payment-system/internal/dto/auth"
	userdto "github.com/ahmadzakyarifin/school-payment-system/internal/dto/user"
	authrepo "github.com/ahmadzakyarifin/school-payment-system/internal/repository/auth"
	"github.com/ahmadzakyarifin/school-payment-system/pkg/utils/mailer"
	"github.com/ahmadzakyarifin/school-payment-system/pkg/utils/password"
	"github.com/ahmadzakyarifin/school-payment-system/pkg/utils/token"
)

var (
	ErrInvalidCredentials = errors.New("email atau password salah")
	ErrAccountInactive    = errors.New("akun tidak aktif")
	ErrEmailNotFound      = errors.New("email tidak terdaftar")
	ErrInvalidToken       = errors.New("token tidak valid atau sudah kedaluwarsa")
)

type AuthService interface {
	Login(ctx context.Context, req authdto.LoginRequest) (*authdto.LoginResponse, error)
	ForgotPassword(ctx context.Context, req authdto.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req authdto.ResetPasswordRequest) error
}

type authService struct {
	repo      authrepo.AuthRepository
	mailer    mailer.Mailer
	jwtSecret string
}

func New(repo authrepo.AuthRepository, mailer mailer.Mailer, jwtSecret string) AuthService {
	return &authService{
		repo:      repo,
		mailer:    mailer,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Login(ctx context.Context, req authdto.LoginRequest) (*authdto.LoginResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("authservice.Login: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}
	if !user.IsActive {
		return nil, ErrAccountInactive
	}
	if !password.Verify(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	t, err := token.Generate(user.ID, user.Email, user.Role, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("authservice.Login gagal generate token: %w", err)
	}

	return &authdto.LoginResponse{
		Token: t,
		User: userdto.UserResponse{
			ID: user.ID, Name: user.Name, Email: user.Email, Role: user.Role, IsActive: user.IsActive, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

func (s *authService) ForgotPassword(ctx context.Context, req authdto.ForgotPasswordRequest) error {
	user, _ := s.repo.FindByEmail(ctx, req.Email)
	if user == nil {
		fmt.Printf("[Forgot Password] Email tidak ditemukan: %s\n", req.Email)
		return nil
	}

	// Generate Random Token
	b := make([]byte, 32)
	rand.Read(b)
	token := hex.EncodeToString(b)

	// Simpan ke DB
	if err := s.repo.SaveResetToken(ctx, user.Email, token); err != nil {
		fmt.Printf("[Forgot Password] ERROR simpan token ke DB: %v\n", err)
		return err
	}

	fmt.Printf("[Forgot Password] Token berhasil dibuat untuk %s. Mengirim email...\n", user.Email)

	// Kirim Email via Resend
	go func() {
		err := s.mailer.SendResetPassword(user.Email, token)
		if err != nil {
			fmt.Printf("[Forgot Password] ERROR kirim email via Resend: %v\n", err)
		} else {
			fmt.Printf("[Forgot Password] Email berhasil terkirim ke %s\n", user.Email)
		}
	}()

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, req authdto.ResetPasswordRequest) error {
	email, err := s.repo.GetResetToken(ctx, req.Token)
	if err != nil || email == "" {
		return ErrInvalidToken
	}

	hashed, _ := password.Hash(req.Password)
	if err := s.repo.UpdatePassword(ctx, email, hashed); err != nil {
		return err
	}

	// Hapus token setelah dipakai
	s.repo.DeleteResetToken(ctx, req.Token)

	return nil
}
