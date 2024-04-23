# dictionaryAi

Launch Postgres Db in docker

Create tables by running the next command:

`$ psql -h localhost -p 5432 -U YOU_NAME postgres < create_tables.sql`

Build server using this command
`$ go build`
`$ go run dictionaryAi`

Server will be available here http://localhost:1234/