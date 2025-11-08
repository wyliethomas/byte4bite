package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/byte4bite/byte4bite/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DonationHandler handles HTTP requests for donations
type DonationHandler struct {
	donationService *services.DonationService
}

// NewDonationHandler creates a new donation handler
func NewDonationHandler(donationService *services.DonationService) *DonationHandler {
	return &DonationHandler{
		donationService: donationService,
	}
}

// CreateDonation creates a new donation (public endpoint)
// @Summary Create donation
// @Description Submit a new donation to a pantry
// @Tags donations
// @Accept json
// @Produce json
// @Param body body services.CreateDonationRequest true "Donation data"
// @Success 201 {object} models.Donation
// @Router /api/v1/donations [post]
func (h *DonationHandler) CreateDonation(c *gin.Context) {
	var req services.CreateDonationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	donation, err := h.donationService.CreateDonation(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, donation)
}

// GetDonations returns a list of donations (admin only)
// @Summary Get donations
// @Description Get list of donations with optional filters
// @Tags donations
// @Produce json
// @Param pantry_id query string false "Filter by pantry ID"
// @Param receipt_sent query boolean false "Filter by receipt sent status"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} services.GetDonationsResponse
// @Router /api/v1/admin/donations [get]
func (h *DonationHandler) GetDonations(c *gin.Context) {
	var req services.GetDonationsRequest

	// Parse pantry_id filter
	if pantryIDStr := c.Query("pantry_id"); pantryIDStr != "" {
		pantryID, err := uuid.Parse(pantryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pantry ID"})
			return
		}
		req.PantryID = &pantryID
	}

	// Parse receipt_sent filter
	if receiptSentStr := c.Query("receipt_sent"); receiptSentStr != "" {
		receiptSent := receiptSentStr == "true"
		req.ReceiptSent = &receiptSent
	}

	// Parse date filters
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date format (use YYYY-MM-DD)"})
			return
		}
		req.StartDate = &startDate
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date format (use YYYY-MM-DD)"})
			return
		}
		req.EndDate = &endDate
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	req.Page = page
	req.PageSize = pageSize

	response, err := h.donationService.GetDonations(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetDonation returns a single donation by ID (admin only)
// @Summary Get donation by ID
// @Description Get detailed information about a specific donation
// @Tags donations
// @Produce json
// @Param id path string true "Donation ID"
// @Success 200 {object} models.Donation
// @Router /api/v1/admin/donations/{id} [get]
func (h *DonationHandler) GetDonation(c *gin.Context) {
	donationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid donation ID"})
		return
	}

	donation, err := h.donationService.GetDonation(donationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, donation)
}

// UpdateDonation updates a donation (admin only)
// @Summary Update donation
// @Description Update donation information
// @Tags donations
// @Accept json
// @Produce json
// @Param id path string true "Donation ID"
// @Param body body services.UpdateDonationRequest true "Donation update data"
// @Success 200 {object} models.Donation
// @Router /api/v1/admin/donations/{id} [put]
func (h *DonationHandler) UpdateDonation(c *gin.Context) {
	donationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid donation ID"})
		return
	}

	var req services.UpdateDonationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	donation, err := h.donationService.UpdateDonation(donationID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, donation)
}

// DeleteDonation deletes a donation (admin only)
// @Summary Delete donation
// @Description Delete a donation record
// @Tags donations
// @Produce json
// @Param id path string true "Donation ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/admin/donations/{id} [delete]
func (h *DonationHandler) DeleteDonation(c *gin.Context) {
	donationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid donation ID"})
		return
	}

	if err := h.donationService.DeleteDonation(donationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "donation deleted successfully"})
}

// MarkReceiptSent marks a donation receipt as sent (admin only)
// @Summary Mark receipt as sent
// @Description Mark a donation receipt as sent
// @Tags donations
// @Produce json
// @Param id path string true "Donation ID"
// @Success 200 {object} models.Donation
// @Router /api/v1/admin/donations/{id}/receipt [patch]
func (h *DonationHandler) MarkReceiptSent(c *gin.Context) {
	donationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid donation ID"})
		return
	}

	donation, err := h.donationService.MarkReceiptSent(donationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, donation)
}

// SearchDonations searches donations (admin only)
// @Summary Search donations
// @Description Search donations by donor name or description
// @Tags donations
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} services.GetDonationsResponse
// @Router /api/v1/admin/donations/search [get]
func (h *DonationHandler) SearchDonations(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search query is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	response, err := h.donationService.SearchDonations(query, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetDonationStats returns donation statistics (admin only)
// @Summary Get donation statistics
// @Description Get statistics about donations
// @Tags donations
// @Produce json
// @Param pantry_id query string false "Filter by pantry ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} services.DonationStatsResponse
// @Router /api/v1/admin/donations/stats [get]
func (h *DonationHandler) GetDonationStats(c *gin.Context) {
	var pantryID *uuid.UUID
	var startDate, endDate *time.Time

	// Parse pantry_id filter
	if pantryIDStr := c.Query("pantry_id"); pantryIDStr != "" {
		id, err := uuid.Parse(pantryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pantry ID"})
			return
		}
		pantryID = &id
	}

	// Parse date filters
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		date, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date format (use YYYY-MM-DD)"})
			return
		}
		startDate = &date
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		date, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date format (use YYYY-MM-DD)"})
			return
		}
		endDate = &date
	}

	stats, err := h.donationService.GetDonationStats(pantryID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetDonationsByDonor gets donations by donor email (admin only)
// @Summary Get donations by donor
// @Description Get all donations from a specific donor by email
// @Tags donations
// @Produce json
// @Param email query string true "Donor email"
// @Success 200 {array} models.Donation
// @Router /api/v1/admin/donations/by-donor [get]
func (h *DonationHandler) GetDonationsByDonor(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	donations, err := h.donationService.GetDonationsByDonor(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, donations)
}
