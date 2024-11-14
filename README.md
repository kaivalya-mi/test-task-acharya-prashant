Here's the updated README.md with the images and their text aligned:

# 🚀 Aacharya Prashant Test Task

[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)](https://golang.org/doc/go1.20)
[![Gin Framework](https://img.shields.io/badge/Gin-Framework-00ADD8?style=flat&logo=go)](https://github.com/gin-gonic/gin)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-336791?style=flat&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7+-DC382D?style=flat&logo=redis&logoColor=white)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=flat&logo=docker&logoColor=white)](https://www.docker.com/)

A robust REST API service built with modern technologies for efficient user management and authentication.

## ⚡️ Tech Stack

<table>
  <tr>
    <td align="center" width="33%">
      <img src="https://raw.githubusercontent.com/gin-gonic/logo/master/color.png" width="80" alt="Gin Framework">
      <h3>Gin Framework</h3>
      <p>• High-performance HTTP web framework<br>
      • Built-in middleware support<br>
      • Custom validation<br>
      • Efficient routing</p>
    </td>
    <td align="center" width="33%">
      <img src="https://wiki.postgresql.org/images/a/a4/PostgreSQL_logo.3colors.svg" width="65" alt="PostgreSQL">
      <h3>PostgreSQL</h3>
      <p>• Robust relational database<br>
      • ACID compliance<br>
      • Strong data integrity</p>
    </td>
    <td align="center" width="33%">
      <img src="https://redis.io/wp-content/uploads/2024/04/Logotype.svg?auto=webp&quality=85,75&width=120" width="200" alt="Redis">
      <h3>Redis</h3>
      <p>• In-memory caching<br>
      • Session management<br>
      • High-performance storage</p>
    </td>
  </tr>
</table>

### 🐳 Docker Setup

Our application uses Docker for containerization, ensuring consistent development and deployment environments.

#### Prerequisites
- Docker
- Docker Compose

#### Components

<table>
  <tr>
    <td align="center" width="33%">
      <img src="https://img.icons8.com/color/48/000000/api.png" width="70" alt="API Service">
      <h3>API Service</h3>
      <p>Golang application</p>
    </td>
    <td align="center" width="33%">
      <img src="https://img.icons8.com/color/48/000000/database.png" width="70" alt="PostgreSQL">
      <h3>PostgreSQL</h3>
      <p>Database service</p>
    </td>
    <td align="center" width="33%">
      <img src="https://img.icons8.com/color/48/000000/redis.png" width="70" alt="Redis">
      <h3>Redis</h3>
      <p>Cache service</p>
    </td>
  </tr>
</table>

#### Commands
```bash
# 🚀 Start all services
docker compose up -d --build

# 📊 View running containers
docker ps

# 📝 View container logs
docker logs <container_name>

# 🛑 Stop all containers
docker compose down

# 🗑️ Remove containers and volumes
docker compose down -v
```

## 📌 API Documentation

### Authentication Endpoints

#### 1. Sign-Up
```bash
curl -X POST http://localhost:8080/api/v1/sign-up \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@mailinator.com",
    "password": "secretpassword"
  }'
```

#### 2. Sign-In
```bash
curl -X POST http://localhost:8080/api/v1/sign-in \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@mailinator.com",
    "password": "secretpassword"
  }'
```

#### 3. Get User Profile
```bash
curl -X GET http://localhost:8080/api/v1/user-profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

#### 4. Sign-Out
```bash
curl -X POST http://localhost:8080/api/v1/sign-out \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

#### 5. Refresh-Token
```bash
curl -X POST http://localhost:8080/api/v1/refresh-token \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "REFRESH_TOKEN"}'
```

### Response Examples

#### ✅ Success Response
```json
{
  "data": {
    "id": "5576b88b-a255-11ef-a40e-408d5cd209e3",
    "first_name": "John",
    "last_name": "Doe",
    "email": "kaiser@mailinator.com",
    "created_at": "2024-11-14T12:24:46.040005+05:30"
  },
  "meta": {
    "message": "user profile fetched successfully",
    "res_code": 200
  }
}
```

#### ❌ Error Response
```json
{
  "message": "email is not registered with us",
  "res_code": 400
}
```

## 🔒 Security Notes

- 🔑 Authentication required for protected endpoints via `Authorization: Bearer TOKEN`
- ⏰ Access tokens valid for 30 minutes
- 🔄 Refresh tokens valid for 24 hours
- 🛡️ Rate limiting and security middleware enabled

## 📈 Performance

- High-performance routing with Gin
- Redis caching for optimal response times
- Connection pooling for database efficiency

To run the application locally, you will need to have PostgreSQL, Golang, and Redis installed and configured with your `config.toml` file.