module github.com/vdaas/vald

go 1.15

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.0+incompatible
	github.com/aws/aws-sdk-go => github.com/aws/aws-sdk-go v1.34.7
	github.com/boltdb/bolt => github.com/boltdb/bolt v1.3.1
	github.com/cockroachdb/errors => github.com/cockroachdb/errors v1.7.4-0.20200821164133-7540597692ae
	github.com/coreos/etcd => go.etcd.io/etcd v0.0.0-20200520232829-54ba9589114f
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v0.4.1
	github.com/go-sql-driver/mysql => github.com/go-sql-driver/mysql v1.5.1-0.20200818111213-46351a889297
	github.com/gocql/gocql => github.com/gocql/gocql v0.0.0-20200815110948-5378c8f664e9
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2-0.20200807193113-deb6fe8ca7c6
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.0
	github.com/gophercloud/gophercloud => github.com/gophercloud/gophercloud v0.12.0
	github.com/gorilla/mux => github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.3-0.20200822210332-78ab81e2420a
	github.com/tensorflow/tensorflow => github.com/tensorflow/tensorflow v2.1.0+incompatible
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	k8s.io/api => k8s.io/api v0.18.8
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.8
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.8
	k8s.io/client-go => k8s.io/client-go v0.18.8
	k8s.io/metrics => k8s.io/metrics v0.18.8
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.6.2
)

require (
	cloud.google.com/go v0.64.0
	code.cloudfoundry.org/bytefmt v0.0.0-20200131002437-cf55d5288a48
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/prometheus v0.2.0
	contrib.go.opencensus.io/exporter/stackdriver v0.13.4
	github.com/aws/aws-sdk-go v1.23.20
	github.com/cespare/xxhash/v2 v2.1.1
	github.com/cockroachdb/errors v1.6.1
	github.com/danielvladco/go-proto-gql/pb v0.6.1
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-redis/redis/v7 v7.4.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gocql/gocql v0.0.0-20200131111108-92af2e088537
	github.com/gocraft/dbr/v2 v2.7.0
	github.com/gogo/protobuf v1.3.1
	github.com/google/go-cmp v0.5.2
	github.com/google/gofuzz v1.2.0
	github.com/gorilla/mux v1.7.1
	github.com/hashicorp/go-version v1.2.1
	github.com/json-iterator/go v1.1.10
	github.com/klauspost/compress v1.10.11
	github.com/kpango/fastime v1.0.16
	github.com/kpango/fuid v0.0.0-20200823100533-287aa95e0641
	github.com/kpango/gache v1.2.3
	github.com/kpango/glg v1.5.1
	github.com/lucasb-eyer/go-colorful v1.0.3
	github.com/pierrec/lz4/v3 v3.3.2
	github.com/scylladb/gocqlx v1.5.0
	github.com/tensorflow/tensorflow v0.0.0-00010101000000-000000000000
	github.com/yahoojapan/gongt v0.0.0-20190517050727-966dcc7aa5e8
	github.com/yahoojapan/ngtd v0.0.0-20200424071638-9872bbae3700
	go.opencensus.io v0.22.4
	go.uber.org/automaxprocs v1.3.0
	go.uber.org/goleak v1.1.10
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6
	gonum.org/v1/hdf5 v0.0.0-20200504100616-496fefe91614
	gonum.org/v1/netlib v0.0.0-20200824093956-f0ca4b3a5ef5 // indirect
	gonum.org/v1/plot v0.7.0
	google.golang.org/api v0.30.0
	google.golang.org/genproto v0.0.0-20200825200019-8632dd797987
	google.golang.org/grpc v1.31.1
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
	k8s.io/metrics v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.0.0-00010101000000-000000000000
)
