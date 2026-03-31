module github.com/vdaas/vald

go 1.26.1

replace (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go => buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.11-20260209202127-80ab13bee0bf.1
	cloud.google.com/go/storage => cloud.google.com/go/storage v1.61.3
	code.cloudfoundry.org/bytefmt => code.cloudfoundry.org/bytefmt v0.67.0
	github.com/akrylysov/pogreb => github.com/akrylysov/pogreb v0.10.2
	github.com/aws/aws-sdk-go => github.com/aws/aws-sdk-go v1.55.8
	github.com/felixge/fgprof => github.com/felixge/fgprof v0.9.5
	github.com/fsnotify/fsnotify => github.com/fsnotify/fsnotify v1.9.0
	github.com/go-redis/redis/v8 => github.com/go-redis/redis/v8 v8.11.5
	github.com/go-sql-driver/mysql => github.com/go-sql-driver/mysql v1.9.3
	github.com/goccy/go-json => github.com/goccy/go-json v0.10.6
	github.com/gocql/gocql => github.com/gocql/gocql v1.7.0
	github.com/gocraft/dbr/v2 => github.com/gocraft/dbr/v2 v2.7.7
	github.com/google/go-cmp => github.com/google/go-cmp v0.7.0
	github.com/google/uuid => github.com/google/uuid v1.6.0
	github.com/gorilla/mux => github.com/gorilla/mux v1.8.1
	github.com/grafana/grafana-foundation-sdk/go => github.com/grafana/grafana-foundation-sdk/go v0.0.0-20260129154346-aba721fdefde
	github.com/grafana/promql-builder/go => github.com/grafana/promql-builder/go v0.0.0-20250916111012-8fa9625b89a3
	github.com/grafana/pyroscope-go/godeltaprof => github.com/grafana/pyroscope-go/godeltaprof v0.1.9
	github.com/hashicorp/go-version => github.com/hashicorp/go-version v1.9.0
	github.com/klauspost/compress => github.com/klauspost/compress v1.18.5
	github.com/kpango/fastime => github.com/kpango/fastime v1.1.10
	github.com/kpango/gache/v2 => github.com/kpango/gache/v2 v2.1.9
	github.com/kpango/glg => github.com/kpango/glg v1.6.15
	github.com/kubernetes-csi/external-snapshotter/client/v6 => github.com/kubernetes-csi/external-snapshotter/client/v6 v6.3.0
	github.com/leanovate/gopter => github.com/leanovate/gopter v0.2.11
	github.com/lucasb-eyer/go-colorful => github.com/lucasb-eyer/go-colorful v1.4.0
	github.com/pierrec/lz4/v3 => github.com/pierrec/lz4/v3 v3.3.5
	github.com/planetscale/vtprotobuf => github.com/planetscale/vtprotobuf v0.6.0
	github.com/quasilyte/go-ruleguard => github.com/quasilyte/go-ruleguard v0.4.5
	github.com/quasilyte/go-ruleguard/dsl => github.com/quasilyte/go-ruleguard/dsl v0.3.23
	github.com/quic-go/quic-go => github.com/quic-go/quic-go v0.59.0
	github.com/scylladb/gocqlx => github.com/scylladb/gocqlx v1.5.0
	github.com/stretchr/testify => github.com/stretchr/testify v1.11.1
	github.com/zeebo/xxh3 => github.com/zeebo/xxh3 v1.1.0
	go.etcd.io/bbolt => go.etcd.io/bbolt v1.4.3
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc => go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.67.0
	go.opentelemetry.io/otel => go.opentelemetry.io/otel v1.42.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc => go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.42.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace => go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.42.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc => go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.42.0
	go.opentelemetry.io/otel/metric => go.opentelemetry.io/otel/metric v1.42.0
	go.opentelemetry.io/otel/sdk => go.opentelemetry.io/otel/sdk v1.42.0
	go.opentelemetry.io/otel/sdk/metric => go.opentelemetry.io/otel/sdk/metric v1.42.0
	go.opentelemetry.io/otel/trace => go.opentelemetry.io/otel/trace v1.42.0
	go.uber.org/automaxprocs => go.uber.org/automaxprocs v1.6.0
	go.uber.org/goleak => go.uber.org/goleak v1.3.0
	go.uber.org/ratelimit => go.uber.org/ratelimit v0.3.1
	go.uber.org/zap => go.uber.org/zap v1.27.1
	gocloud.dev => gocloud.dev v0.45.0
	golang.org/x/net => golang.org/x/net v0.52.0
	golang.org/x/oauth2 => golang.org/x/oauth2 v0.36.0
	golang.org/x/sync => golang.org/x/sync v0.20.0
	golang.org/x/sys => golang.org/x/sys v0.42.0
	golang.org/x/text => golang.org/x/text v0.35.0
	golang.org/x/time => golang.org/x/time v0.15.0
	golang.org/x/tools => golang.org/x/tools v0.43.0
	gonum.org/v1/hdf5 => gonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc6946
	gonum.org/v1/plot => gonum.org/v1/plot v0.16.0
	google.golang.org/genproto/googleapis/api => google.golang.org/genproto/googleapis/api v0.0.0-20260330182312-d5a96adf58d8
	google.golang.org/genproto/googleapis/rpc => google.golang.org/genproto/googleapis/rpc v0.0.0-20260330182312-d5a96adf58d8
	google.golang.org/grpc => google.golang.org/grpc v1.79.3
	google.golang.org/protobuf => google.golang.org/protobuf v1.36.11
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.1
	k8s.io/api => k8s.io/api v0.35.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.35.3
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.35.3
	k8s.io/metrics => k8s.io/metrics v0.35.3
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.23.3
	sigs.k8s.io/yaml => sigs.k8s.io/yaml v1.6.0
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v0.0.0-00010101000000-000000000000
	cloud.google.com/go/storage v1.57.2
	code.cloudfoundry.org/bytefmt v0.0.0-20190710193110-1eb035ffe2b6
	github.com/akrylysov/pogreb v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go v1.55.8
	github.com/felixge/fgprof v0.0.0-00010101000000-000000000000
	github.com/fsnotify/fsnotify v1.9.0
	github.com/go-redis/redis/v8 v8.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.9.3
	github.com/goccy/go-json v0.10.6
	github.com/gocql/gocql v0.0.0-20200131111108-92af2e088537
	github.com/gocraft/dbr/v2 v2.0.0-00010101000000-000000000000
	github.com/google/go-cmp v0.7.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.7.3
	github.com/grafana/grafana-foundation-sdk/go v0.0.0-00010101000000-000000000000
	github.com/grafana/promql-builder/go v0.0.0-00010101000000-000000000000
	github.com/grafana/pyroscope-go/godeltaprof v0.0.0-00010101000000-000000000000
	github.com/hashicorp/go-version v1.2.0
	github.com/klauspost/compress v1.18.0
	github.com/kpango/fastime v1.1.10
	github.com/kpango/gache/v2 v2.0.0-00010101000000-000000000000
	github.com/kpango/glg v1.6.15
	github.com/kubernetes-csi/external-snapshotter/client/v6 v6.0.0-00010101000000-000000000000
	github.com/leanovate/gopter v0.0.0-00010101000000-000000000000
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/pierrec/lz4/v3 v3.0.0-00010101000000-000000000000
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10
	github.com/quasilyte/go-ruleguard v0.0.0-00010101000000-000000000000
	github.com/quasilyte/go-ruleguard/dsl v0.3.22
	github.com/quic-go/quic-go v0.0.0-00010101000000-000000000000
	github.com/scylladb/gocqlx v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
	github.com/zeebo/xxh3 v1.1.0
	go.etcd.io/bbolt v1.4.3
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.63.0
	go.opentelemetry.io/otel v1.42.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.38.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.42.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.38.0
	go.opentelemetry.io/otel/metric v1.42.0
	go.opentelemetry.io/otel/sdk v1.42.0
	go.opentelemetry.io/otel/sdk/metric v1.42.0
	go.opentelemetry.io/otel/trace v1.42.0
	go.uber.org/automaxprocs v1.6.0
	go.uber.org/goleak v1.3.0
	go.uber.org/ratelimit v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.27.1
	gocloud.dev v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.52.0
	golang.org/x/oauth2 v0.36.0
	golang.org/x/sync v0.20.0
	golang.org/x/sys v0.42.0
	golang.org/x/text v0.35.0
	golang.org/x/time v0.15.0
	golang.org/x/tools v0.43.0
	gonum.org/v1/hdf5 v0.0.0-00010101000000-000000000000
	gonum.org/v1/plot v0.15.2
	google.golang.org/genproto/googleapis/api v0.0.0-20260319201613-d00831a3d3e7
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260319201613-d00831a3d3e7
	google.golang.org/grpc v1.79.3
	google.golang.org/protobuf v1.36.11
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/api v0.35.3
	k8s.io/apimachinery v0.35.3
	k8s.io/cli-runtime v0.0.0-00010101000000-000000000000
	k8s.io/client-go v0.35.3
	k8s.io/metrics v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.0.0-00010101000000-000000000000
	sigs.k8s.io/yaml v1.6.0
)

require (
	cel.dev/expr v0.25.1 // indirect
	cloud.google.com/go v0.123.0 // indirect
	cloud.google.com/go/auth v0.19.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	cloud.google.com/go/iam v1.6.0 // indirect
	cloud.google.com/go/monitoring v1.24.3 // indirect
	codeberg.org/go-fonts/liberation v0.5.0 // indirect
	codeberg.org/go-latex/latex v0.2.0 // indirect
	codeberg.org/go-pdf/fpdf v0.11.1 // indirect
	filippo.io/edwards25519 v1.2.0 // indirect
	git.sr.ht/~sbinet/gg v0.7.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.31.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.55.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.55.0 // indirect
	github.com/ajstarks/svgo v0.0.0-20211024235047-1546f124cd8b // indirect
	github.com/benbjohnson/clock v1.3.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/campoy/embedmd v1.0.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cncf/xds/go v0.0.0-20260202195803-dba9d589def2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/emicklei/go-restful/v3 v3.13.0 // indirect
	github.com/envoyproxy/go-control-plane/envoy v1.37.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.3.3 // indirect
	github.com/evanphx/json-patch/v5 v5.9.11 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fxamacker/cbor/v2 v2.9.1 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/go-jose/go-jose/v4 v4.1.3 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.22.5 // indirect
	github.com/go-openapi/jsonreference v0.21.5 // indirect
	github.com/go-openapi/swag v0.25.5 // indirect
	github.com/go-openapi/swag/cmdutils v0.25.5 // indirect
	github.com/go-openapi/swag/conv v0.25.5 // indirect
	github.com/go-openapi/swag/fileutils v0.25.5 // indirect
	github.com/go-openapi/swag/jsonname v0.25.5 // indirect
	github.com/go-openapi/swag/jsonutils v0.25.5 // indirect
	github.com/go-openapi/swag/loading v0.25.5 // indirect
	github.com/go-openapi/swag/mangling v0.25.5 // indirect
	github.com/go-openapi/swag/netutils v0.25.5 // indirect
	github.com/go-openapi/swag/stringutils v0.25.5 // indirect
	github.com/go-openapi/swag/typeutils v0.25.5 // indirect
	github.com/go-openapi/swag/yamlutils v0.25.5 // indirect
	github.com/go-toolsmith/astcopy v1.0.2 // indirect
	github.com/go-toolsmith/astequal v1.0.3 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/google/gnostic-models v0.7.1 // indirect
	github.com/google/pprof v0.0.0-20260302011040-a15ffb7f9dcc // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/google/wire v0.7.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.14 // indirect
	github.com/googleapis/gax-go/v2 v2.20.0 // indirect
	github.com/gorilla/websocket v1.5.4-0.20250319132907-e064f32e3674 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.28.0 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/mailru/easyjson v0.9.2 // indirect
	github.com/moby/spdystream v0.5.0 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_golang v1.23.2 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.67.5 // indirect
	github.com/prometheus/procfs v0.20.1 // indirect
	github.com/quasilyte/gogrep v0.5.0 // indirect
	github.com/quasilyte/stdinfo v0.0.0-20220114132959-f7386bf02567 // indirect
	github.com/scylladb/go-reflectx v1.0.1 // indirect
	github.com/spf13/cobra v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/spiffe/go-spiffe/v2 v2.6.0 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xlab/treeprint v1.2.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/detectors/gcp v1.42.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.67.0 // indirect
	go.opentelemetry.io/proto/otlp v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v2 v2.4.4 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/exp/typeparams v0.0.0-20240213143201-ec583247a57a // indirect
	golang.org/x/image v0.38.0 // indirect
	golang.org/x/mod v0.34.0 // indirect
	golang.org/x/term v0.41.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	gomodules.xyz/jsonpatch/v2 v2.5.0 // indirect
	google.golang.org/api v0.273.0 // indirect
	google.golang.org/genproto v0.0.0-20260330182312-d5a96adf58d8 // indirect
	gopkg.in/evanphx/json-patch.v4 v4.13.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/apiextensions-apiserver v0.35.3 // indirect
	k8s.io/klog/v2 v2.140.0 // indirect
	k8s.io/kube-openapi v0.0.0-20260330154417-16be699c7b31 // indirect
	k8s.io/utils v0.0.0-20260319190234-28399d86e0b5 // indirect
	sigs.k8s.io/json v0.0.0-20250730193827-2d320260d730 // indirect
	sigs.k8s.io/kustomize/api v0.20.1 // indirect
	sigs.k8s.io/kustomize/kyaml v0.20.1 // indirect
	sigs.k8s.io/randfill v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v6 v6.3.2 // indirect
)
