module github.com/vdaas/vald

go 1.13

replace (
	github.com/DataDog/zstd => github.com/DataDog/zstd v1.3.0
	github.com/boltdb/bolt => github.com/boltdb/bolt v1.3.1
	github.com/go-sql-driver/mysql => github.com/go-sql-driver/mysql v1.5.0
	github.com/gocql/gocql => github.com/gocql/gocql v0.0.0-20200203083758-81b8263d9fe5
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.1
	github.com/gophercloud/gophercloud => github.com/gophercloud/gophercloud v0.7.0
	github.com/gorilla/mux => github.com/gorilla/mux v1.7.3
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20200210222208-86ce3cb69678
	k8s.io/api => k8s.io/api v0.17.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.2
	k8s.io/client-go => k8s.io/client-go v0.17.2
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.4.0
)

require (
	bou.ke/monkey v1.0.2 // indirect
	github.com/DataDog/zstd v0.0.0-00010101000000-000000000000
	github.com/certifi/gocertifi v0.0.0-20200104152315-a6d78f326758 // indirect
	github.com/cespare/xxhash/v2 v2.1.1
	github.com/cockroachdb/errors v1.2.4
	github.com/cockroachdb/logtags v0.0.0-20190617123548-eb05cc24525f // indirect
	github.com/danielvladco/go-proto-gql/pb v0.6.1
	github.com/envoyproxy/protoc-gen-validate v0.1.0
	github.com/getsentry/raven-go v0.2.0 // indirect
	github.com/go-redis/redis/v7 v7.0.0-beta.6
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gocql/gocql v0.0.0-20200103014340-68f928edb90a
	github.com/gocraft/dbr/v2 v2.6.3
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
	github.com/google/gofuzz v1.1.0
	github.com/gorilla/mux v1.7.1
	github.com/hashicorp/go-version v1.2.0
	github.com/json-iterator/go v1.1.9
	github.com/kpango/fastime v1.0.16
	github.com/kpango/fuid v0.0.0-20190507064958-80435564606b
	github.com/kpango/gache v1.2.0
	github.com/kpango/glg v1.5.0
	github.com/lucasb-eyer/go-colorful v1.0.3
	github.com/pierrec/lz4/v3 v3.2.1
	github.com/scylladb/gocqlx v1.3.3
	github.com/tensorflow/tensorflow v2.1.0+incompatible
	github.com/valyala/gozstd v1.6.4
	github.com/yahoojapan/gongt v0.0.0-20190517050727-966dcc7aa5e8
	github.com/yahoojapan/ngtd v0.0.0-20190510080733-0c37ddc5e720
	go.uber.org/automaxprocs v1.3.0
	golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5
	gonum.org/v1/hdf5 v0.0.0-20191105085658-fe04b73f3b53
	gonum.org/v1/plot v0.0.0-20200111075622-4abb28f724d5
	google.golang.org/genproto v0.0.0-20200211035748-55294c81d784
	google.golang.org/grpc v1.27.1
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
	k8s.io/metrics v0.17.2
	sigs.k8s.io/controller-runtime v0.0.0-00010101000000-000000000000
)
