module github.com/JinXVIII/BE-Medical-Record

go 1.25.1

require (
	github.com/go-chi/chi/v5 v5.2.3
	github.com/go-sql-driver/mysql v1.9.3
	github.com/golang-jwt/jwt/v5 v5.3.0
	golang.org/x/crypto v0.46.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-chi/cors v1.2.1 // indirect
)

// deps get value .env
require github.com/joho/godotenv v1.5.1

// deps validator
require (
	github.com/gabriel-vasile/mimetype v1.4.11 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.29.0
	github.com/leodido/go-urn v1.4.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
)
