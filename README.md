# AI of the World - Go Backend API

A high-performance REST API built with Go (Golang) and Gin framework for the AI of the World platform.

## 🚀 Features

- ✅ **Fast & Efficient** - Built with Go for maximum performance
- ✅ **RESTful API** - Clean and intuitive API design
- ✅ **JWT Authentication** - Secure token-based authentication
- ✅ **Role-Based Access** - Admin and user roles
- ✅ **MySQL Database** - Using GORM ORM
- ✅ **File Upload Support** - Handle images, GIFs, and videos
- ✅ **CORS Enabled** - Ready for frontend integration
- ✅ **Tag Management** - Full CRUD operations for tags
- ✅ **Environment Config** - Easy configuration via .env

## 📋 Prerequisites

- Go 1.21 or higher
- MySQL 5.7+ or MariaDB 10.2+
- Git

## 🛠️ Installation

### 1. Clone the repository (if not already done)

```bash
cd "/Users/kshitizmaurya/Documents/Projects/AI OF THE WORLD/Backend"
```

### 2. Install Go dependencies

```bash
go mod download
```

### 3. Set up the database

First, create the database using the SQL files in the `../DataBase` folder:

```bash
mysql -u root -p < ../DataBase/00_complete_setup.sql
```

Or run individual files:

```bash
mysql -u root -p ai_of_the_world < ../DataBase/01_users_table.sql
mysql -u root -p ai_of_the_world < ../DataBase/02_tags_table.sql
# ... continue with remaining files
```

### 4. Configure environment

The `.env` file is already configured with your database credentials. Review and update if needed:

```bash
cat .env
```

### 5. Run the server

```bash
go run main.go
```

Or build and run:

```bash
go build -o ai-backend
./ai-backend
```

The server will start on `http://localhost:8080`

## 📚 API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/register` | Register new user | No |
| POST | `/api/v1/auth/login` | Login user | No |
| GET | `/api/v1/profile` | Get user profile | Yes |

### Tags (Public)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/tags` | Get all tags | No |
| GET | `/api/v1/tags/:id` | Get tag by ID | No |
| GET | `/api/v1/tags/search?q=query` | Search tags | No |
| GET | `/api/v1/tags/stats` | Get tag statistics | No |

### Tags (Admin Only)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/admin/tags` | Create new tag | Yes (Admin) |
| PUT | `/api/v1/admin/tags/:id` | Update tag | Yes (Admin) |
| DELETE | `/api/v1/admin/tags/:id` | Delete tag | Yes (Admin) |

### Health Check

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/health` | Server health check | No |

## 📝 API Usage Examples

### Register a new user

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "SecurePass123!",
    "full_name": "John Doe"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@aioftheworld.com",
    "password": "MyNewPassword123!"
  }'
```

### Get all tags

```bash
curl http://localhost:8080/api/v1/tags
```

### Create a tag (Admin only)

```bash
curl -X POST http://localhost:8080/api/v1/admin/tags \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Cyberpunk",
    "category": "Theme",
    "description": "Futuristic cyberpunk aesthetic"
  }'
```

### Search tags

```bash
curl "http://localhost:8080/api/v1/tags/search?q=cyber&limit=5"
```

## 🔐 Authentication

The API uses JWT (JSON Web Tokens) for authentication. After logging in, you'll receive a token that must be included in the `Authorization` header for protected routes:

```
Authorization: Bearer YOUR_JWT_TOKEN
```

## 📁 Project Structure

```
Backend/
├── config/
│   ├── config.go       # Configuration loader
│   └── database.go     # Database connection
├── controllers/
│   ├── auth.go         # Authentication controllers
│   └── tag.go          # Tag controllers
├── middleware/
│   └── auth.go         # Authentication middleware
├── models/
│   ├── auth.go         # Auth models
│   ├── prompt.go       # Prompt models
│   ├── tag.go          # Tag models
│   └── user.go         # User models
├── routes/
│   └── routes.go       # Route definitions
├── utils/
│   ├── jwt.go          # JWT utilities
│   ├── password.go     # Password hashing
│   └── response.go     # Response utilities
├── uploads/            # File upload directory
├── .env                # Environment variables
├── .env.example        # Environment template
├── .gitignore          # Git ignore rules
├── go.mod              # Go module file
├── main.go             # Application entry point
└── README.md           # This file
```

## 🔧 Configuration

All configuration is done through environment variables in the `.env` file:

- `PORT` - Server port (default: 8080)
- `ENV` - Environment (development/production)
- `DB_*` - Database connection settings
- `JWT_SECRET` - Secret key for JWT tokens
- `UPLOAD_DIR` - Directory for file uploads
- `MAX_UPLOAD_SIZE` - Maximum file upload size in bytes
- `ALLOWED_ORIGINS` - CORS allowed origins
- `FRONTEND_URL` - Frontend application URL

## 🧪 Testing

Test the API using curl, Postman, or any HTTP client.

### Health Check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "message": "AI of the World API is running"
}
```

## 🚀 Deployment

### Build for production

```bash
go build -o ai-backend
```

### Run in production

```bash
ENV=production ./ai-backend
```

## 📊 Database Schema

The backend uses the following main tables:

- `users` - User accounts and profiles
- `tags` - Content categorization tags
- `image_prompts` - AI-generated image submissions
- `gif_prompts` - AI-generated GIF submissions
- `video_prompts` - AI-generated video submissions

See the `../DataBase` folder for complete SQL schema.

## 🛡️ Security

- Passwords are hashed using bcrypt
- JWT tokens expire after 24 hours
- CORS is configured for specific origins
- SQL injection protection via GORM
- Input validation on all endpoints

## 📞 Support

For issues or questions, please check the documentation or contact the development team.

## 📄 License

This project is part of the AI of the World platform.

---

**Version**: 1.0  
**Last Updated**: January 2026  
**Built with**: Go 1.21, Gin, GORM, MySQL
# AI_of_the_world_Backend
