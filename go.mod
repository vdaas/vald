module github.com/vdaas/vald

go 1.13

replace (
	github.com/boltdb/bolt => github.com/boltdb/bolt v1.3.1
	github.com/go-sql-driver/mysql => github.com/go-sql-driver/mysql v1.4.1-0.20191121062641-15462c1d60d4
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20191128160524-b544559bb6d1
	k8s.io/api => k8s.io/api v0.0.0-20191114100352-16d7abae0d2a
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20191114105449-027877536833
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191028221656-72ed19daf4bb
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191114101535-6c5935290e33
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.4.0
)

require (
	github.com/certifi/gocertifi v0.0.0-20191021191039-0944d244cd40 // indirect
	github.com/cespare/xxhash v1.1.0
	github.com/cockroachdb/errors v1.2.4
	github.com/cockroachdb/logtags v0.0.0-20190617123548-eb05cc24525f // indirect
	github.com/danielvladco/go-proto-gql/pb v0.6.1
	github.com/envoyproxy/protoc-gen-validate v0.1.0
	github.com/getsentry/raven-go v0.2.0 // indirect
	github.com/go-redis/redis/v7 v7.0.0-beta.4
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gocql/gocql v0.0.0-20191126110522-1982a06ad6b9
	github.com/gocraft/dbr/v2 v2.6.3
	github.com/gogo/protobuf v1.3.1
	github.com/google/gofuzz v1.0.0
	github.com/gorilla/mux v1.7.3
	github.com/hashicorp/go-version v1.2.0
	github.com/json-iterator/go v1.1.8
	github.com/kpango/fastime v1.0.15
	github.com/kpango/fuid v0.0.0-20190507064958-80435564606b
	github.com/kpango/gache v1.1.23
	github.com/kpango/glg v1.4.6
	github.com/scylladb/gocqlx v1.3.1
	github.com/yahoojapan/gongt v0.0.0-20190517050727-966dcc7aa5e8
	github.com/yahoojapan/ngtd v0.0.0-20190510080733-0c37ddc5e720
	golang.org/x/sys v0.0.0-20191128015809-6d18c012aee9
	gonum.org/v1/hdf5 v0.0.0-20191105085658-fe04b73f3b53
	google.golang.org/genproto v0.0.0-20191115221424-83cc0476cb11
	google.golang.org/grpc v1.25.1
	gopkg.in/yaml.v2 v2.2.7
	k8s.io/api v0.0.0-20191114100352-16d7abae0d2a
	k8s.io/apimachinery v0.0.0-20191028221656-72ed19daf4bb
	k8s.io/client-go v0.0.0-20191114101535-6c5935290e33
	sigs.k8s.io/controller-runtime v0.0.0-00010101000000-000000000000
)
