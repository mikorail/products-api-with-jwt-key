# Products API with JWT

## Overview

The **Products API with JWT** is a RESTful API built using Go (Golang) and the Gin framework. This API allows users to manage products while ensuring secure access through JSON Web Tokens (JWT). It utilizes SQLite as the primary database and includes functionality for rate limiting and token invalidation upon logout.

## Features

- User authentication using JWT
- CRUD operations for managing products
- Secure endpoints requiring authentication
- Rate limiting to protect endpoints from excessive requests
- Token invalidation upon logout, ensuring logged-out JWTs cannot be reused
- Swagger documentation for easy API exploration

## Technologies Used

- **Go**: The programming language used for building the API.
- **Gin**: A web framework for Go, used for handling HTTP requests.
- **GORM**: An ORM for Go, used for interacting with the SQLite database.
- **Swagger**: A tool for documenting and testing APIs.
- **Rate Limiter**: Built-in rate limiting to prevent abuse of endpoints.

## Getting Started

### Prerequisites

Make sure you have the following installed:

- [Go](https://golang.org/dl/) (version 1.16 or later)
- SQLite database (included by default)

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/products-api-with-jwt.git
   cd products-api-with-jwt
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the application**:
   ```bash
   go run main.go
   ```

   The server will start on `http://localhost:8080`.

### API Endpoints

#### Authentication

- **Login**
  - **Endpoint**: `/auth/login`
  - **Method**: `POST`
  - **Request Body**:
    ```json
    {
      "username": "admin",
      "password": "password123",
      "rememberMe": true
    }
    ```
  - **Response**:
    ```json
    {
      "status": "success",
      "code": 200,
      "message": "Login successful",
      "data": {
        "token": "your_jwt_token_here"
      }
    }
    ```

- **Logout**
  - **Endpoint**: `/auth/logout`
  - **Method**: `POST`
  - **Headers**: `Authorization: Bearer <jwt_token>`
  - **Response**:
    ```json
    {
      "status": "success",
      "code": 200,
      "message": "Logout successful",
      "data": null
    }
    ```

  Upon logout, the JWT is invalidated, ensuring it cannot be reused to access secure endpoints.

#### Products

All product-related endpoints require a valid JWT token in the `Authorization` header.

- **Get All Products**
  - **Endpoint**: `/products`
  - **Method**: `GET`
  - **Response**:
    ```json
    {
      "status": "success",
      "code": 200,
      "message": "Products retrieved successfully",
      "data": [ ... ]
    }
    ```

- **Get Product by ID**
  - **Endpoint**: `/products/:id`
  - **Method**: `GET`
  - **Response**:
    ```json
    {
      "status": "success",
      "code": 200,
      "message": "Product retrieved successfully",
      "data": { ... }
    }
    ```

- **Create Product**
  - **Endpoint**: `/products`
  - **Method**: `POST`
  - **Request Body**:
    ```json
    {
      "nama_produk": "Produk A",
      "deskripsi": "Deskripsi Produk A",
      "harga": 1000,
      "stok": 10
    }
    ```
  - **Response**:
    ```json
    {
      "status": "success",
      "code": 201,
      "message": "Product created successfully",
      "data": { ... }
    }
    ```

- **Update Product**
  - **Endpoint**: `/products/:id`
  - **Method**: `PUT`
  - **Request Body**: Same as Create Product
  - **Response**: Similar to Create Product response

- **Delete Product**
  - **Endpoint**: `/products/:id`
  - **Method**: `DELETE`
  - **Response**:
    ```json
    {
      "status": "success",
      "code": 204,
      "message": "Product deleted successfully",
      "data": null
    }
    ```

### Rate Limiting

Rate limiting is enabled on certain endpoints to prevent abuse by limiting the number of requests allowed within a specified timeframe. If the rate limit is exceeded, the following response is returned:

```json
{
  "status": "error",
  "code": 429,
  "message": "Too many requests, please try again later.",
  "data": null
}
```

### API Documentation

You can access the Swagger documentation by navigating to `http://localhost:8080/swagger/index.html` in your web browser.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
