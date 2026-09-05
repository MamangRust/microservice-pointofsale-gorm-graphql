module github.com/MamangRust/microservice-point-of-sale-shared

go 1.25.0

require (
	github.com/MamangRust/microservice-point-of-sale-pkg v1.0.5
	github.com/MamangRust/microservice-pointofsale-grpc/pb v0.0.0-00010101000000-000000000000
	github.com/go-playground/validator/v10 v10.30.1
	github.com/labstack/echo/v4 v4.15.0
	github.com/redis/go-redis/v9 v9.21.0
	go.opentelemetry.io/otel v1.46.0
	go.opentelemetry.io/otel/metric v1.46.0
	go.opentelemetry.io/otel/trace v1.46.0
	go.uber.org/zap v1.28.0
	google.golang.org/grpc v1.83.1
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.13 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.55.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260819154853-08b0e4226688 // indirect
)

replace github.com/MamangRust/microservice-pointofsale-grpc/pb => ../pb

replace github.com/MamangRust/microservice-point-of-sale-pkg => ../pkg
