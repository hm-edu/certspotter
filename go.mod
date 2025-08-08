module software.sslmate.com/src/certspotter

go 1.24.4

require (
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.41.0
	golang.org/x/net v0.43.0
	golang.org/x/sync v0.16.0
)

require (
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/text v0.28.0 // indirect
)

retract v0.19.0 // Contains serious bugs.
