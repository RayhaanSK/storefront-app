package storefront

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, handler *Handler) {
	router.GET("/api/user_type/list", handler.ListAllUserTypes)
	router.POST("/api/user_account/create", handler.SignUpUserAccount)
	router.GET("/api/user_account/list", handler.ListAllUserAccounts)
	router.POST("/api/user_account/login", handler.LogIntoUserAccount)
	router.DELETE("/api/user_account/delete/:id", handler.DeleteUserAccount)
	router.GET("/api/product/list", handler.ListProductCatalogue)
	router.GET("/api/product/:id", handler.GetSingleProductDetails)
	router.GET("/api/product/seller/:id", handler.ListAllSellerProducts)
	router.POST("/api/product/create", handler.CreateProduct)
	router.PUT("/api/product/update/:id", handler.UpdateProduct)
	router.DELETE("/api/product/delete/:id", handler.DeleteProduct)
	router.GET("/api/cart_item/list/:customerId", handler.ListAllCustomerCartItems)
	router.POST("/api/cart_item/create", handler.AddProductToCustomerCart)
	router.DELETE("/api/cart_item/delete/:customerId/:productId", handler.RemoveProductFromCustomerCart)
	router.DELETE("/api/cart_item/delete/:customerId", handler.CheckoutCustomerCart)
}
