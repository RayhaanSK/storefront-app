export interface UserType {
    id: number
    type: string
}

export interface UserAccount {
    id: number
    username: string
    passwordHash: string
    userTypeId: number
    userType: string
}

export interface Category {
    id: number
    name: string
}

export interface Product {
    id: number
    name: string
    price: number
    description: string
    listedAt: string
    updatedAt: string
    sold: boolean
    sellerId: number
    sellerUsername: string
    categoryId: number
    categoryName: string
}

export interface CartItem {
    customerId: number
    productId: number
    productName: string
    productPrice: number
}

export interface CreateUserAccountRequest {
    username: string
    passwordHash: string
    userTypeId: number
}

export interface LogIntoAccountRequest {
    username: string
    passwordHash: string
}

export interface LogIntoAccountResponse {
    userAccount: UserAccount
    successStatus: boolean
}

export interface CreateProductRequest {
    name: string
    price: number
    description: string
    sellerId: number
    categoryId: number
}

export interface UpdateProductRequest {
    name: string
    price: number
    description: string
    categoryId: number
}

export interface AddProductToCartRequest {
    customerId: number
    productId: number
}