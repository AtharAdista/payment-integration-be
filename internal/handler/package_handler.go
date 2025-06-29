package handler

import (
	"net/http"
	"payment/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PackageHandler struct {
	packageService *service.PackageService
}

func NewPackageHandler(packageService *service.PackageService) *PackageHandler {
	return &PackageHandler{packageService: packageService}
}

func (h *PackageHandler) FindAllPackages(c *gin.Context) {
	packages, err := h.packageService.GetAllPackages()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve packages",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": packages,
	})
}

func (h *PackageHandler) GetPackageById(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid package ID",
		})
		return
	}

	packages, err := h.packageService.GetPackageById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve packages",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": packages,
	})
}
