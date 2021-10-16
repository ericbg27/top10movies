package movies

const (
	QueryAddMovie = "INSERT INTO movies VALUES ($1,$2,$3,$4,$5,$6,$7);"
	QueryGetMovie = "SELECT * FROM movies WHERE id=$1;"

	ReleaseDateLayout = "2006-02-01"
	CreatedAtLayout   = "2006-02-01T15:04:05Z"
)
