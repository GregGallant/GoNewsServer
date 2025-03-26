module gallantone.com/main

go 1.24.0

require (
	// OAuth
	github.com/coreos/go-oidc/v3 v3.13.0
	// CHI
	github.com/go-chi/chi/v5 v5.2.1
	// JWT (outdated?)
	github.com/golang-jwt/jwt v3.2.2+incompatible
	// CORS
	github.com/rs/cors v1.11.1
	// Cryptography
	golang.org/x/crypto v0.36.0
	golang.org/x/oauth2 v0.28.0
	// ORM
	gorm.io/driver/postgres v1.5.11
	gorm.io/gorm v1.25.12
)

require (
	github.com/go-jose/go-jose/v4 v4.0.5 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.4 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)
