module github.com/vdaas/vald/example/client

go 1.23.1

replace (
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v1.1.0
	github.com/goccy/go-json => github.com/goccy/go-json v0.10.3
	github.com/golang/protobuf => github.com/golang/protobuf v1.5.4
	github.com/kpango/glg => github.com/kpango/glg v1.6.15
	github.com/pkg/sftp => github.com/pkg/sftp v1.13.6
	golang.org/x/crypto => golang.org/x/crypto v0.27.0
	golang.org/x/net => golang.org/x/net v0.29.0
	golang.org/x/text => golang.org/x/text v0.18.0
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20240903143218-8af14fe29dc1
	google.golang.org/genproto/googleapis/api => google.golang.org/genproto/googleapis/api v0.0.0-20240903143218-8af14fe29dc1
	google.golang.org/genproto/googleapis/rpc => google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1
	google.golang.org/grpc => google.golang.org/grpc v1.66.1
	google.golang.org/protobuf => google.golang.org/protobuf v1.34.2
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/kpango/fuid v0.0.0-20221203053508-503b5ad89aa1
	github.com/kpango/glg v1.6.14
	github.com/vdaas/vald-client-go v1.7.13
	gonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc6946
	google.golang.org/grpc v1.66.0
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.34.2-20240717164558-a6c49f84cc0f.2 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/kpango/fastime v1.1.9 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240604185151-ef581f913117 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240827150818-7e3bb234dfed // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
