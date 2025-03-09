# To Run
- Build the application: `go buiid`
- Run the executable: `./service-test-runner.exe`

# Migrate database
- Delete table of migration: `migrate -path ./migrations -database "$env:MYSQL_DSN" down`(windows) or `migrate -path ./migrations -database "$MYSQL_DSN" down` (linux)
- Create table of migration: `migrate -path ./migrations -database "$env:MYSQL_DSN" up` (windows) or `migrate -path ./migrations -database "$MYSQL_DSN" up` (linux)

# Export env db config
`$env:MYSQL_DSN = "mysql://root:@tcp(127.0.0.1:3306)/db_omnirunner"` (windows)
or
`export MYSQL_DSN="mysql://root:@tcp(127.0.0.1:3306)/db_omnirunner"` (linux)