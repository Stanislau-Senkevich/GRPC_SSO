# GRPC SSO Microservice

---------------

## Introduction

This API based on GRPC technology provides interface of SSO service for users
with integration with my another GRPC microservice (https://github.com/Stanislau-Senkevich/GRPC_Family).

To test this API on your own you should download protocols from https://github.com/Stanislau-Senkevich/protocols
and send gRPC requests on <span style="color:blue"> grpc://droplet.senkevichdev.work:33033 </span>

### Models
- Admin
- User


#### User
- Get/Update information about account
- Change account's password

#### Admin
- All user's features
- Delete users with erasing its data in GRPC_Family microservice
- Checking if some another user is admin or not


------------------
## Technologies
- #### Go 1.21
- #### gRPC
- #### MongoDB
- #### Docker
- #### Kubernetes
- #### JWT-tokens
- #### DNS
- #### CI/CD (GitHub Actions)

-----------------
## Realization features
- #### Microservice architecture
- #### Clean architecture
- #### Functional tests for handlers
- #### Linter
- #### Logging with slog package

-----------------
### Tools and libraries

### Protocols and Utilities

- `Stanislau-Senkevich/protocols`: Protocol implementations.
- `badoux/checkmail`: Email address validation.
- `brianvoe/gofakeit`: Random data generation.

### Authentication

- `golang-jwt/jwt`: JWT functionality.

### Middleware

- `grpc-ecosystem/go-grpc-middleware/v2`: provides gRPC middleware.

### Configuration

- `ilyakaznacheev/cleanenv`: Environment variable reading and validation.
- `spf13/viper`: Configuration solution for Go applications.
- `subosito/gotenv`: Loading environment variables from .env files.

### Database

- `go.mongodb.org/mongo-driver`: Go package providing driver and functinality to interact with MongoDB.

### Cryptography

- `golang.org/x/crypto`: Cryptographic algorithms for hashing passwords.

### gRPC

- `google.golang.org/grpc`: gRPC Go implementation.

### Testing

- `stretchr/testify`: Assertion functions.

### Logging

- `log/slog`: standard Go library for logging.

### Protocol Buffers

- `google.golang.org/protobuf`: Protocol Buffers serialization.

