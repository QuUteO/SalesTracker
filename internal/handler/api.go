package handler

import (
	"SalesTracker/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
)

func (h *SalesHandler) CreateItem(c *ginext.Context) {
	var req model.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
		return
	}

	item, err := h.srv.Create(c.Request.Context(), &req)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *SalesHandler) ListItems(c *ginext.Context) {
	var filters model.AnalyticsFilter

	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
		return
	}

	items, err := h.srv.List(c, filters)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *SalesHandler) GetItemByID(c *ginext.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	item, err := h.srv.GetByID(c, id)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *SalesHandler) UpdateItem(c *ginext.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	var req model.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
		return
	}

	item, err := h.srv.Update(c, id, &req)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *SalesHandler) DeleteItem(c *ginext.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	if err := h.srv.Delete(c, id); err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": idStr})
}
