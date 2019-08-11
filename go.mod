module github.com/twistedogic/doom

go 1.12

require (
	github.com/fatih/structs v1.1.0
	github.com/google/go-cmp v0.3.0
	github.com/iancoleman/strcase v0.0.0-20190422225806-e506e3ef7365
	github.com/mitchellh/mapstructure v1.1.2
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/prometheus/client_golang v1.1.0
	github.com/robfig/cron v1.2.0
	github.com/timshannon/bolthold v0.0.0-20190715185903-b73eaf0ecf37
	github.com/twistedogic/jsonpath v0.0.0-20190721104144-ea2b063cd8af
	github.com/uber-go/atomic v1.4.0 // indirect
	github.com/urfave/cli v1.21.0
	go.etcd.io/bbolt v1.3.3
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/ratelimit v0.0.0-20180316092928-c15da0234277
)

replace go.uber.org/ratelimit v0.0.0-20180316092928-c15da0234277 => github.com/uber-go/ratelimit v0.0.0-20180316092928-c15da0234277
