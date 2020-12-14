module github.com/vdaas/vald

go 1.15

replace (
	cloud.google.com/go => cloud.google.com/go v0.74.0
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.0+incompatible
	github.com/aws/aws-sdk-go => github.com/aws/aws-sdk-go v1.36.7
	github.com/boltdb/bolt => github.com/boltdb/bolt v1.3.1
	github.com/chzyer/logex => github.com/chzyer/logex v1.1.11-0.20170329064859-445be9e134b2
	github.com/coreos/etcd => go.etcd.io/etcd v3.3.25+incompatible
	github.com/docker/docker => github.com/moby/moby v20.10.0+incompatible
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v0.4.1
	github.com/go-sql-driver/mysql => github.com/go-sql-driver/mysql v1.5.0
	github.com/gocql/gocql => github.com/gocql/gocql v0.0.0-20201209090715-f485b5f9159c
	github.com/gogo/googleapis => github.com/gogo/googleapis v1.4.0
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.1
	github.com/google/go-cmp => github.com/google/go-cmp v0.5.4
	github.com/google/pprof => github.com/google/pprof v0.0.0-20201211104106-9bd6f8a8ed4b
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.5.3
	github.com/gophercloud/gophercloud => github.com/gophercloud/gophercloud v0.14.0
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
	github.com/hailocab/go-hostpool => github.com/monzo/go-hostpool v0.0.0-20200724120130-287edbb29340
	github.com/klauspost/compress => github.com/klauspost/compress v1.11.4-0.20201208122001-8c54b4233d2e
	github.com/tensorflow/tensorflow => github.com/tensorflow/tensorflow v2.1.2+incompatible
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20201208171446-5f87f3452ae9
	google.golang.org/grpc => google.golang.org/grpc v1.34.0
	google.golang.org/protobuf => google.golang.org/protobuf v1.25.0
	k8s.io/api => k8s.io/api v0.20.0
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.20.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.0
	k8s.io/client-go => k8s.io/client-go v0.20.0
	k8s.io/metrics => k8s.io/metrics v0.20.0
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.7.0
)

require (
	cloud.google.com/go v0.74.0
	cloud.google.com/go/storage v1.12.0 // indirect
	code.cloudfoundry.org/bytefmt v0.0.0-20200131002437-cf55d5288a48
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/prometheus v0.2.0
	contrib.go.opencensus.io/exporter/stackdriver v0.13.4
	github.com/DATA-DOG/go-sqlmock v1.5.0 // indirect
	github.com/ajstarks/svgo v0.0.0-20200725142600-7a3c8b57fecb // indirect
	github.com/alecthomas/units v0.0.0-20201120081800-1786d5ef83d4 // indirect
	github.com/aws/aws-sdk-go v1.36.7
	github.com/census-instrumentation/opencensus-proto v0.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.1
	github.com/denisenkom/go-mssqldb v0.9.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.4.1
	github.com/frankban/quicktest v1.11.2 // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-latex/latex v0.0.0-20201211173324-01eae8bd88f2 // indirect
	github.com/go-logr/zapr v0.3.0 // indirect
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/go-redis/redis/v8 v8.4.2
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gocql/gocql v0.0.0-20201209090715-f485b5f9159c
	github.com/gocraft/dbr/v2 v2.7.1
	github.com/gogo/protobuf v1.3.1
	github.com/golang/snappy v0.0.2 // indirect
	github.com/google/go-cmp v0.5.4
	github.com/google/gofuzz v1.2.0
	github.com/google/pprof v0.0.0-20201211104106-9bd6f8a8ed4b // indirect
	github.com/googleapis/gnostic v0.5.3 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-version v1.2.1
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/json-iterator/go v1.1.10
	github.com/klauspost/compress v1.11.3
	github.com/kpango/fastime v1.0.16
	github.com/kpango/fuid v0.0.0-20200823100533-287aa95e0641
	github.com/kpango/gache v1.2.3
	github.com/kpango/glg v1.5.1
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.9.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.0.3
	github.com/mattn/go-sqlite3 v1.14.5 // indirect
	github.com/nxadm/tail v1.4.5 // indirect
	github.com/onsi/gomega v1.10.4 // indirect
	github.com/pierrec/lz4/v3 v3.3.2
	github.com/prometheus/client_golang v1.8.0 // indirect
	github.com/prometheus/common v0.15.0 // indirect
	github.com/prometheus/statsd_exporter v0.18.0 // indirect
	github.com/scylladb/gocqlx v1.5.0
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/tensorflow/tensorflow v2.3.1+incompatible
	github.com/yahoojapan/gongt v0.0.0-20190517050727-966dcc7aa5e8
	github.com/yahoojapan/ngtd v0.0.0-20200424071638-9872bbae3700
	go.opencensus.io v0.22.5
	go.opentelemetry.io/otel v0.15.0 // indirect
	go.uber.org/automaxprocs v1.3.0
	go.uber.org/goleak v1.1.10
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20201208171446-5f87f3452ae9 // indirect
	golang.org/x/exp v0.0.0-20201210212021-a20c86df00b4 // indirect
	golang.org/x/net v0.0.0-20201209123823-ac852fbbde11
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	golang.org/x/sys v0.0.0-20201214095126-aec9a390925b
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
	golang.org/x/tools v0.0.0-20201211185031-d93e913c1a58 // indirect
	gonum.org/v1/gonum v0.8.2 // indirect
	gonum.org/v1/hdf5 v0.0.0-20200504100616-496fefe91614
	gonum.org/v1/netlib v0.0.0-20201012070519-2390d26c3658 // indirect
	gonum.org/v1/plot v0.8.1
	google.golang.org/api v0.36.0
	google.golang.org/genproto v0.0.0-20201211151036-40ec1c210f7a
	google.golang.org/grpc v1.34.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0
	honnef.co/go/tools v0.0.1-2020.1.6 // indirect
	k8s.io/api v0.20.0
	k8s.io/apiextensions-apiserver v0.20.0 // indirect
	k8s.io/apimachinery v0.20.0
	k8s.io/client-go v1.5.1
	k8s.io/metrics v0.20.0
	sigs.k8s.io/controller-runtime v0.7.0
)
