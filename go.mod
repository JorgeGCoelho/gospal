module github.com/jujuyuki/gospal

go 1.14

require (
	github.com/fatih/color v1.13.0
	github.com/jujuyuki/migo/v3 v3.0.4
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/pkg/errors v0.9.1
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.23.0
	golang.org/x/tools v0.3.0
)

replace github.com/jujuyuki/migo/v3 => ../gospal-migo/
