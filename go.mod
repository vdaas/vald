module github.com/vdaas/vald

go 1.14

replace (
	github.com/boltdb/bolt => github.com/boltdb/bolt v1.3.1
	github.com/go-sql-driver/mysql => github.com/go-sql-driver/mysql v1.5.1-0.20200311113236-681ffa848bae
	github.com/gocql/gocql => github.com/gocql/gocql v0.0.0-20200324094621-6d895e38b0a5
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.1
	github.com/gophercloud/gophercloud => github.com/gophercloud/gophercloud v0.9.0
	github.com/gorilla/mux => github.com/gorilla/mux v1.7.4
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20200323165209-0ec3e9974c59
	k8s.io/api => k8s.io/api v0.17.4
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.4
	k8s.io/client-go => k8s.io/client-go v0.17.4
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.5.2
)

require (
	bou.ke/monkey v1.0.2 // indirect
	contrib.go.opencensus.io/exporter/jaeger v0.2.0
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	github.com/certifi/gocertifi v0.0.0-20200211180108-c7c1fbc02894 // indirect
	github.com/cespare/xxhash/v2 v2.1.1
	github.com/cockroachdb/errors v1.2.4
	github.com/cockroachdb/logtags v0.0.0-20190617123548-eb05cc24525f // indirect
	github.com/danielvladco/go-proto-gql/pb v0.6.1
	github.com/envoyproxy/protoc-gen-validate v0.3.0
	github.com/getsentry/raven-go v0.2.0 // indirect
	github.com/go-redis/redis/v7 v7.2.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gocql/gocql v0.0.0-20200131111108-92af2e088537
	github.com/gocraft/dbr/v2 v2.7.0
	github.com/gogo/protobuf v1.3.1
	github.com/google/gofuzz v1.1.0
	github.com/gorilla/mux v1.7.1
	github.com/hashicorp/go-version v1.2.0
	github.com/json-iterator/go v1.1.9
	github.com/klauspost/compress v1.10.3
	github.com/kpango/fastime v1.0.16
	github.com/kpango/fuid v0.0.0-20190507064958-80435564606b
	github.com/kpango/gache v1.2.0
	github.com/kpango/glg v1.5.1
	github.com/lucasb-eyer/go-colorful v1.0.3
	github.com/pierrec/lz4/v3 v3.2.1
	github.com/scylladb/gocqlx v1.4.0
	github.com/tensorflow/tensorflow v2.1.0+incompatible
	github.com/yahoojapan/gongt v0.0.0-20190517050727-966dcc7aa5e8
	github.com/yahoojapan/ngtd v0.0.0-20190510080733-0c37ddc5e720
	go.opencensus.io v0.22.3
	go.uber.org/automaxprocs v1.3.0
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys v0.0.0-20191128015809-6d18c012aee9
	golang.org/x/tools v0.0.0-20200414032229-332987a829c3 // indirect
	gonum.org/v1/hdf5 v0.0.0-20191105085658-fe04b73f3b53
	gonum.org/v1/plot v0.7.0
	google.golang.org/genproto v0.0.0-20200401122417-09ab7b7031d2
	google.golang.org/grpc v1.28.0
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.18.0
	k8s.io/apimachinery v0.18.0
	k8s.io/client-go v0.18.0
	k8s.io/metrics v0.18.0
	sigs.k8s.io/controller-runtime v0.0.0-00010101000000-000000000000
)
