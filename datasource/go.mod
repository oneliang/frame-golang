module github.com/oneliang/frame-golang/datasource

go 1.22.2

require (
	github.com/jmoiron/sqlx v1.4.0
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/shopsprint/decimal v1.3.3
)

require (
	github.com/ClickHouse/ch-go v0.61.5 // indirect
	github.com/ClickHouse/clickhouse-go/v2 v2.30.0 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/oneliang/util-golang/base v0.0.0-20241218054501-34db8a98816d // indirect
	github.com/oneliang/util-golang/common v0.0.0-20241218054501-34db8a98816d // indirect
	github.com/oneliang/util-golang/constants v0.0.0-20241218054501-34db8a98816d // indirect
	github.com/oneliang/util-golang/logging v0.0.0-20241218054501-34db8a98816d // indirect
	github.com/paulmach/orb v0.11.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	go.opentelemetry.io/otel v1.26.0 // indirect
	go.opentelemetry.io/otel/trace v1.26.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/oneliang/frame-golang/context v0.0.0 => ../context
	github.com/oneliang/frame-golang/http v0.0.0 => ./../http
	github.com/oneliang/frame-golang/http/action v0.0.0 => ./../http/action
	github.com/oneliang/frame-golang/ioc v0.0.0 => ./../ioc
	github.com/oneliang/frame-golang/query v0.0.0 => ./../query
)
