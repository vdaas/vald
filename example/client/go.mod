module github.com/vdaas/vald/example/client

go 1.22.0

replace (
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v1.0.4
	github.com/goccy/go-json => github.com/goccy/go-json v0.10.2
	github.com/golang/protobuf => github.com/golang/protobuf v1.5.3
	github.com/kpango/gache/v2 => github.com/kpango/gache/v2 v2.0.9
	github.com/kpango/glg => github.com/kpango/glg v1.6.15
	github.com/pkg/sftp => github.com/pkg/sftp v1.13.6
	github.com/vdaas/vald => ../../../vald
	golang.org/x/crypto => golang.org/x/crypto v0.20.0
	golang.org/x/net => golang.org/x/net v0.21.0
	golang.org/x/text => golang.org/x/text v0.14.0
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20240228224816-df926f6c8641
	google.golang.org/genproto/googleapis/api => google.golang.org/genproto/googleapis/api v0.0.0-20240228224816-df926f6c8641
	google.golang.org/genproto/googleapis/rpc => google.golang.org/genproto/googleapis/rpc v0.0.0-20240228224816-df926f6c8641
	google.golang.org/grpc => google.golang.org/grpc v1.62.0
	google.golang.org/protobuf => google.golang.org/protobuf v1.32.0
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/kpango/fuid v0.0.0-20221203053508-503b5ad89aa1
	github.com/kpango/glg v1.6.15
	github.com/vdaas/vald v0.0.0-00010101000000-000000000000
	github.com/vdaas/vald-client-go v1.7.12
	gonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc6946
	google.golang.org/grpc v1.61.1
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.32.0-20240221180331-f05a6f4403ce.1 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/kpango/fastime v1.1.9 // indirect
	github.com/kpango/gache/v2 v2.0.0-00010101000000-000000000000 // indirect
	github.com/planetscale/vtprotobuf v0.6.0 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240213162025-012b6fc9bca9 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240228201840-1f18d85a4ec2 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)
