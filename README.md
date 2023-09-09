# Features

- [x] User Auth Endpoints
- [x] User Profile Endpoints
- [x] Products Endpoints
- [] Orders Endpoints
- [] Admin Users/Orders Endpoints
- [] Cart Endpoints

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
