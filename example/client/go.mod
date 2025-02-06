module github.com/vdaas/vald/example/client

go 1.23.6

replace (
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v1.2.1
	github.com/goccy/go-json => github.com/goccy/go-json v0.10.5
	github.com/golang/protobuf => github.com/golang/protobuf v1.5.4
	github.com/kpango/glg => github.com/kpango/glg v1.6.15
	github.com/pkg/sftp => github.com/pkg/sftp v1.13.7
	golang.org/x/crypto => golang.org/x/crypto v0.32.0
	golang.org/x/net => golang.org/x/net v0.34.0
	golang.org/x/text => golang.org/x/text v0.22.0
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20250204164813-702378808489
	google.golang.org/genproto/googleapis/api => google.golang.org/genproto/googleapis/api v0.0.0-20250204164813-702378808489
	google.golang.org/genproto/googleapis/rpc => google.golang.org/genproto/googleapis/rpc v0.0.0-20250204164813-702378808489
	google.golang.org/grpc => google.golang.org/grpc v1.70.0
	google.golang.org/protobuf => google.golang.org/protobuf v1.36.4
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.1
	sigs.k8s.io/yaml => sigs.k8s.io/yaml v1.4.0
)

require (
	github.com/kpango/fuid v0.0.0-20221203053508-503b5ad89aa1
	github.com/kpango/glg v1.6.14
	github.com/vdaas/vald-client-go v1.7.16
	gonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc6946
	google.golang.org/grpc v1.67.1
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.1-20241127180247-a33202765966.1 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/kpango/fastime v1.1.9 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241202173237-19429a94021a // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250124145028-65684f501c47 // indirect
	google.golang.org/protobuf v1.36.4 // indirect
)
