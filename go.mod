module software.sslmate.com/src/certspotter

go 1.26.0

require (
	go.uber.org/zap v1.28.0
	golang.org/x/crypto v0.53.0
	golang.org/x/net v0.55.0
	golang.org/x/sync v0.21.0
)

require (
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/text v0.38.0 // indirect
)

retract v0.19.0 // Contains serious bugs.
