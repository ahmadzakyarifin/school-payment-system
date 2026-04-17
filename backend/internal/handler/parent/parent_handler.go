package parenthandler

import (
	parentservice "github.com/ahmadzakyarifin/school-payment-system/internal/service/parent"
	"github.com/ahmadzakyarifin/school-payment-system/pkg/response"
	"github.com/gin-gonic/gin"
)

type ParentHandler struct {
	service parentservice.ParentService
}

func New(service parentservice.ParentService) *ParentHandler {
	return &ParentHandler{service: service}
}

func (h *ParentHandler) GetDashboard(c *gin.Context) {
	response.OK(c, "Selamat Datang di Portal Wali Murid", nil)
}
