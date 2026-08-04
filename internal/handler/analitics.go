package handler

import (
	"SalesTracker/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

func (h *SalesHandler) GetAnalytics(c *ginext.Context) {
	var req model.AnalyticsFilter
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	analytics, err := h.srv.GetAnalytics(c, req)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, analytics)
}

func (h *SalesHandler) ExportCSV(c *ginext.Context) {
	var req model.AnalyticsFilter
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	export, err := h.srv.ExportCSV(c, req)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, export)
}
