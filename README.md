## Smart Serve Backend

Smart Serve Backend is built with **Gin**, **Gorm**, and **MySQL**.

## How to Run

- `docker compose up -d`: Start MySQL service and phpMyAdmin.
- Create a database with the name `smart_serve`.
- `go mod tidy`: Download all required packages.
- `air`: Start the server with hot reload (use `go run .` to start the server without hot reload).
- `swag init`: Generate Swagger documentation.
