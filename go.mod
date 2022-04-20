module github.com/vdaas/vald

go 1.18

replace (
	cloud.google.com/go => cloud.google.com/go v0.100.2
	cloud.google.com/go/iam => cloud.google.com/go/iam v0.3.0
	cloud.google.com/go/monitoring => cloud.google.com/go/monitoring v1.4.0
	cloud.google.com/go/profiler => cloud.google.com/go/profiler v0.2.0
	cloud.google.com/go/storage => cloud.google.com/go/storage v1.22.0
	cloud.google.com/go/trace => cloud.google.com/go/trace v1.2.0
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.1-0.20220331191732-bad3b7cbb013+incompatible
	github.com/Azure/go-autorest/autorest => github.com/Azure/go-autorest/autorest v0.11.25
	github.com/Azure/go-autorest/autorest/adal => github.com/Azure/go-autorest/autorest/adal v0.9.19-0.20220331191732-bad3b7cbb013
	github.com/aws/aws-sdk-go => github.com/aws/aws-sdk-go v1.43.35
	github.com/chzyer/logex => github.com/chzyer/logex v1.2.0
	github.com/coreos/etcd => go.etcd.io/etcd v3.3.27+incompatible
	github.com/docker/docker => github.com/moby/moby v20.10.14+incompatible
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v0.6.7
	github.com/fsnotify/fsnotify => github.com/fsnotify/fsnotify v1.5.1
	github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.7.7
	github.com/go-logr/logr => github.com/go-logr/logr v1.2.3
	github.com/go-sql-driver/mysql => github.com/go-sql-driver/mysql v1.6.0
	github.com/goccy/go-json => github.com/goccy/go-json v0.9.6
	github.com/gocql/gocql => github.com/gocql/gocql v1.0.0
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
	github.com/golang/groupcache => github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da
	github.com/golang/protobuf => github.com/golang/protobuf v1.5.2
	github.com/golang/snappy => github.com/golang/snappy v0.0.4
	github.com/google/btree => github.com/google/btree v1.0.1
	github.com/google/go-cmp => github.com/google/go-cmp v0.5.7
	github.com/google/pprof => github.com/google/pprof v0.0.0-20220401020641-b5a4dc8f4f2a
	github.com/google/uuid => github.com/google/uuid v1.3.0
	github.com/gorilla/mux => github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.5.0
	github.com/hailocab/go-hostpool => github.com/kpango/go-hostpool v0.0.0-20210303030322-aab80263dcd0
	github.com/hashicorp/go-version => github.com/hashicorp/go-version v1.4.0
	github.com/jackc/chunkreader => github.com/jackc/chunkreader v1.0.0
	github.com/jackc/pgconn => github.com/jackc/pgconn v1.11.0
	github.com/jackc/pgmock => github.com/jackc/pgmock v0.0.0-20210724152146-4ad1a8207f65
	github.com/jackc/pgproto3/v2 => github.com/jackc/pgproto3/v2 v2.2.0
	github.com/jackc/pgservicefile => github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b
	github.com/jackc/pgtype => github.com/jackc/pgtype v1.10.0
	github.com/jackc/puddle => github.com/jackc/puddle v1.2.1
	github.com/json-iterator/go => github.com/json-iterator/go v1.1.12
	github.com/klauspost/compress => github.com/klauspost/compress v1.15.2-0.20220407105542-c06ba5f93c87
	github.com/kpango/glg => github.com/kpango/glg v1.6.10
	github.com/onsi/ginkgo => github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega => github.com/onsi/gomega v1.19.0
	github.com/opentracing/opentracing-go => github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/sftp => github.com/pkg/sftp v1.13.4
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v1.12.1
	github.com/prometheus/client_model => github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common => github.com/prometheus/common v0.33.0
	github.com/prometheus/procfs => github.com/prometheus/procfs v0.7.3
	github.com/rs/xid => github.com/rs/xid v1.4.0
	github.com/sirupsen/logrus => github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra => github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag => github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify => github.com/stretchr/testify v1.7.1
	github.com/tensorflow/tensorflow => github.com/tensorflow/tensorflow v2.1.2+incompatible
	github.com/zeebo/xxh3 => github.com/zeebo/xxh3 v1.0.2
	go.etcd.io/bbolt => go.etcd.io/bbolt v1.3.6
	go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3 v3.5.2
	go.etcd.io/etcd/server/v3 => go.etcd.io/etcd/server/v3 v3.5.2
	go.opencensus.io => go.opencensus.io v0.23.0
	go.opentelemetry.io/otel => go.opentelemetry.io/otel v1.6.3
	go.opentelemetry.io/otel/metric => go.opentelemetry.io/otel/metric v0.28.0
	go.opentelemetry.io/otel/sdk => go.opentelemetry.io/otel/sdk v1.6.3
	go.opentelemetry.io/otel/trace => go.opentelemetry.io/otel/trace v1.6.3
	go.uber.org/atomic => go.uber.org/atomic v1.9.0
	go.uber.org/goleak => go.uber.org/goleak v1.1.12
	go.uber.org/multierr => go.uber.org/multierr v1.8.0
	go.uber.org/zap => go.uber.org/zap v1.21.0
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20220331220935-ae2d96664a29
	golang.org/x/exp => golang.org/x/exp v0.0.0-20220407100705-7b9b53b0aca4
	golang.org/x/image => golang.org/x/image v0.0.0-20220321031419-a8550c1d254a
	golang.org/x/lint => golang.org/x/lint v0.0.0-20210508222113-6edffad5e616
	golang.org/x/mod => golang.org/x/mod v0.5.1
	golang.org/x/net => golang.org/x/net v0.0.0-20220407224826-aac1ed45d8e3
	golang.org/x/oauth2 => golang.org/x/oauth2 v0.0.0-20220309155454-6242fa91716a
	golang.org/x/sync => golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys => golang.org/x/sys v0.0.0-20220406163625-3f8b81556e12
	golang.org/x/term => golang.org/x/term v0.0.0-20210927222741-03fcf44c2211
	golang.org/x/text => golang.org/x/text v0.3.7
	golang.org/x/time => golang.org/x/time v0.0.0-20220224211638-0e9765cccd65
	golang.org/x/tools => golang.org/x/tools v0.1.10
	golang.org/x/xerrors => golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	gonum.org/v1/gonum => gonum.org/v1/gonum v0.11.0
	gonum.org/v1/plot => gonum.org/v1/plot v0.11.0
	google.golang.org/api => google.golang.org/api v0.74.0
	google.golang.org/appengine => google.golang.org/appengine v1.6.7
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20220407144326-9054f6ed7bac
	google.golang.org/grpc => google.golang.org/grpc v1.45.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc => google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.2.0
	google.golang.org/protobuf => google.golang.org/protobuf v1.28.0
	gopkg.in/check.v1 => gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api => k8s.io/api v0.23.5
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.23.5
	k8s.io/apimachinery => k8s.io/apimachinery v0.23.5
	k8s.io/apiserver => k8s.io/apiserver v0.23.5
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.23.5
	k8s.io/client-go => k8s.io/client-go v0.23.5
	k8s.io/code-generator => k8s.io/code-generator v0.23.5
	k8s.io/component-base => k8s.io/component-base v0.23.5
	k8s.io/gengo => k8s.io/gengo v0.0.0-20220307231824-4627b89bbf1b
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.60.1
	k8s.io/metrics => k8s.io/metrics v0.23.5
	k8s.io/utils => k8s.io/utils v0.0.0-20220210201930-3a6ce19ff2f9
	sigs.k8s.io/apiserver-network-proxy/konnectivity-client => sigs.k8s.io/apiserver-network-proxy/konnectivity-client v0.0.30
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.11.2
	sigs.k8s.io/kustomize/api => sigs.k8s.io/kustomize/api v0.11.4
	sigs.k8s.io/kustomize/kyaml => sigs.k8s.io/kustomize/kyaml v0.13.6
	sigs.k8s.io/structured-merge-diff/v4 => sigs.k8s.io/structured-merge-diff/v4 v4.2.1
	sigs.k8s.io/yaml => sigs.k8s.io/yaml v1.3.0
)

require (
	cloud.google.com/go/profiler v0.0.0-00010101000000-000000000000
	cloud.google.com/go/storage v1.21.0
	code.cloudfoundry.org/bytefmt v0.0.0-20211005130812-5bb3c17173e5
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/prometheus v0.4.1
	contrib.go.opencensus.io/exporter/stackdriver v0.13.12
	github.com/aws/aws-sdk-go v1.43.31
	github.com/envoyproxy/protoc-gen-validate v0.1.0
	github.com/fsnotify/fsnotify v1.5.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-sql-driver/mysql v1.6.0
	github.com/goccy/go-json v0.9.4
	github.com/gocql/gocql v0.0.0-20200131111108-92af2e088537
	github.com/gocraft/dbr/v2 v2.7.3
	github.com/google/go-cmp v0.5.7
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-version v0.0.0-00010101000000-000000000000
	github.com/klauspost/compress v1.15.1
	github.com/kpango/fastime v1.1.4
	github.com/kpango/fuid v0.0.0-20220209050620-e5987ba1ea5e
	github.com/kpango/gache v1.2.7
	github.com/kpango/glg v1.6.10
	github.com/leanovate/gopter v0.2.9
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/pierrec/lz4/v3 v3.3.4
	github.com/quasilyte/go-ruleguard v0.3.15
	github.com/quasilyte/go-ruleguard/dsl v0.3.19
	github.com/scylladb/gocqlx v1.5.0
	github.com/tensorflow/tensorflow v0.0.0-00010101000000-000000000000
	github.com/zeebo/xxh3 v1.0.1
	go.opencensus.io v0.23.0
	go.uber.org/automaxprocs v1.5.1
	go.uber.org/goleak v1.1.12
	go.uber.org/zap v1.21.0
	gocloud.dev v0.25.0
	golang.org/x/net v0.0.0-20220401154927-543a649e0bdd
	golang.org/x/oauth2 v0.0.0-20220309155454-6242fa91716a
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20220330033206-e17cdc41300f
	golang.org/x/tools v0.1.6-0.20210820212750-d4cc65f0b2ff
	gonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc6946
	gonum.org/v1/plot v0.0.0-00010101000000-000000000000
	google.golang.org/api v0.74.0
	google.golang.org/genproto v0.0.0-20220405205423-9d709892a2bf
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v2 v2.4.0
	inet.af/netaddr v0.0.0-20211027220019-c74959edd3b6
	k8s.io/api v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/cli-runtime v0.0.0-00010101000000-000000000000
	k8s.io/client-go v0.23.5
	k8s.io/metrics v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.0.0-00010101000000-000000000000
)

require (
	cloud.google.com/go v0.100.2 // indirect
	cloud.google.com/go/compute v1.5.0 // indirect
	cloud.google.com/go/iam v0.3.0 // indirect
	cloud.google.com/go/monitoring v1.4.0 // indirect
	cloud.google.com/go/trace v1.2.0 // indirect
	git.sr.ht/~sbinet/gg v0.3.1 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/ajstarks/svgo v0.0.0-20211024235047-1546f124cd8b // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/census-instrumentation/opencensus-proto v0.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/go-errors/errors v1.0.1 // indirect
	github.com/go-fonts/liberation v0.2.0 // indirect
	github.com/go-kit/log v0.2.0 // indirect
	github.com/go-latex/latex v0.0.0-20210823091927-c0d11ff05a81 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/go-pdf/fpdf v0.6.0 // indirect
	github.com/go-toolsmith/astcopy v1.0.0 // indirect
	github.com/go-toolsmith/astequal v1.0.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/pprof v0.0.0-20220113144219-d25a53d42d00 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/google/wire v0.5.0 // indirect
	github.com/googleapis/gax-go/v2 v2.2.0 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/googleapis/go-type-adapters v1.0.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20180305231024-9cad4c3443a7 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.12.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/prometheus/prometheus v2.5.0+incompatible // indirect
	github.com/prometheus/statsd_exporter v0.21.0 // indirect
	github.com/quasilyte/gogrep v0.0.0-20220103110004-ffaa07af02e3 // indirect
	github.com/scylladb/go-reflectx v1.0.1 // indirect
	github.com/spf13/cobra v1.2.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	github.com/uber/jaeger-client-go v2.25.0+incompatible // indirect
	github.com/xlab/treeprint v0.0.0-20181112141820-a009c3971eca // indirect
	go.starlark.net v0.0.0-20200306205701-8dd3e2ee1dd5 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go4.org/intern v0.0.0-20211027215823-ae77deb06f29 // indirect
	go4.org/unsafe/assume-no-moving-gc v0.0.0-20211027215541-db492cf91b37 // indirect
	golang.org/x/image v0.0.0-20220302094943-723b81ca9867 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220106191415-9b9b3d81d5e3 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20220224211638-0e9765cccd65 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	k8s.io/apiextensions-apiserver v0.23.5 // indirect
	k8s.io/component-base v0.23.5 // indirect
	k8s.io/klog/v2 v2.30.0 // indirect
	k8s.io/kube-openapi v0.0.0-20211115234752-e816edb12b65 // indirect
	k8s.io/utils v0.0.0-20211116205334-6203023598ed // indirect
	sigs.k8s.io/json v0.0.0-20211020170558-c049b76a60c6 // indirect
	sigs.k8s.io/kustomize/api v0.10.1 // indirect
	sigs.k8s.io/kustomize/kyaml v0.13.6 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
