module github.com/shengyanli1982/events/test

go 1.23

replace github.com/shengyanli1982/events => ../

replace github.com/shengyanli1982/events/contrib/karta => ../contrib/karta

require (
	github.com/shengyanli1982/events v0.0.0-00010101000000-000000000000
	github.com/shengyanli1982/events/contrib/karta v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/shengyanli1982/gs v0.1.6 // indirect
	github.com/shengyanli1982/karta/v2 v2.0.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
