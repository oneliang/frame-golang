module github.com/oneliang/frame-golang/transfer

go 1.22.2

require (
	github.com/jmoiron/sqlx v1.4.0
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/shopsprint/decimal v1.3.3
)

require (
	github.com/oneliang/util-golang/base v0.0.0-20250107091240-ffda63863297 // indirect
	github.com/oneliang/util-golang/common v0.0.0-20250107091240-ffda63863297 // indirect
	github.com/oneliang/util-golang/constants v0.0.0-20250107091240-ffda63863297 // indirect
	github.com/oneliang/util-golang/logging v0.0.0-20241218054501-34db8a98816d // indirect
)

replace (
	github.com/oneliang/frame-golang/context v0.0.0 => ../context
	github.com/oneliang/frame-golang/datasource v0.0.0 => ./../datasource
	github.com/oneliang/frame-golang/http v0.0.0 => ./../http
	github.com/oneliang/frame-golang/http/action v0.0.0 => ./../http/action
	github.com/oneliang/frame-golang/ioc v0.0.0 => ./../ioc
	github.com/oneliang/frame-golang/query v0.0.0 => ./../query
	github.com/oneliang/frame-golang/parallel v0.0.0 => ./../parallel
)
