import {
    apiRequest,
} from "./apiClient"

import type {
    UserType,
    UserAccount,
    Category,
    Product,
    CartItem,
    CreateUserAccountRequest,
    LogIntoAccountRequest,
    LogIntoAccountResponse,
    CreateProductRequest,
    UpdateProductRequest,
    AddProductToCartRequest,
} from "../models/storefront"

// region UserType

export function listAllUserTypes(): Promise<UserType[]> {
    return apiRequest<UserType[]>("/api/user_type/list")
}

// endregion

// region UserAccount

export function signUpUserAccount(request: CreateUserAccountRequest): Promise<UserAccount> {
    return apiRequest<UserAccount>("/api/user_account/create", {
        method: "POST",
        body: JSON.stringify(request),
    })
}

export function listAllUserAccounts(): Promise<UserAccount[]> {
    return apiRequest<UserAccount[]>("/api/user_account/list")
}

export function logIntoUserAccount(request: LogIntoAccountRequest): Promise<LogIntoAccountResponse> {
    return apiRequest<LogIntoAccountResponse>("/api/user_account/login", {
        method: "POST",
        body: JSON.stringify(request),
    })
}

export function deleteUserAccount(id: number): Promise<void> {
    return apiRequest<void>(`/api/user_account/delete/${id}`, {
        method: "DELETE",
    })
}

// endregion


// region Product

export function listProductCatalogue(): Promise<Product[]> {
    return apiRequest<Product[]>("/api/product/list")
}

export function getSingleProductDetails(id: number): Promise<Product> {
    return apiRequest<Product>(`/api/product/${id}`)
}

export function listAllSellerProducts(sellerId: number): Promise<Product[]> {
    return apiRequest<Product[]>(`/api/product/seller/${sellerId}`)
}

export function createProduct(request: CreateProductRequest): Promise<Product> {
    return apiRequest<Product>("/api/product/create", {
        method: "POST",
        body: JSON.stringify(request)
    })
}

export function updateProduct(id: number, request: UpdateProductRequest): Promise<Product> {
    return apiRequest<Product>(`/api/product/update/${id}`, {
        method: "PUT",
        body: JSON.stringify(request)
    })
}

export function deleteProduct(id: number): Promise<void> {
    return apiRequest<void>(`/api/product/delete/${id}`, {
        method: "DELETE",
    })
}

// endregion

// region Cart Item

export function listAllCustomerCartItems(customerId: number): Promise<Product[]> {
    return apiRequest<Product[]>(`/api/cart_item/list/${customerId}`)
}

export function addProductToCustomerCart(request: AddProductToCartRequest): Promise<CartItem> {
    return apiRequest<CartItem>("/api/cart_item/create", {
        method: "POST",
        body: JSON.stringify(request)
    })
}

export function removeProductFromCustomerCart(customerId: number, productId: number): Promise<void> {
    return apiRequest<void>(`/api/cart_item/delete/${customerId}/${productId}`, {
        method: "DELETE",
    })
}

export function checkoutCustomerCart(customerId: number): Promise<void> {
    return apiRequest<void>(`/api/cart_item/delete/${customerId}`, {
        method: "DELETE",
    })
}

// endregion