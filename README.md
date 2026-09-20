# Bill Splitter

A small web application for splitting a bill between people by percentage. The
application calculates each person's euro amount and assigns any remainder cents
deterministically.

## Requirements

- [Docker](https://docs.docker.com/get-docker/)
- Docker Compose (included with Docker Desktop)

## Run the application

1. Clone the repository and enter the project directory:

   ```sh
   git clone https://github.com/lniklison/bill-splitter-go.git
   cd bill-splitter-go
   ```

2. Build and start the application and PostgreSQL database:

   ```sh
   docker compose up --build
   ```

3. Open [http://localhost:8080](http://localhost:8080) in a browser.

The database migrations run automatically when the application starts. A sample
bill with ID `1` is created so the application is ready to use immediately.

To stop the application, press `Ctrl+C`. To also remove the database volume and
start again with fresh data, run:

```sh
docker compose down --volumes
```

## Verification

Run the backend checks from the repository root:

```sh
go test ./...
```

Run the frontend type check and production build:

```sh
cd frontend
npm ci
npm run build
```
