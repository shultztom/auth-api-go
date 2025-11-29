# auth-api-go

A RESTful authentication API built with Go, providing user registration, login, JWT-based authentication, session management, and role-based access control.

### Tech Stack

- **Go** - Programming language
- **Gin** - HTTP web framework
- **GORM** - ORM library for database operations
- **PostgreSQL** - Primary database
- **Redis** - Session/token caching
- **JWT** - JSON Web Tokens for authentication
- **bcrypt** - Password hashing

### Setup

#### Prerequisites

- Go 1.18+
- PostgreSQL
- Redis

#### Environment Variables

Create a `.env` file in the project root with the following variables:

```env
# PostgreSQL Configuration
PG_USER=your_postgres_user
PG_PASS=your_postgres_password
PG_DB=your_database_name
PG_HOST=localhost

# Redis Configuration
REDIS_URL=localhost:6379

# JWT Secrets
JWT_SECRET=your_jwt_secret_key
JWT_APP_SECRET=your_app_jwt_secret_key

# Cloud Deployment (optional)
IS_CLOUD=false
```

#### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd auth-api-go
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables (see above)

4. Run the application:
   ```bash
   go run main.go
   ```

#### Running Tests

```bash
go test ./... -v
```

### API Routes

#### User Authentication

##### POST - /register

Register a new user.

Body:
```json
{
    "username": "test",
    "password": "123"
}
```

Response: `201 Created`
```json
{
    "token": "<jwt_token>"
}
```

##### POST - /login

Authenticate an existing user.

Body:
```json
{
    "username": "test",
    "password": "123"
}
```

Response: `200 OK`
```json
{
    "token": "<jwt_token>"
}
```

##### GET - /verify

Verify a user's JWT token.

Headers:
```
x-auth-token: <jwt_token>
```

Response: `200 OK`
```json
{
    "message": "success"
}
```

##### DELETE - /

Delete the authenticated user's account.

Headers:
```
x-auth-token: <jwt_token>
```

Response: `200 OK`
```json
{
    "Deleted user": "<username>"
}
```

##### DELETE - /session

Delete the authenticated user's session (logout).

Headers:
```
x-auth-token: <jwt_token>
```

Response: `200 OK`
```json
{
    "Deleted session for user": "<username>"
}
```

#### Role Management

##### GET - /roles

Get all roles for the authenticated user.

Headers:
```
x-auth-token: <jwt_token>
```

Response: `200 OK`
```json
{
    "Roles": [
        {
            "username": "test",
            "role": "admin"
        }
    ]
}
```

##### GET - /roles/:role

Check if the authenticated user has a specific role.

Headers:
```
x-auth-token: <jwt_token>
```

Response: `200 OK`
```json
{
    "hasRoleAlready": true
}
```

##### POST - /roles

Add a role to the authenticated user.

Headers:
```
x-auth-token: <jwt_token>
```

Body:
```json
{
    "role": "role-name"
}
```

Response: `201 Created`
```json
{
    "Added Role": "role-name"
}
```

#### App Authentication (Service-to-Service)

##### GET - /app/verify

Verify an application JWT token.

Headers:
```
X-API-Token: <app_jwt_token>
```

Response: `200 OK`
```json
{
    "message": "success"
}
```

##### DELETE - /app/user/:username

Delete a user by username (app-level access).

Headers:
```
X-API-Token: <app_jwt_token>
```

Response: `200 OK`
```json
{
    "Deleted user": "<username>"
}
```