module github.com/vdaas/vald/example/client

go 1.25.3

replace (
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v1.2.1
	github.com/goccy/go-json => github.com/goccy/go-json v0.10.5
	github.com/golang/protobuf => github.com/golang/protobuf v1.5.4
	github.com/kpango/glg => github.com/kpango/glg v1.6.15
	github.com/pkg/sftp => github.com/pkg/sftp v1.13.9
	golang.org/x/crypto => golang.org/x/crypto v0.43.0
	golang.org/x/net => golang.org/x/net v0.46.0
	golang.org/x/text => golang.org/x/text v0.30.0
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20251014184007-4626949a642f
	google.golang.org/genproto/googleapis/api => google.golang.org/genproto/googleapis/api v0.0.0-20251014184007-4626949a642f
	google.golang.org/genproto/googleapis/rpc => google.golang.org/genproto/googleapis/rpc v0.0.0-20251014184007-4626949a642f
	google.golang.org/grpc => google.golang.org/grpc v1.76.0
	google.golang.org/protobuf => google.golang.org/protobuf v1.36.10
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.1
	sigs.k8s.io/yaml => sigs.k8s.io/yaml v1.6.0
)

require (
	github.com/kpango/fuid v0.0.0-20221203053508-503b5ad89aa1
	github.com/kpango/glg v1.6.14
	github.com/vdaas/vald-client-go v1.7.17
	gonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc6946
	google.golang.org/grpc v1.71.0
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.6-20250625184727-c923a0c2a132.1 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/kpango/fastime v1.1.9 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250804133106-a7a43d27e69b // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251007200510-49b9836ed3ff // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
