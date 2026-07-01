module github.com/shengyanli1982/events/examples/lazy

go 1.23

replace github.com/shengyanli1982/events => ../../

replace github.com/shengyanli1982/events/contrib/lazy => ../../contrib/lazy

require github.com/shengyanli1982/events/contrib/lazy v0.0.0-00010101000000-000000000000

require (
	github.com/shengyanli1982/events v0.0.0-00010101000000-000000000000 // indirect
	github.com/shengyanli1982/events/contrib/karta v0.0.0-00010101000000-000000000000 // indirect
	github.com/shengyanli1982/gs v0.1.6 // indirect
	github.com/shengyanli1982/karta/v2 v2.0.0 // indirect
	golang.org/x/time v0.5.0 // indirect
)
