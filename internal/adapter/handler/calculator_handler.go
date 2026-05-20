package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"goapp/internal/domain/entity"
	"goapp/internal/domain/service"
)

type CalculatorHandler struct {
	pricingEngine *service.PricingEngine
}

func NewCalculatorHandler(pricingEngine *service.PricingEngine) *CalculatorHandler {
	return &CalculatorHandler{pricingEngine: pricingEngine}
}

func (h *CalculatorHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/api/calculate", h.Calculate)
	r.POST("/api/calculate/discount", h.CalculateDiscount)
}

func (h *CalculatorHandler) Calculate(c *gin.Context) {
	expr := c.Query("expr")
	if expr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expr parameter required"})
		return
	}

	// Create DTO and call service
	req := &entity.CalculationRequest{
		Expression: expr,
		RuleName:   "default",
	}

	// VULNERABLE: Code Injection (CWE-94)
	result, err := h.pricingEngine.Evaluate(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"expression": expr, "result": result})
}

func (h *CalculatorHandler) CalculateDiscount(c *gin.Context) {
	var body struct {
		Price   float64 `json:"price" form:"price"`
		Formula string  `json:"formula" form:"formula"`
	}
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price and formula required"})
		return
	}

	// Create DTO and call service
	req := &entity.DiscountRequest{
		Price:    body.Price,
		Formula:  body.Formula,
		RuleName: "discount",
	}

	// VULNERABLE: Code Injection (CWE-94)
	result, err := h.pricingEngine.CalculateDiscount(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"price": body.Price, "formula": body.Formula, "result": result})
}
