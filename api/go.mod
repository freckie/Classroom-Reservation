module classroom

go 1.14

require (
	github.com/go-sql-driver/mysql v1.5.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/rs/cors v1.7.0

	classroom/endpoints v0.0.0
	classroom/models v0.0.0
	classroom/functions v0.0.0
)

replace (
	classroom/endpoints v0.0.0 => ./endpoints
	classroom/models v0.0.0 => ./models
	classroom/functions v0.0.0 => ./functions
)