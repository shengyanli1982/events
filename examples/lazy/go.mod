module github.com/shengyanli1982/events/examples/lazy

go 1.19

replace github.com/shengyanli1982/events => ../../

replace github.com/shengyanli1982/events/contrib/lazy => ../../contrib/lazy

require github.com/shengyanli1982/events/contrib/lazy v0.0.0-00010101000000-000000000000

require (
	github.com/shengyanli1982/events v0.0.0-00010101000000-000000000000 // indirect
	github.com/shengyanli1982/karta v0.2.4 // indirect
	github.com/shengyanli1982/workqueue/v2 v2.2.4 // indirect
	golang.org/x/time v0.5.0 // indirect
)
