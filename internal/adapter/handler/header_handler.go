package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"goapp/internal/domain/entity"
	"goapp/internal/domain/service"
)

type HeaderHandler struct {
	customizationService *service.ResponseCustomizationService
	localizationService  *service.LocalizationService
}

func NewHeaderHandler(customizationService *service.ResponseCustomizationService, localizationService *service.LocalizationService) *HeaderHandler {
	return &HeaderHandler{
		customizationService: customizationService,
		localizationService:  localizationService,
	}
}

func (h *HeaderHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/api/header/set", h.SetHeader)
	r.GET("/api/header/redirect", h.RedirectWithLang)
}

func (h *HeaderHandler) SetHeader(c *gin.Context) {
	name := c.Query("name")
	value := c.Query("value")

	if name == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and value parameters required"})
		return
	}

	// Create DTO and call service
	req := &entity.HeaderRequest{
		HeaderName:  name,
		HeaderValue: value,
		PolicyName:  "custom",
	}

	// VULNERABLE: HTTP Header Injection (CWE-113) - raw header write with user input
	headerName, headerValue := h.customizationService.SetCustomHeader(req)

	// Use raw response writer to bypass framework sanitization
	c.Writer.Header().Set(headerName, headerValue)
	c.Writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(c.Writer, `{"success":true,"header":"%s","value":"%s"}`, headerName, headerValue)
}

func (h *HeaderHandler) RedirectWithLang(c *gin.Context) {
	lang := c.Query("lang")
	if lang == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lang parameter required"})
		return
	}

	// Create DTO and call service
	req := &entity.LocaleRequest{
		Lang: lang,
	}

	// VULNERABLE: HTTP Header Injection (CWE-113) - CRLF in redirect Location
	redirectURL := h.localizationService.GetRedirectURL(req)

	c.Writer.Header().Set("X-Language", lang)
	c.Writer.Header().Set("Location", redirectURL)
	c.Writer.WriteHeader(http.StatusFound)
	fmt.Fprintf(c.Writer, `Redirecting to %s`, redirectURL)
}
