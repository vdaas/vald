module github.com/vdaas/vald/example/client

go 1.27.4

replace (
	github.com/kpango/fuid => github.com/kpango/fuid v0.0.0-20221203053508-503b5ad89aa1
	github.com/kpango/glg => github.com/kpango/glg v1.6.15
	github.com/vdaas/vald-client-go => github.com/vdaas/vald-client-go v1.8.1
	golang.org/x/net => golang.org/x/net v0.62.0
	gonum.org/v1/hdf5 => gonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc6946
	google.golang.org/grpc => google.golang.org/grpc v1.85.0
)

require (
	github.com/kpango/fuid v0.0.0-00010101000000-000000000000
	github.com/kpango/glg v1.6.14
	github.com/vdaas/vald-client-go v1.8.1
	gonum.org/v1/hdf5 v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.83.1
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.12-20260825204119-511051f7f437.2 // indirect
	github.com/goccy/go-json v0.12.0 // indirect
	github.com/kpango/fastime v1.1.10 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	golang.org/x/net v0.60.0 // indirect
	golang.org/x/sys v0.50.0 // indirect
	golang.org/x/text v0.44.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20261110000943-09e21b237f2c // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261110000943-09e21b237f2c // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)
