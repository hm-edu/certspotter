module software.sslmate.com/src/certspotter

go 1.24

require (
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.37.0
	golang.org/x/net v0.39.0
	golang.org/x/sync v0.13.0
)

require (
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/text v0.24.0 // indirect
)

retract v0.19.0 // Contains serious bugs.
