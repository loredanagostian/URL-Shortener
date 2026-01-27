# URL Shortener

A full-stack URL shortening application with a Go backend and React frontend.

## Project Structure

```
url-shortener/
├── backend/                 # Go backend API
│   ├── cmd/server/          # Application entry point
│   ├── internal/            # Internal packages
│   ├── tests/               # Test files
│   ├── docs/                # Documentation
│   ├── docker-compose.yml   # Database containers
│   └── ...
└── frontend/
    └── url-shortener/       # React + Vite frontend
        ├── src/
        └── ...
```

## Prerequisites

- **Go** 1.25.3 or higher
- **Node.js** 18+ and npm
- **Docker Desktop** (for PostgreSQL database)
- **Git**

## Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd url-shortener
```

### 2. Start the Database

Make sure Docker Desktop is running, then:

```bash
cd backend
docker-compose up -d
```

This starts:

- **PostgreSQL** on port `5432`
- **Adminer** (DB admin UI) on port `8081`

### 3. Run the Backend

```bash
cd backend
cp .env.example .env    # First time only
go run cmd/server/main.go
```

The API server starts on `http://localhost:8080`

### 4. Run the Frontend

```bash
cd frontend/url-shortener
npm install             # First time only
npm run dev
```

The frontend starts on `http://localhost:5173`

---

## Detailed Setup

### Database

#### Using Docker (Recommended)

```bash
cd backend
docker-compose up -d
```

| Service    | URL                     | Credentials         |
| ---------- | ----------------------- | ------------------- |
| PostgreSQL | `localhost:5432`        | postgres / postgres |
| Adminer    | `http://localhost:8081` | postgres / postgres |

#### Stop Database

```bash
docker-compose down
```

#### Using Local PostgreSQL

If you have PostgreSQL installed locally:

1. Create a database named `urlshortener`
2. Update `backend/.env` with your connection string:
    ```
    DATABASE_URL=postgres://user:password@localhost:5432/urlshortener?sslmode=disable
    ```

### Backend

#### Environment Variables

Copy the example file and customize if needed:

```bash
cd backend
cp .env.example .env
```

| Variable     | Default                                                                    | Description           |
| ------------ | -------------------------------------------------------------------------- | --------------------- |
| DATABASE_URL | `postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable` | PostgreSQL connection |
| PORT         | `8080`                                                                     | API server port       |

#### Run in Development

```bash
cd backend
go run cmd/server/main.go
```

#### Build for Production

```bash
cd backend
go build -o bin/server cmd/server/main.go
./bin/server
```

#### Run Tests

```bash
cd backend
go test ./tests/...
```

### Frontend

#### Install Dependencies

```bash
cd frontend/url-shortener
npm install
```

#### Run in Development

```bash
npm run dev
```

Opens at `http://localhost:5173`

#### Build for Production

```bash
npm run build
```

Output is in `dist/` folder.

#### Preview Production Build

```bash
npm run preview
```

---

## API Endpoints

| Method | Endpoint                     | Description              |
| ------ | ---------------------------- | ------------------------ |
| POST   | `/api/shorten`               | Create a shortened URL   |
| GET    | `/{shortCode}`               | Redirect to original URL |
| GET    | `/api/urls`                  | List all URLs            |
| GET    | `/api/urls/{shortCode}`      | Get URL details          |
| DELETE | `/api/urls/{shortCode}`      | Delete a URL             |
| GET    | `/api/analytics/{shortCode}` | Get URL analytics        |

---

## Development Workflow

### Running Everything Together

Open three terminal windows:

**Terminal 1 - Database:**

```bash
cd backend
docker-compose up
```

**Terminal 2 - Backend:**

```bash
cd backend
go run cmd/server/main.go
```

**Terminal 3 - Frontend:**

```bash
cd frontend/url-shortener
npm run dev
```

### Stopping Everything

1. Press `Ctrl+C` in the frontend terminal
2. Press `Ctrl+C` in the backend terminal
3. Stop database: `docker-compose down` (in backend folder)

---

## Troubleshooting

### Database Connection Refused

```
dial tcp [::1]:5432: connectex: No connection could be made because the target machine actively refused it
```

**Solution:** Start Docker Desktop and run `docker-compose up -d`

### Docker Not Found

```
unable to get image 'postgres:15-alpine': error during connect
```

**Solution:** Install and start Docker Desktop

### Port Already in Use

**Backend (8080):**

```bash
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

**Frontend (5173):**

```bash
# Change port in vite.config.ts or use:
npm run dev -- --port 3000
```

---

## Tech Stack

### Backend

- **Go** - Programming language
- **gorilla/mux** - HTTP router
- **lib/pq** - PostgreSQL driver
- **PostgreSQL** - Database

### Frontend

- **React 19** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool
- **QRCode** - QR code generation

---

## Screenshots

Main URL shortener form

![alt text](image.png)

![alt text](image-2.png)

URL shortener view

![alt text](image-3.png)

![alt text](image-4.png)

URL history list

![alt text](image-5.png)

## License

MIT License
