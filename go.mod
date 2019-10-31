module github.com/vdaas/vald

go 1.13

replace (
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20191029031824-8986dd9e96cf
	k8s.io/api => k8s.io/api v0.0.0-20191016110246-af539daaa43a
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20191016113439-b64f2075a530
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191004115701-31ade1b30762
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191016110837-54936ba21026
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.3.0
)

require (
	github.com/certifi/gocertifi v0.0.0-20191021191039-0944d244cd40 // indirect
	github.com/cespare/xxhash v1.1.0
	github.com/cockroachdb/errors v1.2.3
	github.com/cockroachdb/logtags v0.0.0-20190617123548-eb05cc24525f // indirect
	github.com/danielvladco/go-proto-gql/pb v0.6.1
	github.com/envoyproxy/protoc-gen-validate v0.1.0
	github.com/getsentry/raven-go v0.2.0 // indirect
	github.com/go-redis/redis/v7 v7.0.0-beta.4
	github.com/gogo/protobuf v1.3.1
	github.com/gorilla/mux v1.7.3
	github.com/hashicorp/go-version v1.2.0
	github.com/json-iterator/go v1.1.8
	github.com/kpango/fastime v1.0.15
	github.com/kpango/fuid v0.0.0-20190507064958-80435564606b
	github.com/kpango/gache v1.1.23
	github.com/kpango/glg v1.4.6
	github.com/yahoojapan/gongt v0.0.0-20190517050727-966dcc7aa5e8
	golang.org/x/sys v0.0.0-20191029155521-f43be2a4598c
	gonum.org/v1/hdf5 v0.0.0-20190920010848-b0d662f53d94
	google.golang.org/genproto v0.0.0-20191028173616-919d9bdd9fe6
	google.golang.org/grpc v1.24.0
	gopkg.in/yaml.v2 v2.2.4
	k8s.io/api v0.0.0-20191016110246-af539daaa43a
	k8s.io/apimachinery v0.0.0-20191004115701-31ade1b30762
	k8s.io/client-go v0.0.0-20191016110837-54936ba21026
	sigs.k8s.io/controller-runtime v0.0.0-00010101000000-000000000000
)
