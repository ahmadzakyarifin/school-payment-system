package authservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/ahmadzakyarifin/school-payment-system/internal/dto"
	authrepo "github.com/ahmadzakyarifin/school-payment-system/internal/repository/auth"
	"github.com/ahmadzakyarifin/school-payment-system/pkg/utils/password"
	"github.com/ahmadzakyarifin/school-payment-system/pkg/utils/token"
)

var (
	ErrInvalidCredentials = errors.New("username atau password salah")
	ErrAccountInactive    = errors.New("akun tidak aktif")
)

// AuthService mendefinisikan kontrak business logic untuk domain auth.
type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
}

type authService struct {
	repo      authrepo.AuthRepository
	jwtSecret string
}

// New membuat instance baru dari authService.
func New(repo authrepo.AuthRepository, jwtSecret string) AuthService {
	return &authService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

// Login memvalidasi kredensial dan mengembalikan JWT token.
func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.FindByUsername(ctx, req.Username)
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

	t, err := token.Generate(user.ID, user.Username, user.Role, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("authservice.Login gagal generate token: %w", err)
	}

	return &dto.LoginResponse{
		Token: t,
		Role:  user.Role,
	}, nil
}
