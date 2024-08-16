module github.com/oneliang/frame-golang/test

go 1.22.2

require (
	github.com/oneliang/util-golang/base v0.0.0-20240807054646-8b8ed6868f35 // indirect
	github.com/oneliang/util-golang/common v0.0.0-20240816043657-0a3706677929 // indirect
	github.com/oneliang/util-golang/constants v0.0.0-20240807054646-8b8ed6868f35 // indirect
	github.com/oneliang/util-golang/logging v0.0.0-20240807054646-8b8ed6868f35 // indirect
)

replace (
	github.com/oneliang/frame-golang/context v0.0.0 => ../context
	github.com/oneliang/frame-golang/http v0.0.0 => ./../http
	github.com/oneliang/frame-golang/http/action v0.0.0 => ./../http/action
	github.com/oneliang/frame-golang/query v0.0.0 => ./../query
	github.com/oneliang/frame-golang/ioc v0.0.0 => ./../ioc
)
