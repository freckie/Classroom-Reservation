module classroom/endpoints

go 1.14

require (
	classroom/functions v0.0.0
	classroom/models v0.0.0
	classroom/utils v0.0.0
	github.com/julienschmidt/httprouter v1.3.0
)

replace (
	classroom/functions v0.0.0 => ../functions
	classroom/models v0.0.0 => ../models
	classroom/utils v0.0.0 => ../utils
)
