package storefront

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	domain "storefrontapp/domain/storefront"
	storefrontservice "storefrontapp/services/storefront"
)

type CreateUserAccountRequest struct {
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
	UserTypeId   int64  `json:"userTypeId"`
}

type LogIntoAccountRequest struct {
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
}

type LogIntoAccountResponse struct {
	UserAccount   domain.UserAccount `json:"userAccount"`
	SuccessStatus bool               `json:"successStatus"`
}

type Handler struct {
	service *storefrontservice.Service
}

// ListAllUserTypes handles GET /api/user_types
func (h *Handler) ListAllUserTypes(c *gin.Context) {
	userTypes, err := h.service.ListAllUserTypes(c.Request.Context())
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, userTypes)
}

// SignUpUserAccount handles POST /api/user_account/create
func (h *Handler) SignUpUserAccount(c *gin.Context) {
	var req CreateUserAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	userAccount, err := h.service.SignUpUserAccount(
		c.Request.Context(),
		req.Username,
		req.PasswordHash,
		req.UserTypeId,
	)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, userAccount)
}

// ListAllUserAccounts handles GET /api/user_account/list
func (h *Handler) ListAllUserAccounts(c *gin.Context) {
	userAccounts, err := h.service.ListAllUserAccounts(c.Request.Context())
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, userAccounts)
}

// LogIntoUserAccount handles POST /api/user_account/login
func (h *Handler) LogIntoUserAccount(c *gin.Context) {
	var req LogIntoAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	userAccount, successStatus, err := h.service.LogIntoUserAccount(
		c.Request.Context(),
		req.Username,
		req.PasswordHash,
	)
	if err != nil {
		h.writeError(c, err)
		return
	}

	response := LogIntoAccountResponse{
		UserAccount:   userAccount,
		SuccessStatus: successStatus,
	}

	c.JSON(http.StatusCreated, response)
}

// DeleteUserAccount handles DELETE /api/user_account/delete/:id
func (h *Handler) DeleteUserAccount(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user account id",
		})
		return
	}

	if err := h.service.DeleteUserAccount(c.Request.Context(), id); err != nil {
		h.writeError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrTitleRequired):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "title is required",
		})

	case errors.Is(err, storefrontservice.ErrUserAccountNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user account not found",
		})

	case errors.Is(err, storefrontservice.ErrProductNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": "product not found",
		})

	case errors.Is(err, storefrontservice.ErrCartItemNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": "cart item not found",
		})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}
}
