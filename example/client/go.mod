module github.com/vdaas/vald/example/client

go 1.18

replace (
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v0.6.7
	github.com/goccy/go-json => github.com/goccy/go-json v0.9.7
	github.com/golang/protobuf => github.com/golang/protobuf v1.5.2
	github.com/kpango/glg => github.com/kpango/glg v1.6.10
	github.com/pkg/sftp => github.com/pkg/sftp v1.13.4
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20220518034528-6f7dac969898
	golang.org/x/net => golang.org/x/net v0.0.0-20220520000938-2e3eb7b945c2
	golang.org/x/text => golang.org/x/text v0.3.7
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20220519153652-3a47de7e79bd
	google.golang.org/grpc => google.golang.org/grpc v1.46.2
	google.golang.org/protobuf => google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.0
)

require (
	github.com/kpango/fuid v0.0.0-20220209050620-e5987ba1ea5e
	github.com/kpango/glg v1.6.10
	github.com/vdaas/vald-client-go v1.5.2
	gonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc6946
	google.golang.org/grpc v1.46.0
)

require (
	github.com/envoyproxy/protoc-gen-validate v0.1.0 // indirect
	github.com/goccy/go-json v0.9.4 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/kpango/fastime v1.1.4 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211104193956-4c6863e31247 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
