# 📝 Todo App - Microservices Architecture

A modern, full-stack todo application built with Go backend, vanilla JavaScript frontend, and PostgreSQL database. Designed as microservices that can be deployed independently.

## 🏗️ Architecture

```
Frontend (JavaScript/HTML)  ←→  Backend (Go)  ←→  Database (PostgreSQL)
     Port 80                    Port 8080           Port 5432
```

## 🚀 Features

- ✅ Add, edit, delete todos
- ✅ Mark todos as complete/incomplete  
- ✅ Beautiful, responsive UI with animations
- ✅ Collapsible sidebar (like Claude's chat history)
- ✅ Real-time statistics
- ✅ Microservices architecture
- ✅ Docker containerization
- ✅ Independent deployment capability

## 📁 Project Structure

```
todoapp/
├── backend/
│   ├── main.go              # Go server with REST API
│   ├── go.mod               # Go dependencies
│   └── Dockerfile           # Backend containerization
├── frontend/
│   ├── index.html           # Single-page application
│   ├── nginx.conf           # Nginx configuration
│   └── Dockerfile           # Frontend containerization
├── docker-compose.yml       # Multi-service orchestration
├── Makefile                 # Development commands
└── README.md               # This file
```

## 🛠️ Quick Start

### Option 1: Docker Compose (Recommended)

1. **Clone and start all services:**
   ```bash
   # Build and start all services
   make build
   make up
   
   # Or manually:
   docker-compose up --build -d
   ```

2. **Access the application:**
   - Frontend: http://localhost
   - Backend API: http://localhost:8080
   - Database: localhost:5432

3. **View logs:**
   ```bash
   make logs
   # Or: docker-compose logs -f
   ```

### Option 2: Development Mode

1. **Start PostgreSQL:**
   ```bash
   # Using Docker
   docker run --name postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=todoapp -p 5432:5432 -d postgres:15-alpine
   
   # Or install locally and create database
   createdb todoapp
   ```

2. **Run Backend:**
   ```bash
   cd backend
   go mod tidy
   go run main.go
   
   # Or use make command
   make dev-backend
   ```

3. **Serve Frontend:**
   ```bash
   cd frontend
   python3 -m http.server 8000
   
   # Or use make command
   make dev-frontend
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

## 🐳 Docker Deployment

### Individual Service Deployment

**Backend only:**
```bash
cd backend
docker build -t todo-backend .
docker run -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-password \
  todo-backend
```

**Frontend only:**
```bash
cd frontend
docker build -t todo-frontend .
docker run -p 80:80 todo-frontend
```

**Database only:**
```bash
docker run --name postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=todoapp \
  -p 5432:5432 \
  -d postgres:15-alpine
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

## 🔄 Development Workflow

```bash
# Start development environment
make up

# View logs
make logs

# Stop services
make down

# Clean up everything
make clean

# Run backend tests
make test
```
