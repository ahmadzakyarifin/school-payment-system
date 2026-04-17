package adminuserservice

import (
	"context"
	"errors"
	"math"

	userdto "github.com/ahmadzakyarifin/school-payment-system/internal/dto/user"
	"github.com/ahmadzakyarifin/school-payment-system/internal/entity"
	adminuserrepo "github.com/ahmadzakyarifin/school-payment-system/internal/repository/admin/user"
	"github.com/ahmadzakyarifin/school-payment-system/pkg/response"
	"github.com/ahmadzakyarifin/school-payment-system/pkg/utils/password"
)

type AdminUserService interface {
	GetAllUsers(ctx context.Context, req userdto.UserListRequest) ([]userdto.UserResponse, response.Pagination, error)
	GetUserByID(ctx context.Context, id int64) (*userdto.UserResponse, error)
	GetRoles(ctx context.Context) ([]string, error)
	CreateUser(ctx context.Context, req userdto.CreateUserRequest) (*userdto.UserResponse, error)
	UpdateUser(ctx context.Context, id int64, req userdto.UpdateUserRequest) (*userdto.UserResponse, error)
	ToggleUserStatus(ctx context.Context, id int64) (*userdto.UserResponse, error)
	DeleteUser(ctx context.Context, id int64) error
}

type adminUserService struct {
	repo adminuserrepo.AdminUserRepository
}

func New(repo adminuserrepo.AdminUserRepository) AdminUserService {
	return &adminUserService{repo: repo}
}

func (s *adminUserService) GetAllUsers(ctx context.Context, req userdto.UserListRequest) ([]userdto.UserResponse, response.Pagination, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	users, total, err := s.repo.FindAll(ctx, req)
	if err != nil {
		return nil, response.Pagination{}, err
	}

	res := []userdto.UserResponse{}
	for _, u := range users {
		res = append(res, s.toResponse(u))
	}

	pg := response.Pagination{
		TotalRows:   total,
		TotalPages:  int(math.Ceil(float64(total) / float64(req.Limit))),
		CurrentPage: req.Page,
		Limit:       req.Limit,
	}
	return res, pg, nil
}

func (s *adminUserService) GetUserByID(ctx context.Context, id int64) (*userdto.UserResponse, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil || u == nil {
		return nil, err
	}
	res := s.toResponse(*u)
	return &res, nil
}

func (s *adminUserService) CreateUser(ctx context.Context, req userdto.CreateUserRequest) (*userdto.UserResponse, error) {
	exists, _ := s.repo.CheckDuplicate(ctx, req.Email, 0)
	if exists {
		return nil, errors.New("email sudah digunakan")
	}
	hashed, _ := password.Hash(req.Password)
	user := &entity.User{
		Name: req.Name, Email: req.Email, PasswordHash: hashed, Phone: req.Phone, Role: req.Role, IsActive: true,
	}
	id, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return s.GetUserByID(ctx, id)
}

func (s *adminUserService) UpdateUser(ctx context.Context, id int64, req userdto.UpdateUserRequest) (*userdto.UserResponse, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil || u == nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if req.Email != nil && *req.Email != u.Email {
		exists, _ := s.repo.CheckDuplicate(ctx, *req.Email, id)
		if exists {
			return nil, errors.New("email sudah digunakan oleh user lain")
		}
		u.Email = *req.Email
	}

	if req.Name != nil {
		u.Name = *req.Name
	}
	if req.Phone != nil {
		u.Phone = req.Phone
	}
	if req.Role != nil {
		u.Role = *req.Role
	}
	if req.IsActive != nil {
		u.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}
	return s.GetUserByID(ctx, id)
}

func (s *adminUserService) ToggleUserStatus(ctx context.Context, id int64) (*userdto.UserResponse, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil || u == nil {
		return nil, errors.New("user tidak ditemukan")
	}
	newStatus := !u.IsActive
	if err := s.repo.UpdateStatus(ctx, id, newStatus); err != nil {
		return nil, err
	}
	return s.GetUserByID(ctx, id)
}

func (s *adminUserService) GetRoles(ctx context.Context) ([]string, error) {
	return s.repo.GetRoles(ctx)
}

func (s *adminUserService) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *adminUserService) toResponse(u entity.User) userdto.UserResponse {
	return userdto.UserResponse{
		ID: u.ID, Name: u.Name, Email: u.Email, Phone: u.Phone, Role: u.Role, IsActive: u.IsActive, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt,
	}
}
