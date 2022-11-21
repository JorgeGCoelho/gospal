module github.com/jujuyuki/gospal

go 1.14

require (
	github.com/fatih/color v1.7.0
	github.com/jujuyuki/migo/v3 v3.0.4
	github.com/mattn/go-colorable v0.0.9 // indirect
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.3.0 // indirect
	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.9.1
	golang.org/x/sys v0.0.0-20190109145017-48ac38b7c8cb // indirect
	golang.org/x/tools v0.0.0-20190110163146-51295c7ec13a
)

replace github.com/jujuyuki/migo/v3 => ../gospal-migo/
