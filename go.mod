module github.com/twistedogic/doom

go 1.12

require (
	github.com/alecthomas/jsonschema v0.0.0-20190626084004-00dfc6288dec
	github.com/fatih/structs v1.1.0
	github.com/google/go-cmp v0.3.1
	github.com/iancoleman/strcase v0.0.0-20190422225806-e506e3ef7365
	github.com/json-iterator/go v1.1.7
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/prometheus/client_golang v1.1.0
	github.com/robfig/cron v1.2.0
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/timshannon/bolthold v0.0.0-20190715185903-b73eaf0ecf37
	github.com/twistedogic/jsonpath v0.0.0-20190817133755-3ea90c63ec25
	github.com/uber-go/atomic v1.4.0 // indirect
	github.com/urfave/cli v1.21.0
	go.etcd.io/bbolt v1.3.3
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/ratelimit v0.0.0-20180316092928-c15da0234277
	golang.org/x/net v0.0.0-20190613194153-d28f0bde5980
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	google.golang.org/api v0.8.0
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace go.uber.org/ratelimit v0.0.0-20180316092928-c15da0234277 => github.com/uber-go/ratelimit v0.0.0-20180316092928-c15da0234277
