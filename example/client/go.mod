module github.com/vdaas/vald/example/client

go 1.24.2

replace (
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v1.2.1
	github.com/goccy/go-json => github.com/goccy/go-json v0.10.5
	github.com/golang/protobuf => github.com/golang/protobuf v1.5.4
	github.com/kpango/glg => github.com/kpango/glg v1.6.15
	github.com/pkg/sftp => github.com/pkg/sftp v1.13.9
	golang.org/x/crypto => golang.org/x/crypto v0.37.0
	golang.org/x/net => golang.org/x/net v0.38.0
<<<<<<< HEAD
	golang.org/x/text => golang.org/x/text v0.23.0
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20250324211829-b45e905df463
	google.golang.org/genproto/googleapis/api => google.golang.org/genproto/googleapis/api v0.0.0-20250324211829-b45e905df463
	google.golang.org/genproto/googleapis/rpc => google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463
=======
	golang.org/x/text => golang.org/x/text v0.24.0
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20250404141209-ee84b53bf3d0
	google.golang.org/genproto/googleapis/api => google.golang.org/genproto/googleapis/api v0.0.0-20250404141209-ee84b53bf3d0
	google.golang.org/genproto/googleapis/rpc => google.golang.org/genproto/googleapis/rpc v0.0.0-20250404141209-ee84b53bf3d0
>>>>>>> 96661e014 (Refactor: add Health Check for Range over gRPC Connection Loop (#2924))
	google.golang.org/grpc => google.golang.org/grpc v1.71.1
	google.golang.org/protobuf => google.golang.org/protobuf v1.36.6
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.1
	sigs.k8s.io/yaml => sigs.k8s.io/yaml v1.4.0
)

require (
	github.com/kpango/fuid v0.0.0-20221203053508-503b5ad89aa1
	github.com/kpango/glg v1.6.14
	github.com/vdaas/vald-client-go v1.7.16
	gonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc6946
	google.golang.org/grpc v1.71.0
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.1-20241127180247-a33202765966.1 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/kpango/fastime v1.1.9 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250106144421-5f5ef82da422 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250313205543-e70fdf4c4cb4 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
