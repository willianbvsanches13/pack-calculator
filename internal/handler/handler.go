package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/willianbsanches13/pack-calculator/internal/calculator"
	"github.com/willianbsanches13/pack-calculator/internal/storage"
)

type Handler struct {
	storage storage.Storage
}

func New(s storage.Storage) *Handler {
	return &Handler{storage: s}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type PackSizesResponse struct {
	PackSizes []int  `json:"pack_sizes"`
	Message   string `json:"message,omitempty"`
}

type CalculateRequest struct {
	Amount    int   `json:"amount" binding:"required,gt=0"`
	PackSizes []int `json:"pack_sizes,omitempty"`
}

type CalculateResponse struct {
	OrderAmount int         `json:"order_amount"`
	TotalItems  int         `json:"total_items"`
	TotalPacks  int         `json:"total_packs"`
	Packs       map[int]int `json:"packs"`
	PackSizes   []int       `json:"pack_sizes_used"`
}

type AddPackSizeRequest struct {
	Size int `json:"size" binding:"required,gt=0"`
}

type RemovePackSizeRequest struct {
	Size int `json:"size" binding:"required,gt=0"`
}

func (h *Handler) GetPackSizes(c *gin.Context) {
	sizes := h.storage.GetPackSizes()
	c.JSON(http.StatusOK, PackSizesResponse{PackSizes: sizes})
}

func (h *Handler) SetPackSizes(c *gin.Context) {
	var req PackSizesResponse
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_json",
			Message: "Failed to parse request body",
		})
		return
	}

	if len(req.PackSizes) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_pack_sizes",
			Message: "Pack sizes cannot be empty",
		})
		return
	}

	for _, size := range req.PackSizes {
		if size <= 0 {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "invalid_pack_size",
				Message: "All pack sizes must be greater than zero",
			})
			return
		}
	}

	h.storage.SetPackSizes(req.PackSizes)
	c.JSON(http.StatusOK, PackSizesResponse{
		PackSizes: h.storage.GetPackSizes(),
		Message:   "Pack sizes updated successfully",
	})
}

func (h *Handler) Calculate(c *gin.Context) {
	var amount int
	var packSizes []int

	if c.Request.Method == http.MethodGet {
		amountQuery := c.Query("amount")
		if amountQuery == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "missing_amount",
				Message: "Amount query parameter is required",
			})
			return
		}

		if n, err := parsePositiveInt(amountQuery); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "invalid_amount",
				Message: "Amount must be a valid positive integer",
			})
			return
		} else {
			amount = n
		}

		packSizes = h.storage.GetPackSizes()
	} else {
		var req CalculateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "invalid_json",
				Message: "Failed to parse request body or amount must be greater than zero",
			})
			return
		}

		amount = req.Amount

		if len(req.PackSizes) > 0 {
			packSizes = req.PackSizes
		} else {
			packSizes = h.storage.GetPackSizes()
		}
	}

	if amount <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_amount",
			Message: "Amount must be greater than zero",
		})
		return
	}

	calc, err := calculator.New(packSizes)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "calculator_error",
			Message: err.Error(),
		})
		return
	}

	result, err := calc.CalculateWithDetails(amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "calculation_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CalculateResponse{
		OrderAmount: result.OrderAmount,
		TotalItems:  result.TotalItems,
		TotalPacks:  result.TotalPacks,
		Packs:       result.Packs,
		PackSizes:   packSizes,
	})
}

func (h *Handler) AddPackSize(c *gin.Context) {
	var req AddPackSizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_json",
			Message: "Failed to parse request body or size must be greater than zero",
		})
		return
	}

	if !h.storage.AddPackSize(req.Size) {
		c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "already_exists",
			Message: "Pack size already exists",
		})
		return
	}

	c.JSON(http.StatusCreated, PackSizesResponse{
		PackSizes: h.storage.GetPackSizes(),
		Message:   "Pack size added successfully",
	})
}

func (h *Handler) RemovePackSize(c *gin.Context) {
	var req RemovePackSizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_json",
			Message: "Failed to parse request body or size must be greater than zero",
		})
		return
	}

	if !h.storage.RemovePackSize(req.Size) {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "not_found",
			Message: "Pack size not found",
		})
		return
	}

	c.JSON(http.StatusOK, PackSizesResponse{
		PackSizes: h.storage.GetPackSizes(),
		Message:   "Pack size removed successfully",
	})
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.Health)

	api := r.Group("/api")
	{
		api.GET("/pack-sizes", h.GetPackSizes)
		api.PUT("/pack-sizes", h.SetPackSizes)
		api.POST("/pack-sizes", h.SetPackSizes)
		api.POST("/pack-sizes/add", h.AddPackSize)
		api.POST("/pack-sizes/remove", h.RemovePackSize)
		api.DELETE("/pack-sizes/remove", h.RemovePackSize)
		api.GET("/calculate", h.Calculate)
		api.POST("/calculate", h.Calculate)
	}
}

func parsePositiveInt(s string) (int, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, calculator.ErrInvalidAmount
		}
		n = n*10 + int(c-'0')
	}
	if n <= 0 {
		return 0, calculator.ErrInvalidAmount
	}
	return n, nil
}
