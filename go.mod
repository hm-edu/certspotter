module software.sslmate.com/src/certspotter

go 1.26.0

require (
	golang.org/x/crypto v0.54.0
	golang.org/x/net v0.57.0
	golang.org/x/sync v0.22.0
	golang.org/x/time v0.15.0
)

require (
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.28.0
	golang.org/x/text v0.40.0 // indirect
)

retract v0.19.0 // Contains serious bugs.
