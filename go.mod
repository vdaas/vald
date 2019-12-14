module github.com/vdaas/vald

go 1.13

replace (
	github.com/boltdb/bolt => github.com/boltdb/bolt v1.3.1
	github.com/go-sql-driver/mysql => github.com/go-sql-driver/mysql v1.4.1-0.20191212001955-b66d043e6c89
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.1
	github.com/gophercloud/gophercloud => github.com/gophercloud/gophercloud v0.7.0
	github.com/gorilla/mux => github.com/gorilla/mux v1.7.3
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413
	k8s.io/api => k8s.io/api v0.16.4
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.16.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.5-beta.1
	k8s.io/client-go => k8s.io/client-go v0.16.4
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.4.0
)

require (
	github.com/ajstarks/svgo v0.0.0-20191124160048-bd5c74aaa11c // indirect
	github.com/certifi/gocertifi v0.0.0-20191021191039-0944d244cd40 // indirect
	github.com/cespare/xxhash v1.1.0
	github.com/cockroachdb/errors v1.2.4
	github.com/cockroachdb/logtags v0.0.0-20190617123548-eb05cc24525f // indirect
	github.com/danielvladco/go-proto-gql/pb v0.6.1
	github.com/envoyproxy/protoc-gen-validate v0.1.0
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/getsentry/raven-go v0.2.0 // indirect
	github.com/go-redis/redis/v7 v7.0.0-beta.4
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gocql/gocql v0.0.0-20191126110522-1982a06ad6b9
	github.com/gocraft/dbr/v2 v2.6.3
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
	github.com/google/gofuzz v1.0.0
	github.com/gorilla/mux v1.7.1
	github.com/hashicorp/go-version v1.2.0
	github.com/json-iterator/go v1.1.8
	github.com/jung-kurt/gofpdf v1.16.2 // indirect
	github.com/kpango/fastime v1.0.15
	github.com/kpango/fuid v0.0.0-20190507064958-80435564606b
	github.com/kpango/gache v1.1.23
	github.com/kpango/glg v1.4.6
	github.com/scylladb/gocqlx v1.3.1
	github.com/wcharczuk/go-chart v2.0.1+incompatible
	github.com/yahoojapan/gongt v0.0.0-20190517050727-966dcc7aa5e8
	github.com/yahoojapan/ngtd v0.0.0-20190510080733-0c37ddc5e720
	go.uber.org/zap v1.10.0
	golang.org/x/image v0.0.0-20191206065243-da761ea9ff43 // indirect
	golang.org/x/sys v0.0.0-20191206220618-eeba5f6aabab
	gonum.org/v1/hdf5 v0.0.0-20191105085658-fe04b73f3b53
	gonum.org/v1/plot v0.0.0-20191107103940-ca91d9d40d0a
	google.golang.org/genproto v0.0.0-20191206224255-0243a4be9c8f
	google.golang.org/grpc v1.25.1
	gopkg.in/yaml.v2 v2.2.7
	k8s.io/api v0.16.4
	k8s.io/apimachinery v0.16.4
	k8s.io/client-go v0.16.4
	sigs.k8s.io/controller-runtime v0.0.0-00010101000000-000000000000
)
