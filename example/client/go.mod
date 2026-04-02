module github.com/vdaas/vald/example/client

go 1.26.1

replace (
	github.com/kpango/fuid => github.com/kpango/fuid v0.0.0-20221203053508-503b5ad89aa1
	github.com/kpango/glg => github.com/kpango/glg v1.6.15
	github.com/vdaas/vald-client-go => github.com/vdaas/vald-client-go v1.7.17
	gonum.org/v1/hdf5 => gonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc6946
	google.golang.org/grpc => google.golang.org/grpc v1.80.0
)

require (
	github.com/kpango/fuid v0.0.0-00010101000000-000000000000
	github.com/kpango/glg v1.6.14
	github.com/vdaas/vald-client-go v1.7.17
	gonum.org/v1/hdf5 v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.79.3
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.11-20260209202127-80ab13bee0bf.1 // indirect
	github.com/goccy/go-json v0.10.6 // indirect
	github.com/kpango/fastime v1.1.10 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260401024825-9d38bb4040a9 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260401024825-9d38bb4040a9 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
