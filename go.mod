module github.com/formulator2

go 1.19

require (
	github.com/formulator2/combinator/step1/zeroOneTwoTree v0.0.0-00010101000000-000000000000 // indirect
	github.com/formulator2/combinator/combinator v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
)

replace github.com/formulator2/combinator/combinator => ./combinator

replace github.com/formulator2/combinator/step1/zeroOneTwoTree => ./combinator/step1
