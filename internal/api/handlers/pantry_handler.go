package handlers

import (
	"net/http"
	"strconv"

	"github.com/byte4bite/byte4bite/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PantryHandler handles HTTP requests for pantries
type PantryHandler struct {
	pantryService *services.PantryService
}

// NewPantryHandler creates a new pantry handler
func NewPantryHandler(pantryService *services.PantryService) *PantryHandler {
	return &PantryHandler{
		pantryService: pantryService,
	}
}

// GetPantries returns a list of pantries
// @Summary Get pantries
// @Description Get list of pantries (optionally filter by active status)
// @Tags pantries
// @Produce json
// @Param is_active query boolean false "Filter by active status"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} services.GetPantriesResponse
// @Router /api/v1/pantries [get]
func (h *PantryHandler) GetPantries(c *gin.Context) {
	var req services.GetPantriesRequest

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		req.IsActive = &isActive
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	req.Page = page
	req.PageSize = pageSize

	response, err := h.pantryService.GetPantries(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetPantry returns a single pantry by ID
// @Summary Get pantry by ID
// @Description Get detailed information about a specific pantry
// @Tags pantries
// @Produce json
// @Param id path string true "Pantry ID"
// @Success 200 {object} models.Pantry
// @Router /api/v1/pantries/{id} [get]
func (h *PantryHandler) GetPantry(c *gin.Context) {
	pantryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pantry ID"})
		return
	}

	pantry, err := h.pantryService.GetPantry(pantryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pantry)
}

// SearchPantries searches for pantries by query
// @Summary Search pantries
// @Description Search pantries by name or city
// @Tags pantries
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {array} models.Pantry
// @Router /api/v1/pantries/search [get]
func (h *PantryHandler) SearchPantries(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search query is required"})
		return
	}

	pantries, err := h.pantryService.SearchPantries(query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pantries)
}

// GetPantriesByCity gets pantries in a specific city
// @Summary Get pantries by city
// @Description Get pantries in a specific city
// @Tags pantries
// @Produce json
// @Param city query string true "City name"
// @Success 200 {array} models.Pantry
// @Router /api/v1/pantries/by-city [get]
func (h *PantryHandler) GetPantriesByCity(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city is required"})
		return
	}

	pantries, err := h.pantryService.GetPantriesByCity(city)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pantries)
}

// GetPantriesByZipCode gets pantries in a specific zip code
// @Summary Get pantries by zip code
// @Description Get pantries in a specific zip code
// @Tags pantries
// @Produce json
// @Param zip query string true "Zip code"
// @Success 200 {array} models.Pantry
// @Router /api/v1/pantries/by-zip [get]
func (h *PantryHandler) GetPantriesByZipCode(c *gin.Context) {
	zipCode := c.Query("zip")
	if zipCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "zip code is required"})
		return
	}

	pantries, err := h.pantryService.GetPantriesByZipCode(zipCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pantries)
}

// CreatePantry creates a new pantry (admin only)
// @Summary Create pantry
// @Description Create a new pantry
// @Tags pantries
// @Accept json
// @Produce json
// @Param body body services.CreatePantryRequest true "Pantry data"
// @Success 201 {object} models.Pantry
// @Router /api/v1/admin/pantries [post]
func (h *PantryHandler) CreatePantry(c *gin.Context) {
	var req services.CreatePantryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pantry, err := h.pantryService.CreatePantry(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pantry)
}

// UpdatePantry updates a pantry (admin only)
// @Summary Update pantry
// @Description Update pantry information
// @Tags pantries
// @Accept json
// @Produce json
// @Param id path string true "Pantry ID"
// @Param body body services.UpdatePantryRequest true "Pantry update data"
// @Success 200 {object} models.Pantry
// @Router /api/v1/admin/pantries/{id} [put]
func (h *PantryHandler) UpdatePantry(c *gin.Context) {
	pantryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pantry ID"})
		return
	}

	var req services.UpdatePantryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pantry, err := h.pantryService.UpdatePantry(pantryID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pantry)
}

// DeletePantry deletes a pantry (admin only)
// @Summary Delete pantry
// @Description Delete a pantry
// @Tags pantries
// @Produce json
// @Param id path string true "Pantry ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/admin/pantries/{id} [delete]
func (h *PantryHandler) DeletePantry(c *gin.Context) {
	pantryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pantry ID"})
		return
	}

	if err := h.pantryService.DeletePantry(pantryID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "pantry deleted successfully"})
}

// TogglePantryStatus toggles the active status of a pantry (admin only)
// @Summary Toggle pantry status
// @Description Toggle the active status of a pantry
// @Tags pantries
// @Produce json
// @Param id path string true "Pantry ID"
// @Success 200 {object} models.Pantry
// @Router /api/v1/admin/pantries/{id}/toggle [patch]
func (h *PantryHandler) TogglePantryStatus(c *gin.Context) {
	pantryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pantry ID"})
		return
	}

	pantry, err := h.pantryService.TogglePantryStatus(pantryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pantry)
}
