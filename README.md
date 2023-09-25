# Features

- [x] User Auth Endpoints
- [x] User Profile Endpoints
- [x] Products Endpoints
- [x] Orders Endpoints
- [x] Admin Users/Orders Endpoints
- [x] Cart Endpoints
- [x] JWT-Based Authentication w/ server-side cookie-setting
- [x] Multiple checkpoints and validation

# API Documentation

The API documentation can be found [here](https://documenter.getpostman.com/view/28952349/2s9YJW4Qna)

# Run Locally (with Docker)

- Create a `.env` file with the following entries:

  ```
  DB_HOST = 'postgres'
  DB_USER = 'postgres'
  DB_PASSWORD = 'postgres'
  DB_NAME = 'postman_amzn'
  DB_SOURCE='postgresql://postgres:postgres@postgres:5432/postman_amzn?sslmode=disable'
  PORT='8080'
  ```

- Run `docker compose up`

- The API will be available at `localhost:8080`

# Run Locally (without Docker)

- Run `npm i`

- Spin up a Postgres database (preferably use Docker)

- Create a `.env` file with the following entries:

  ```
  DB_HOST = ''
  DB_USER = ''
  DB_PASSWORD = ''
  DATABASE_URL = ''
  ```

- Run `npx dbmate up`

- Next, run `go run main.go` (or use `CompileDaemon --command="./amzn"`)
