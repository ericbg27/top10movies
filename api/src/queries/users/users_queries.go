package users

const (
	QueryInsertUser     = "INSERT INTO users (first_name,last_name,email,date_created,status,password) VALUES ($1,$2,$3,$4,$5,$6);"
	QueryInsertUserName = "insert-user-query"

	QueryGetUser     = "SELECT id, first_name, status, password FROM users WHERE email=$1;"
	QueryGetUserName = "get-user-query"

	QueryGetUserById     = "SELECT first_name, last_name, email, status, password FROM users WHERE id=$1;"
	QueryGetUserByIdName = "get-user-by-id-query"

	QueryUpdateUser     = "UPDATE users SET first_name=$1, last_name=$2, email=$3 WHERE id=$4;"
	QueryUpdateUserName = "update-user-query"

	QueryDeleteUser     = "DELETE FROM users WHERE id=$1;"
	QueryDeleteUserName = "delete-user-query"

	QuerySearchUser     = "SELECT id, first_name, last_name, email FROM users WHERE first_name ILIKE '' || $1 || '%' AND last_name ILIKE '%' || $2 || '%';"
	QuerySearchUserName = "search-user-query"
)
