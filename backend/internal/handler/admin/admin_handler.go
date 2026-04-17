package adminhandler

import (
	"strconv"

	"github.com/ahmadzakyarifin/school-payment-system/internal/dto"
	adminservice "github.com/ahmadzakyarifin/school-payment-system/internal/service/admin"
	"github.com/ahmadzakyarifin/school-payment-system/pkg/response"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	service adminservice.AdminService
}

func New(service adminservice.AdminService) *AdminHandler {
	return &AdminHandler{service: service}
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	var req dto.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	users, pagination, err := h.service.GetAllUsers(c.Request.Context(), req)
	if err != nil {
		response.InternalServerError(c, "Gagal mengambil data user")
		return
	}
	response.OKWithPagination(c, "Data user berhasil diambil", users, pagination)
}

func (h *AdminHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), id)
	if err != nil || user == nil {
		response.NotFound(c, "User tidak ditemukan")
		return
	}
	response.OK(c, "Data user berhasil diambil", user)
}

func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, "User berhasil dibuat", user)
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	user, err := h.service.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, "Data user berhasil diupdate", user)
}

func (h *AdminHandler) ToggleStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	user, err := h.service.ToggleUserStatus(c.Request.Context(), id)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, "Status user berhasil diubah", user)
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	if err := h.service.DeleteUser(c.Request.Context(), id); err != nil {
		response.InternalServerError(c, "Gagal menghapus user")
		return
	}
	response.OK(c, "User berhasil dihapus", gin.H{"id": id})
}

func (h *AdminHandler) GetRoles(c *gin.Context) {
	roles, err := h.service.GetRoles(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Gagal mengambil daftar role")
		return
	}
	response.OK(c, "Daftar role berhasil diambil", roles)
}
