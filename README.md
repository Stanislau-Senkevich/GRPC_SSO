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

