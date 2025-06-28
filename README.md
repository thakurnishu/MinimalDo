# 📝 Todo App - Microservices Architecture

A modern, full-stack todo application built with Go backend, vanilla JavaScript frontend, and PostgreSQL database. Designed as microservices that can be deployed independently.

## 🏗️ Architecture

```
Frontend (JavaScript/HTML)  ←→  Backend (Go)  ←→  Database (PostgreSQL)
     Port 80                    Port 8080           Port 5432
```

## 🛠️ Quick Start

### Docker Compose 

1. **Clone and start all services:**
   ```bash
   # Build and start all services
   make build
   make up
   ```

2. **Access the application:**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Database: localhost:5432

3. **View logs:**
   ```bash
   make logs
   ```

## 🔧 Configuration

### Environment Variables

**Backend (`backend/main.go`):**
```bash
DB_HOST=localhost      # Database host
DB_PORT=5432          # Database port
DB_USER=postgres      # Database user
DB_PASSWORD=password  # Database password
DB_NAME=todoapp      # Database name
PORT=8080            # Server port
```

**Frontend (`frontend/index.html`):**
```javascript
// Update the API URL in the TodoApp constructor
this.apiUrl = 'http://localhost:8080/api';
```

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/todos` | Get all todos |
| POST | `/api/todos` | Create new todo |
| PUT | `/api/todos/{id}` | Update todo |
| DELETE | `/api/todos/{id}` | Delete todo |
| GET | `/api/health` | Health check |

### Example API Usage

**Create Todo:**
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Go","description":"Build a todo app","completed":false}'
```

**Get All Todos:**
```bash
curl http://localhost:8080/api/todos
```

## 🎨 UI Features

- **Responsive Design:** Adapts to mobile and desktop
- **Collapsible Sidebar:** Toggle with the floating button (like Claude's interface)
- **Smooth Animations:** Hover effects and transitions
- **Statistics Dashboard:** Real-time task counters
- **Modern Styling:** Glassmorphism effects and gradients

## 🛡️ Production Considerations

### Security
- Add authentication/authorization
- Use HTTPS in production
- Implement rate limiting
- Add input validation and sanitization

### Performance
- Add database connection pooling
- Implement caching (Redis)
- Add database indices
- Use CDN for frontend assets

### Monitoring
- Add logging middleware
- Implement health checks
- Add metrics collection
- Set up alerting
