module github.com/shengyanli1982/events/examples/default

go 1.19

replace github.com/shengyanli1982/events => ../../

require (
	github.com/shengyanli1982/events v0.0.0-00010101000000-000000000000
	github.com/shengyanli1982/karta v0.2.1
	github.com/shengyanli1982/workqueue/v2 v2.1.3
)

require golang.org/x/time v0.5.0 // indirect
