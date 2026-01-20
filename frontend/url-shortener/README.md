# URL Shortener

A simple URL shortener application with Go backend and React frontend.

## Prerequisites

- Go 1.19 or higher
- Node.js 18 or higher
- npm or yarn

## Running the Application

### Backend (Go)

1. Navigate to the backend directory:

    ```
    cd backend
    ```

2. Install dependencies:

    ```
    go mod tidy
    ```

3. Run the server:
    ```
    go run cmd/server/main.go
    ```

The backend will start on `http://localhost:8080`

### Frontend (React)

1. Navigate to the frontend directory:

    ```
    cd frontend/url-shortener
    ```

2. Install dependencies:

    ```
    npm install
    ```

3. Start the development server:
    ```
    npm run dev
    ```

The frontend will start on `http://localhost:3000`

## Usage

1. Open your browser and go to `http://localhost:3000`
2. Enter a long URL in the input field
3. Click "Shorten URL" to generate a short URL
4. Use the generated short URL to redirect to the original URL
