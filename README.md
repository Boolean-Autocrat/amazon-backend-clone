# Run Locally

<ul>
    <li> Run <code>npm i</code></li><br>
    <li>Spin up a Postgres database (preferably use Docker)</li><br>
    <li>Create a <code>.env</code> file with the following entries: 
    <code>DB_HOST = ''<br>DB_USER = ''<br>DB_PASSWORD = ''<br>DATABASE_URL = ''
    </code>
    </li><br>
    <li>Run <code>npx dbmate up</code></li><br>
    <li>Next, run <code>go run main.go</code> or use <code>CompileDaemon --command="./amzn"</code></li>
