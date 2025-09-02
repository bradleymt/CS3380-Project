# CS3380-Docker

A Dockerized environment for Mizzou's CS3380 Databases class. This repository provides a quick way to run a PostgreSQL database, a Go web application, and a pgAdmin GUI for database management.

## Services

- **PostgreSQL**: Persistent database storage using Docker volumes.
- **Go Webapp**: Connects to PostgreSQL and demonstrates database-backed web development. Go module cache is persisted for faster builds. Air is included for hot-reloading.
- **pgAdmin**: Web-based GUI for managing and viewing PostgreSQL data.

## Getting Started

1. Clone this repository:
   ```sh
   git clone https://github.com/notSam25/CS3380-Docker.git
   cd CS3380-Docker
   ```
2. Modify `.env` for required environment variables.
3. Build and start all services:
   ```sh
   docker compose up --build
   ```

## Usage
- Access the webapp at `http://localhost:<webapp-port>` (set in your `.env` file).
- Access pgAdmin at `http://localhost:8080` (default).
- Database data and Go module cache are persisted between runs for faster startup and reliability.

## Troubleshooting
- If you encounter issues, try rebuilding containers:
  ```sh
  docker compose down -v
  docker compose up --build
  ```
- Ensure your `.env` file is correctly configured.

## License
MIT