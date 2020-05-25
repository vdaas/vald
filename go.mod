module github.com/vdaas/vald

go 1.14

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.1.1+incompatible
	github.com/boltdb/bolt => github.com/boltdb/bolt v1.3.1
	github.com/cockroachdb/errors => github.com/cockroachdb/errors v1.2.5-0.20200521081835-9833fdfac54c
	github.com/coreos/etcd => go.etcd.io/etcd v0.5.0-alpha.5.0.20200425165423-262c93980547
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200309214505-aa6a9891b09c+incompatible
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v0.3.0-java
	github.com/go-sql-driver/mysql => github.com/go-sql-driver/mysql v1.5.1-0.20200521132419-128a6737a202
	github.com/gocql/gocql => github.com/gocql/gocql v0.0.0-20200519160334-799061058e31
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.1
	github.com/gophercloud/gophercloud => github.com/gophercloud/gophercloud v0.11.0
	github.com/gorilla/mux => github.com/gorilla/mux v1.7.5-0.20200517040254-948bec34b516
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
	github.com/tensorflow/tensorflow => github.com/tensorflow/tensorflow v2.1.0+incompatible
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
	k8s.io/api => k8s.io/api v0.18.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.2
	k8s.io/client-go => k8s.io/client-go v0.18.2
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.6.0
)

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.0
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	github.com/aws/aws-sdk-go v1.31.4
	github.com/cespare/xxhash/v2 v2.1.1
	github.com/cockroachdb/errors v0.0.0-00010101000000-000000000000
	github.com/danielvladco/go-proto-gql/pb v0.6.1
	github.com/envoyproxy/protoc-gen-validate v0.1.0
	github.com/go-redis/redis/v7 v7.3.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gocql/gocql v0.0.0-20200131111108-92af2e088537
	github.com/gocraft/dbr/v2 v2.7.0
	github.com/gogo/protobuf v1.3.1
	github.com/google/gofuzz v1.1.0
	github.com/gorilla/mux v1.7.1
	github.com/hashicorp/go-version v1.2.0
	github.com/json-iterator/go v1.1.9
	github.com/klauspost/compress v1.10.6
	github.com/kpango/fastime v1.0.16
	github.com/kpango/fuid v0.0.0-20190507064958-80435564606b
	github.com/kpango/gache v1.2.1
	github.com/kpango/glg v1.5.1
	github.com/lucasb-eyer/go-colorful v1.0.3
	github.com/pierrec/lz4/v3 v3.3.2
	github.com/scylladb/gocqlx v1.5.0
	github.com/tensorflow/tensorflow v0.0.0-00010101000000-000000000000
	github.com/yahoojapan/gongt v0.0.0-20190517050727-966dcc7aa5e8
	github.com/yahoojapan/ngtd v0.0.0-20200424071638-9872bbae3700
	go.opencensus.io v0.22.3
	go.uber.org/automaxprocs v1.3.0
	go.uber.org/goleak v1.0.0
	gocloud.dev v0.19.0
	golang.org/x/mod v0.3.0 // indirect
	golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5
	golang.org/x/tools v0.0.0-20200522201501-cb1345f3a375 // indirect
	gonum.org/v1/hdf5 v0.0.0-20200504100616-496fefe91614
	gonum.org/v1/netlib v0.0.0-20200317120129-c5a04cffd98a // indirect
	gonum.org/v1/plot v0.7.0
	google.golang.org/genproto v0.0.0-20200521103424-e9a78aa275b7
	google.golang.org/grpc v1.29.1
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.18.3
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v0.18.3
	k8s.io/metrics v0.18.3
	sigs.k8s.io/controller-runtime v0.0.0-00010101000000-000000000000
)
