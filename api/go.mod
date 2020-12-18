module classroom

go 1.14

require (
	classroom/endpoints v0.0.0
	classroom/functions v0.0.0
	classroom/models v0.0.0
	classroom/utils v0.0.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/rs/cors v1.7.0
	google.golang.org/api v0.36.0
)

replace (
	classroom/endpoints v0.0.0 => ./endpoints
	classroom/functions v0.0.0 => ./functions
	classroom/models v0.0.0 => ./models
	classroom/utils v0.0.0 => ./utils
)
