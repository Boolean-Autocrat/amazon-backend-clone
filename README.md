# Features

- [x] User Auth Endpoints
- [x] User Profile Endpoints
- [x] Products Endpoints
- [x] Orders Endpoints
- [x] Admin Users/Orders Endpoints
- [x] Cart Endpoints
- [x] JWT-Based Authentication w/ server-side cookie-setting
- [x] Multiple checkpoints and validation

# Run Locally

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
