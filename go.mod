module github.com/pasataleo/go-logging

go 1.23

replace github.com/pasataleo/go-testing => ../go-testing
replace github.com/pasataleo/go-errors => ../go-errors

require (
	github.com/pasataleo/go-errors v0.1.2
	github.com/pasataleo/go-testing v0.1.1
)

require github.com/google/go-cmp v0.5.9 // indirect
