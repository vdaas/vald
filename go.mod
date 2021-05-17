module github.com/vdaas/vald

go 1.16

replace (
	cloud.google.com/go => cloud.google.com/go v0.81.1-0.20210512214559-2ec191233ccc
	cloud.google.com/go/storage => cloud.google.com/go/storage v1.15.1-0.20210512214559-2ec191233ccc
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.1-0.20210210161804-c7f947c0610d+incompatible
	github.com/aws/aws-sdk-go => github.com/aws/aws-sdk-go v1.38.39
	github.com/boltdb/bolt => github.com/boltdb/bolt v1.3.1
	github.com/chzyer/logex => github.com/chzyer/logex v1.1.11-0.20170329064859-445be9e134b2
	github.com/coreos/etcd => go.etcd.io/etcd v3.3.25+incompatible
	github.com/docker/docker => github.com/moby/moby v20.10.6+incompatible
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v0.6.1
	github.com/go-sql-driver/mysql => github.com/go-sql-driver/mysql v1.6.0
	github.com/gocql/gocql => github.com/gocql/gocql v0.0.0-20210504150947-558dfae50b5d
	github.com/gogo/googleapis => github.com/gogo/googleapis v1.4.1
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
	github.com/google/go-cmp => github.com/google/go-cmp v0.5.5
	github.com/google/pprof => github.com/google/pprof v0.0.0-20210506205249-923b5ab0fc1a
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.5.5
	github.com/gophercloud/gophercloud => github.com/gophercloud/gophercloud v0.17.0
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
	github.com/hailocab/go-hostpool => github.com/kpango/go-hostpool v0.0.0-20210303030322-aab80263dcd0
	github.com/klauspost/compress => github.com/klauspost/compress v1.12.3-0.20210511091518-2748482b33e2
	github.com/kpango/glg => github.com/kpango/glg v1.5.1
	github.com/tensorflow/tensorflow => github.com/tensorflow/tensorflow v2.1.2+incompatible
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20210506145944-38f3c27a63bf
	google.golang.org/grpc => google.golang.org/grpc v1.37.1
	google.golang.org/protobuf => google.golang.org/protobuf v1.26.0
	k8s.io/api => k8s.io/api v0.20.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.20.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.6
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.20.6
	k8s.io/client-go => k8s.io/client-go v0.20.6
	k8s.io/metrics => k8s.io/metrics v0.20.6
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.8.3
)

require (
	cloud.google.com/go v0.81.0
	cloud.google.com/go/storage v1.15.0
	code.cloudfoundry.org/bytefmt v0.0.0-20200131002437-cf55d5288a48
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/prometheus v0.3.0
	contrib.go.opencensus.io/exporter/stackdriver v0.13.6
	github.com/aws/aws-sdk-go v1.38.35
	github.com/cespare/xxhash/v2 v2.1.1
	github.com/envoyproxy/protoc-gen-validate v0.1.0
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-redis/redis/v8 v8.8.2
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gocql/gocql v0.0.0-20200131111108-92af2e088537
	github.com/gocraft/dbr/v2 v2.7.1
	github.com/gogo/googleapis v0.0.0-20180223154316-0cd9801be74a
	github.com/gogo/protobuf v1.3.2
	github.com/gogo/status v1.1.0
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.5
	github.com/google/gofuzz v1.2.0
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-version v1.3.0
	github.com/json-iterator/go v1.1.11
	github.com/klauspost/compress v1.12.2
	github.com/kpango/fastime v1.0.16
	github.com/kpango/fuid v0.0.0-20210407064122-2990e29e1ea5
	github.com/kpango/gache v1.2.5
	github.com/kpango/glg v1.5.4
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/pierrec/lz4/v3 v3.3.2
	github.com/scylladb/gocqlx v1.5.0
	github.com/tensorflow/tensorflow v0.0.0-00010101000000-000000000000
	github.com/yahoojapan/gongt v0.0.0-20190517050727-966dcc7aa5e8
	github.com/yahoojapan/ngtd v0.0.0-20200424071638-9872bbae3700
	go.opencensus.io v0.23.0
	go.uber.org/automaxprocs v1.4.0
	go.uber.org/goleak v1.1.10
	go.uber.org/zap v1.16.0
	gocloud.dev v0.23.0
	golang.org/x/net v0.0.0-20210510120150-4163338589ed
	golang.org/x/oauth2 v0.0.0-20210427180440-81ed05c6b58c
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210511113859-b0526f3d8744
	gonum.org/v1/hdf5 v0.0.0-20200504100616-496fefe91614
	gonum.org/v1/plot v0.9.0
	google.golang.org/api v0.46.0
	google.golang.org/genproto v0.0.0-20210510173355-fb37daa5cd7a
	google.golang.org/grpc v1.37.0
	gopkg.in/yaml.v2 v2.4.0
	inet.af/netaddr v0.0.0-20210511181906-37180328850c
	k8s.io/api v0.20.6
	k8s.io/apimachinery v0.20.6
	k8s.io/cli-runtime v0.0.0-00010101000000-000000000000
	k8s.io/client-go v0.20.6
	k8s.io/metrics v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.0.0-00010101000000-000000000000
)
