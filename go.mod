module github.com/vdaas/vald

go 1.16

replace (
	cloud.google.com/go => cloud.google.com/go v0.81.0
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.1-0.20210210161804-c7f947c0610d+incompatible
	github.com/aws/aws-sdk-go => github.com/aws/aws-sdk-go v1.38.13
	github.com/boltdb/bolt => github.com/boltdb/bolt v1.3.1
	github.com/chzyer/logex => github.com/chzyer/logex v1.1.11-0.20170329064859-445be9e134b2
	github.com/coreos/etcd => go.etcd.io/etcd v3.3.25+incompatible
	github.com/docker/docker => github.com/moby/moby v20.10.5+incompatible
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v0.5.1
	github.com/go-sql-driver/mysql => github.com/go-sql-driver/mysql v1.6.0
	github.com/gocql/gocql => github.com/gocql/gocql v0.0.0-20210401103645-80ab1e13e309
	github.com/gogo/googleapis => github.com/gogo/googleapis v1.4.1
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
	github.com/google/go-cmp => github.com/google/go-cmp v0.5.5
	github.com/google/pprof => github.com/google/pprof v0.0.0-20210323184331-8eee2492667d
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.5.4
	github.com/gophercloud/gophercloud => github.com/gophercloud/gophercloud v0.16.0
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
	github.com/hailocab/go-hostpool => github.com/kpango/go-hostpool v0.0.0-20210303030322-aab80263dcd0
	github.com/klauspost/compress => github.com/klauspost/compress v1.11.14-0.20210330161112-b2eb8368e62b
	github.com/kpango/glg => github.com/kpango/glg v1.5.1
	github.com/tensorflow/tensorflow => github.com/tensorflow/tensorflow v2.1.2+incompatible
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	google.golang.org/grpc => google.golang.org/grpc v1.36.1
	google.golang.org/protobuf => google.golang.org/protobuf v1.26.0
	k8s.io/api => k8s.io/api v0.20.5
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.20.5
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.5
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.20.5
	k8s.io/client-go => k8s.io/client-go v0.20.5
	k8s.io/metrics => k8s.io/metrics v0.20.5
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.8.3
)

require (
	cloud.google.com/go v0.79.0
	cloud.google.com/go/storage v1.12.0 // indirect
	code.cloudfoundry.org/bytefmt v0.0.0-20200131002437-cf55d5288a48
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.4.12 // indirect
	contrib.go.opencensus.io/exporter/prometheus v0.3.0
	contrib.go.opencensus.io/exporter/stackdriver v0.13.5
	contrib.go.opencensus.io/resource v0.0.0-20190131005048-21591786a5e0 // indirect
	github.com/Azure/azure-amqp-common-go v1.1.4 // indirect
	github.com/aws/aws-sdk-go v1.36.1
	github.com/cespare/xxhash/v2 v2.1.1
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/envoyproxy/protoc-gen-validate v0.1.0
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-redis/redis/v8 v8.8.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gocql/gocql v0.0.0-20200131111108-92af2e088537
	github.com/gocraft/dbr/v2 v2.7.1
	github.com/gogo/googleapis v0.0.0-20180223154316-0cd9801be74a
	github.com/gogo/protobuf v1.3.2
	github.com/gogo/status v1.1.0
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.5
	github.com/google/gofuzz v1.2.0
	github.com/googleapis/gax-go v2.0.2+incompatible // indirect
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-version v1.3.0
	github.com/json-iterator/go v1.1.10
	github.com/klauspost/compress v1.11.3
	github.com/kpango/fastime v1.0.16
	github.com/kpango/fuid v0.0.0-20200823100533-287aa95e0641
	github.com/kpango/gache v1.2.5
	github.com/kpango/glg v1.5.4
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/pierrec/lz4/v3 v3.3.2
	github.com/scylladb/gocqlx v1.5.0
	github.com/tensorflow/tensorflow v0.0.0-00010101000000-000000000000
	github.com/tidwall/pretty v0.0.0-20190325153808-1166b9ac2b65 // indirect
	github.com/uber-go/atomic v1.3.2 // indirect
	github.com/uber/jaeger-lib v1.5.0 // indirect
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/xdg/stringprep v1.0.0 // indirect
	github.com/yahoojapan/gongt v0.0.0-20190517050727-966dcc7aa5e8
	github.com/yahoojapan/ngtd v0.0.0-20200424071638-9872bbae3700
	go.mongodb.org/mongo-driver v1.0.1 // indirect
	go.opencensus.io v0.23.0
	go.uber.org/automaxprocs v1.4.0
	go.uber.org/goleak v1.1.10
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210403161142-5e06dd20ab57
	gonum.org/v1/hdf5 v0.0.0-20200504100616-496fefe91614
	gonum.org/v1/plot v0.9.0
	google.golang.org/api v0.43.0
	google.golang.org/genproto v0.0.0-20210405174219-a39eb2f71cb9
	google.golang.org/grpc v1.36.1
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.20.5
	k8s.io/apimachinery v0.20.5
	k8s.io/cli-runtime v0.0.0-00010101000000-000000000000
	k8s.io/client-go v0.20.5
	k8s.io/metrics v0.0.0-00010101000000-000000000000
	pack.ag/amqp v0.11.0 // indirect
	sigs.k8s.io/controller-runtime v0.0.0-00010101000000-000000000000
)
