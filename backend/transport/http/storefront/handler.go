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

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	SellerId    int64   `json:"sellerId"`
	CategoryId  int64   `json:"categoryId"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	CategoryId  int64   `json:"categoryId"`
}

type AddProductToCartRequest struct {
	CustomerId int64 `json:"customerId"`
	ProductId  int64 `json:"productId"`
}

type Handler struct {
	service *storefrontservice.Service
}

func NewHandler(service *storefrontservice.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// ListAllUserTypes handles GET /api/user_type/list
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

	c.JSON(http.StatusOK, response)
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

// ListProductCatalogue handles GET /api/product/list
func (h *Handler) ListProductCatalogue(c *gin.Context) {
	products, err := h.service.ListProductCatalogue(c.Request.Context())
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetSingleProductDetails handles GET /api/product/:id
func (h *Handler) GetSingleProductDetails(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product id",
		})
		return
	}

	product, err := h.service.GetSingleProductDetails(c.Request.Context(), id)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

// ListAllSellerProducts handles GET /api/product/seller/:id
func (h *Handler) ListAllSellerProducts(c *gin.Context) {
	sellerId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid seller id",
		})
		return
	}

	products, err := h.service.ListAllSellerProducts(c.Request.Context(), sellerId)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, products)
}

// CreateProduct handles POST /api/product/create
func (h *Handler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	product, err := h.service.CreateProduct(
		c.Request.Context(),
		req.Name,
		req.Price,
		req.Description,
		req.SellerId,
		req.CategoryId,
	)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct handles PUT /api/product/update/:id
func (h *Handler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product id",
		})
		return
	}

	var req UpdateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	product, err := h.service.UpdateProduct(
		c.Request.Context(),
		id,
		req.Name,
		req.Price,
		req.Description,
		req.CategoryId,
	)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct handles DELETE /api/product/delete/:id
func (h *Handler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product id",
		})
		return
	}

	if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		h.writeError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListAllCustomerCartItems handles GET /api/cart_item/list/:customerId
func (h *Handler) ListAllCustomerCartItems(c *gin.Context) {
	customerId, err := strconv.ParseInt(c.Param("customerId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid customer id",
		})
		return
	}

	products, err := h.service.ListAllCustomerCartItems(c.Request.Context(), customerId)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, products)
}

// AddProductToCustomerCart handles POST /api/cart_item/create
func (h *Handler) AddProductToCustomerCart(c *gin.Context) {
	var req AddProductToCartRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	cartItem, err := h.service.AddProductToCustomerCart(
		c.Request.Context(),
		req.CustomerId,
		req.ProductId,
	)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, cartItem)
}

// RemoveProductFromCustomerCart handles DELETE /api/cart_item/delete/:customerId/:productId
func (h *Handler) RemoveProductFromCustomerCart(c *gin.Context) {
	customerId, err := strconv.ParseInt(c.Param("customerId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid customer id",
		})
		return
	}

	productId, err := strconv.ParseInt(c.Param("productId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product id",
		})
		return
	}

	if err := h.service.RemoveProductFromCustomerCart(c.Request.Context(), customerId, productId); err != nil {
		h.writeError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// CheckoutCustomerCart handles DELETE /api/cart_item/delete/:customerId
func (h *Handler) CheckoutCustomerCart(c *gin.Context) {
	customerId, err := strconv.ParseInt(c.Param("customerId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid customer id",
		})
		return
	}

	if err := h.service.CheckoutCustomerCart(c.Request.Context(), customerId); err != nil {
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

	case errors.Is(err, storefrontservice.ErrIncorrectUsernameOrPassword):
		c.JSON(http.StatusNotFound, gin.H{
			"error": "incorrect username or password",
		})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}
}
