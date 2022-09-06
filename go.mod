module aces/plankton

go 1.19

require (
	github.com/go-sql-driver/mysql v1.6.0
	uta.edu/aces/jadesdk v0.0.0
)

require github.com/ant0ine/go-json-rest v3.3.2+incompatible // indirect

replace uta.edu/aces/jadesdk => ../jadesdk
