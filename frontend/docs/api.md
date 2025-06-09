<!-- frontend/docs/api.md -->
# Autosysadmin API Documentation

## Authentication

### Login
`POST /api/v1/auth/login`

Request body:
```json
{
  "email": "user@example.com",
  "password": "password123"
}