module github.com/vdaas/vald

go 1.13

replace (
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20191002192127-34f69633bfdc
	k8s.io/api => k8s.io/api v0.0.0-20190918195907-bd6ac527cfd2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190918201827-3de75813f604
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190918200256-06eb1244587a
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.2.2
)

require (
	bou.ke/monkey v1.0.1 // indirect
	github.com/certifi/gocertifi v0.0.0-20190905060710-a5e0173ced67 // indirect
	github.com/cespare/xxhash v1.1.0
	github.com/cockroachdb/errors v1.2.3
	github.com/cockroachdb/logtags v0.0.0-20190617123548-eb05cc24525f // indirect
	github.com/danielvladco/go-proto-gql/pb v0.6.1
	github.com/envoyproxy/protoc-gen-validate v0.1.0
	github.com/getsentry/raven-go v0.2.0 // indirect
	github.com/go-redis/redis v6.15.5+incompatible
	github.com/gogo/protobuf v1.3.0
	github.com/gorilla/mux v1.7.3
	github.com/hashicorp/go-version v1.2.0
	github.com/json-iterator/go v1.1.7
	github.com/kpango/fastime v1.0.15
	github.com/kpango/fuid v0.0.0-20190507064958-80435564606b
	github.com/kpango/gache v1.1.22
	github.com/kpango/glg v1.4.6
	github.com/yahoojapan/gongt v0.0.0-20190517050727-966dcc7aa5e8
	golang.org/x/sys v0.0.0-20191007092633-5f54ce542709
	gonum.org/v1/hdf5 v0.0.0-20190920010848-b0d662f53d94
	google.golang.org/genproto v0.0.0-20191002211648-c459b9ce5143
	google.golang.org/grpc v1.24.0
	gopkg.in/yaml.v2 v2.2.4
	k8s.io/api v0.0.0-20190918195907-bd6ac527cfd2
	k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	sigs.k8s.io/controller-runtime v0.0.0-00010101000000-000000000000
)
