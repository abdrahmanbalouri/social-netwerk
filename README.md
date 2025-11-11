# Social Network Project

This is a full-stack social network application with a Go backend and a Next.js frontend. It's set up to be easily run with Docker.

## Project Structure

*   `backend/`: The Go API.
*   `frontend/`: The Next.js web application.
*   `docker-compose.yml`: Defines the services for running the application.

## Getting Started

To get the application running, you'll need Docker and Docker Compose installed.

1.  **Build and run the services:**

    ```bash
    docker-compose up --build
    ```

    This command will build the Docker images for both the frontend and backend and start the containers.

2.  **Access the application:**
    *   Frontend: [http://localhost:3000](http://localhost:3000)
    *   Backend: [http://localhost:8080](http://localhost:8080)

## Development

You can also run the frontend and backend services locally for development.

### Backend (Go)

1.  Navigate to the backend directory:
    ```bash
    cd backend
    ```
2.  Install dependencies:
    ```bash
    go mod tidy
    ```
3.  Run the server:
    ```bash
    go run ./cmd/
    ```
    The backend will be running on `http://localhost:8080`.

### Frontend (Next.js)

1.  Navigate to the frontend directory:
    ```bash
    cd frontend
    ```
2.  Install dependencies:
    ```bash
    npm install
    ```
3.  Run the development server:
    ```bash
    npm run dev
    ```
    The frontend will be running on `http://localhost:3000`.

    **Note:** When running the frontend locally, you'll need to create a `.env.local` file in the `frontend` directory and set the API URL:
    ```
    NEXT_PUBLIC_API_URL=http://localhost:8080
    ```

## Docker Commands

*   **Run in detached mode:**
    ```bash
    docker-compose up --build -d
    ```
*   **Stop the services:**
    ```bash
    docker-compose down
    ```
*   **View logs:**
    ```bash
    docker-compose logs -f <service_name>  # e.g., backend or frontend
    ```
# Authors :
```bash
@azraji, @abalouri, @abaid, @mennas, @ranniz, @ychatoua 
```