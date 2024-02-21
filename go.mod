module github.com/vdaas/vald

go 1.22.0

replace (
	cloud.google.com/go => cloud.google.com/go v0.112.0
	cloud.google.com/go/bigquery => cloud.google.com/go/bigquery v1.59.1
	cloud.google.com/go/compute => cloud.google.com/go/compute v1.24.0
	cloud.google.com/go/datastore => cloud.google.com/go/datastore v1.15.0
	cloud.google.com/go/firestore => cloud.google.com/go/firestore v1.14.0
	cloud.google.com/go/iam => cloud.google.com/go/iam v1.1.6
	cloud.google.com/go/kms => cloud.google.com/go/kms v1.15.7
	cloud.google.com/go/monitoring => cloud.google.com/go/monitoring v1.18.0
	cloud.google.com/go/pubsub => cloud.google.com/go/pubsub v1.36.1
	cloud.google.com/go/secretmanager => cloud.google.com/go/secretmanager v1.11.5
	cloud.google.com/go/storage => cloud.google.com/go/storage v1.38.0
	cloud.google.com/go/trace => cloud.google.com/go/trace v1.10.5
	code.cloudfoundry.org/bytefmt => code.cloudfoundry.org/bytefmt v0.0.0-20231017140541-3b893ed0421b
	contrib.go.opencensus.io/exporter/aws => contrib.go.opencensus.io/exporter/aws v0.0.0-20230502192102-15967c811cec
	contrib.go.opencensus.io/exporter/prometheus => contrib.go.opencensus.io/exporter/prometheus v0.4.2
	contrib.go.opencensus.io/integrations/ocsql => contrib.go.opencensus.io/integrations/ocsql v0.1.7
	git.sr.ht/~sbinet/gg => git.sr.ht/~sbinet/gg v0.5.0
	github.com/Azure/azure-amqp-common-go/v3 => github.com/Azure/azure-amqp-common-go/v3 v3.2.3
	github.com/Azure/azure-sdk-for-go => github.com/Azure/azure-sdk-for-go v68.0.0+incompatible
	github.com/Azure/azure-sdk-for-go/sdk/azcore => github.com/Azure/azure-sdk-for-go/sdk/azcore v1.9.2
	github.com/Azure/azure-sdk-for-go/sdk/azidentity => github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.5.1
	github.com/Azure/azure-sdk-for-go/sdk/internal => github.com/Azure/azure-sdk-for-go/sdk/internal v1.5.2
	github.com/Azure/go-amqp => github.com/Azure/go-amqp v1.0.4
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.1-0.20240104002855-4c0e21ca2bbb+incompatible
	github.com/Azure/go-autorest/autorest => github.com/Azure/go-autorest/autorest v0.11.30-0.20240104002855-4c0e21ca2bbb
	github.com/Azure/go-autorest/autorest/adal => github.com/Azure/go-autorest/autorest/adal v0.9.23
	github.com/Azure/go-autorest/autorest/date => github.com/Azure/go-autorest/autorest/date v0.3.1-0.20240104002855-4c0e21ca2bbb
	github.com/Azure/go-autorest/autorest/mocks => github.com/Azure/go-autorest/autorest/mocks v0.4.3-0.20240104002855-4c0e21ca2bbb
	github.com/Azure/go-autorest/autorest/to => github.com/Azure/go-autorest/autorest/to v0.4.1-0.20240104002855-4c0e21ca2bbb
	github.com/Azure/go-autorest/logger => github.com/Azure/go-autorest/logger v0.2.2-0.20240104002855-4c0e21ca2bbb
	github.com/Azure/go-autorest/tracing => github.com/Azure/go-autorest/tracing v0.6.1-0.20240104002855-4c0e21ca2bbb
	github.com/BurntSushi/toml => github.com/BurntSushi/toml v1.3.2
	github.com/DATA-DOG/go-sqlmock => github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/GoogleCloudPlatform/cloudsql-proxy => github.com/GoogleCloudPlatform/cloudsql-proxy v1.33.16
	github.com/Masterminds/semver/v3 => github.com/Masterminds/semver/v3 v3.2.1
	github.com/ajstarks/deck => github.com/ajstarks/deck v0.0.0-20240215023147-77e8381c6641
	github.com/ajstarks/deck/generate => github.com/ajstarks/deck/generate v0.0.0-20240215023147-77e8381c6641
	github.com/ajstarks/svgo => github.com/ajstarks/svgo v0.0.0-20211024235047-1546f124cd8b
	github.com/akrylysov/pogreb => github.com/akrylysov/pogreb v0.10.2
	github.com/antihax/optional => github.com/antihax/optional v1.0.0
	github.com/armon/go-socks5 => github.com/armon/go-socks5 v0.0.0-20160902184237-e75332964ef5
	github.com/aws/aws-sdk-go => github.com/aws/aws-sdk-go v1.50.21
	github.com/aws/aws-sdk-go-v2 => github.com/aws/aws-sdk-go-v2 v1.25.0
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.0
	github.com/aws/aws-sdk-go-v2/config => github.com/aws/aws-sdk-go-v2/config v1.27.0
	github.com/aws/aws-sdk-go-v2/credentials => github.com/aws/aws-sdk-go-v2/credentials v1.17.0
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds => github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.15.0
	github.com/aws/aws-sdk-go-v2/feature/s3/manager => github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.16.2
	github.com/aws/aws-sdk-go-v2/internal/configsources => github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.0
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.0
	github.com/aws/aws-sdk-go-v2/internal/ini => github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.0
	github.com/aws/aws-sdk-go-v2/service/internal/checksum => github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.0
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.0
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared => github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.0
	github.com/aws/aws-sdk-go-v2/service/kms => github.com/aws/aws-sdk-go-v2/service/kms v1.28.1
	github.com/aws/aws-sdk-go-v2/service/s3 => github.com/aws/aws-sdk-go-v2/service/s3 v1.50.1
	github.com/aws/aws-sdk-go-v2/service/secretsmanager => github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.27.1
	github.com/aws/aws-sdk-go-v2/service/sns => github.com/aws/aws-sdk-go-v2/service/sns v1.28.0
	github.com/aws/aws-sdk-go-v2/service/sqs => github.com/aws/aws-sdk-go-v2/service/sqs v1.30.1
	github.com/aws/aws-sdk-go-v2/service/ssm => github.com/aws/aws-sdk-go-v2/service/ssm v1.47.0
	github.com/aws/aws-sdk-go-v2/service/sso => github.com/aws/aws-sdk-go-v2/service/sso v1.19.0
	github.com/aws/aws-sdk-go-v2/service/sts => github.com/aws/aws-sdk-go-v2/service/sts v1.27.0
	github.com/aws/smithy-go => github.com/aws/smithy-go v1.20.0
	github.com/benbjohnson/clock => github.com/benbjohnson/clock v1.3.5
	github.com/beorn7/perks => github.com/beorn7/perks v1.0.1
	github.com/bmizerany/assert => github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869
	github.com/boombuler/barcode => github.com/boombuler/barcode v1.0.1
	github.com/buger/jsonparser => github.com/buger/jsonparser v1.1.1
	github.com/cenkalti/backoff/v4 => github.com/cenkalti/backoff/v4 v4.2.1
	github.com/census-instrumentation/opencensus-proto => github.com/census-instrumentation/opencensus-proto v0.4.1
	github.com/cespare/xxhash/v2 => github.com/cespare/xxhash/v2 v2.2.0
	github.com/chzyer/logex => github.com/chzyer/logex v1.2.1
	github.com/chzyer/readline => github.com/chzyer/readline v1.5.1
	github.com/chzyer/test => github.com/chzyer/test v1.0.0
	github.com/cncf/udpa/go => github.com/cncf/udpa/go v0.0.0-20220112060539-c52dc94e7fbe
	github.com/cncf/xds/go => github.com/cncf/xds/go v0.0.0-20231128003011-0fa0005c9caa
	github.com/cockroachdb/apd => github.com/cockroachdb/apd v1.1.0
	github.com/coreos/go-systemd/v22 => github.com/coreos/go-systemd/v22 v22.5.0
	github.com/cpuguy83/go-md2man/v2 => github.com/cpuguy83/go-md2man/v2 v2.0.3
	github.com/creack/pty => github.com/creack/pty v1.1.21
	github.com/davecgh/go-spew => github.com/davecgh/go-spew v1.1.1
	github.com/denisenkom/go-mssqldb => github.com/denisenkom/go-mssqldb v0.12.3
	github.com/devigned/tab => github.com/devigned/tab v0.1.1
	github.com/dgryski/go-rendezvous => github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f
	github.com/dnaeon/go-vcr => github.com/dnaeon/go-vcr v1.2.0
	github.com/docopt/docopt-go => github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
	github.com/dustin/go-humanize => github.com/dustin/go-humanize v1.0.1
	github.com/emicklei/go-restful/v3 => github.com/emicklei/go-restful/v3 v3.11.2
	github.com/envoyproxy/go-control-plane => github.com/envoyproxy/go-control-plane v0.12.0
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v1.0.4
	github.com/evanphx/json-patch => github.com/evanphx/json-patch v0.5.2
	github.com/fogleman/gg => github.com/fogleman/gg v1.3.0
	github.com/fortytw2/leaktest => github.com/fortytw2/leaktest v1.3.0
	github.com/frankban/quicktest => github.com/frankban/quicktest v1.14.6
	github.com/fsnotify/fsnotify => github.com/fsnotify/fsnotify v1.7.0
	github.com/gin-contrib/sse => github.com/gin-contrib/sse v0.1.0
	github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.9.1
	github.com/go-errors/errors => github.com/go-errors/errors v1.5.1
	github.com/go-fonts/dejavu => github.com/go-fonts/dejavu v0.3.3
	github.com/go-fonts/latin-modern => github.com/go-fonts/latin-modern v0.3.2
	github.com/go-fonts/liberation => github.com/go-fonts/liberation v0.3.2
	github.com/go-fonts/stix => github.com/go-fonts/stix v0.2.2
	github.com/go-gl/gl => github.com/go-gl/gl v0.0.0-20231021071112-07e5d0ea2e71
	github.com/go-gl/glfw/v3.3/glfw => github.com/go-gl/glfw/v3.3/glfw v0.0.0-20240118000515-a250818d05e3
	github.com/go-kit/log => github.com/go-kit/log v0.2.1
	github.com/go-latex/latex => github.com/go-latex/latex v0.0.0-20231108140139-5c1ce85aa4ea
	github.com/go-logfmt/logfmt => github.com/go-logfmt/logfmt v0.6.0
	github.com/go-logr/logr => github.com/go-logr/logr v1.4.1
	github.com/go-logr/stdr => github.com/go-logr/stdr v1.2.2
	github.com/go-logr/zapr => github.com/go-logr/zapr v1.3.0
	github.com/go-openapi/jsonpointer => github.com/go-openapi/jsonpointer v0.20.2
	github.com/go-openapi/jsonreference => github.com/go-openapi/jsonreference v0.20.4
	github.com/go-openapi/swag => github.com/go-openapi/swag v0.22.9
	github.com/go-pdf/fpdf => github.com/go-pdf/fpdf v1.4.3
	github.com/go-playground/assert/v2 => github.com/go-playground/assert/v2 v2.2.0
	github.com/go-playground/locales => github.com/go-playground/locales v0.14.1
	github.com/go-playground/universal-translator => github.com/go-playground/universal-translator v0.18.1
	github.com/go-playground/validator/v10 => github.com/go-playground/validator/v10 v10.18.0
	github.com/go-redis/redis/v8 => github.com/go-redis/redis/v8 v8.11.5
	github.com/go-sql-driver/mysql => github.com/go-sql-driver/mysql v1.7.1
	github.com/go-task/slim-sprig => github.com/go-task/slim-sprig v2.20.0+incompatible
	github.com/go-toolsmith/astcopy => github.com/go-toolsmith/astcopy v1.1.0
	github.com/go-toolsmith/astequal => github.com/go-toolsmith/astequal v1.2.0
	github.com/go-toolsmith/strparse => github.com/go-toolsmith/strparse v1.1.0
	github.com/gobwas/httphead => github.com/gobwas/httphead v0.1.0
	github.com/gobwas/pool => github.com/gobwas/pool v0.2.1
	github.com/gobwas/ws => github.com/gobwas/ws v1.3.2
	github.com/goccy/go-json => github.com/goccy/go-json v0.10.2
	github.com/gocql/gocql => github.com/gocql/gocql v1.6.0
	github.com/gocraft/dbr/v2 => github.com/gocraft/dbr/v2 v2.7.6
	github.com/godbus/dbus/v5 => github.com/godbus/dbus/v5 v5.1.0
	github.com/gofrs/uuid => github.com/gofrs/uuid v4.4.0+incompatible
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
	github.com/golang-jwt/jwt/v4 => github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/golang-sql/civil => github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9
	github.com/golang-sql/sqlexp => github.com/golang-sql/sqlexp v0.1.0
	github.com/golang/freetype => github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/golang/glog => github.com/golang/glog v1.2.0
	github.com/golang/groupcache => github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da
	github.com/golang/mock => github.com/golang/mock v1.6.0
	github.com/golang/protobuf => github.com/golang/protobuf v1.5.3
	github.com/golang/snappy => github.com/golang/snappy v0.0.4
	github.com/google/btree => github.com/google/btree v1.1.2
	github.com/google/gnostic => github.com/google/gnostic v0.7.0
	github.com/google/go-cmp => github.com/google/go-cmp v0.6.0
	github.com/google/go-replayers/grpcreplay => github.com/google/go-replayers/grpcreplay v1.1.0
	github.com/google/go-replayers/httpreplay => github.com/google/go-replayers/httpreplay v1.2.0
	github.com/google/gofuzz => github.com/google/gofuzz v1.2.0
	github.com/google/martian => github.com/google/martian v2.1.0+incompatible
	github.com/google/martian/v3 => github.com/google/martian/v3 v3.3.2
	github.com/google/pprof => github.com/google/pprof v0.0.0-20240207164012-fb44976bdcd5
	github.com/google/shlex => github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/google/subcommands => github.com/google/subcommands v1.2.0
	github.com/google/uuid => github.com/google/uuid v1.6.0
	github.com/google/wire => github.com/google/wire v0.6.0
	github.com/googleapis/gax-go/v2 => github.com/googleapis/gax-go/v2 v2.12.1
	github.com/gorilla/mux => github.com/gorilla/mux v1.8.1
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.5.1
	github.com/gregjones/httpcache => github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79
	github.com/grpc-ecosystem/grpc-gateway/v2 => github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1
	github.com/hailocab/go-hostpool => github.com/kpango/go-hostpool v0.0.0-20210303030322-aab80263dcd0
	github.com/hanwen/go-fuse/v2 => github.com/hanwen/go-fuse/v2 v2.5.0
	github.com/hashicorp/go-uuid => github.com/hashicorp/go-uuid v1.0.3
	github.com/hashicorp/go-version => github.com/hashicorp/go-version v1.6.0
	github.com/iancoleman/strcase => github.com/iancoleman/strcase v0.3.0
	github.com/ianlancetaylor/demangle => github.com/ianlancetaylor/demangle v0.0.0-20240205174729-1f824a1a9b87
	github.com/inconshreveable/mousetrap => github.com/inconshreveable/mousetrap v1.1.0
	github.com/jackc/chunkreader/v2 => github.com/jackc/chunkreader/v2 v2.0.1
	github.com/jackc/pgconn => github.com/jackc/pgconn v1.14.1
	github.com/jackc/pgio => github.com/jackc/pgio v1.0.0
	github.com/jackc/pgmock => github.com/jackc/pgmock v0.0.0-20210724152146-4ad1a8207f65
	github.com/jackc/pgpassfile => github.com/jackc/pgpassfile v1.0.0
	github.com/jackc/pgproto3/v2 => github.com/jackc/pgproto3/v2 v2.3.2
	github.com/jackc/pgservicefile => github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9
	github.com/jackc/pgtype => github.com/jackc/pgtype v1.14.2
	github.com/jackc/pgx/v4 => github.com/jackc/pgx/v4 v4.18.1
	github.com/jackc/puddle => github.com/jackc/puddle v1.3.0
	github.com/jessevdk/go-flags => github.com/jessevdk/go-flags v1.5.0
	github.com/jmespath/go-jmespath => github.com/jmespath/go-jmespath v0.4.0
	github.com/jmespath/go-jmespath/internal/testify => github.com/jmespath/go-jmespath/internal/testify v1.5.1
	github.com/jmoiron/sqlx => github.com/jmoiron/sqlx v1.3.5
	github.com/joho/godotenv => github.com/joho/godotenv v1.5.1
	github.com/josharian/intern => github.com/josharian/intern v1.0.0
	github.com/json-iterator/go => github.com/json-iterator/go v1.1.12
	github.com/jstemmer/go-junit-report => github.com/jstemmer/go-junit-report v1.0.0
	github.com/kisielk/errcheck => github.com/kisielk/errcheck v1.7.0
	github.com/kisielk/gotool => github.com/kisielk/gotool v1.0.0
	github.com/klauspost/compress => github.com/klauspost/compress v1.17.7-0.20240212204355-5895eb4259d0
	github.com/klauspost/cpuid/v2 => github.com/klauspost/cpuid/v2 v2.2.6
	github.com/kpango/fastime => github.com/kpango/fastime v1.1.9
	github.com/kpango/fuid => github.com/kpango/fuid v0.0.0-20221203053508-503b5ad89aa1
	github.com/kpango/gache/v2 => github.com/kpango/gache/v2 v2.0.9
	github.com/kpango/glg => github.com/kpango/glg v1.6.15
	github.com/kr/fs => github.com/kr/fs v0.1.0
	github.com/kr/pretty => github.com/kr/pretty v0.3.1
	github.com/kr/text => github.com/kr/text v0.2.0
	github.com/kubernetes-csi/external-snapshotter/client/v6 => github.com/kubernetes-csi/external-snapshotter/client/v6 v6.3.0
	github.com/kylelemons/godebug => github.com/kylelemons/godebug v1.1.0
	github.com/leanovate/gopter => github.com/leanovate/gopter v0.2.9
	github.com/leodido/go-urn => github.com/leodido/go-urn v1.4.0
	github.com/lib/pq => github.com/lib/pq v1.10.9
	github.com/liggitt/tabwriter => github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de
	github.com/lucasb-eyer/go-colorful => github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/mailru/easyjson => github.com/mailru/easyjson v0.7.7
	github.com/mattn/go-colorable => github.com/mattn/go-colorable v0.1.13
	github.com/mattn/go-isatty => github.com/mattn/go-isatty v0.0.20
	github.com/mattn/go-sqlite3 => github.com/mattn/go-sqlite3 v1.14.22
	github.com/matttproud/golang_protobuf_extensions => github.com/matttproud/golang_protobuf_extensions v1.0.4
	github.com/mitchellh/colorstring => github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db
	github.com/moby/spdystream => github.com/moby/spdystream v0.2.0
	github.com/moby/sys/mountinfo => github.com/moby/sys/mountinfo v0.7.1
	github.com/modern-go/concurrent => github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
	github.com/modern-go/reflect2 => github.com/modern-go/reflect2 v1.0.2
	github.com/modocache/gover => github.com/modocache/gover v0.0.0-20171022184752-b58185e213c5
	github.com/monochromegane/go-gitignore => github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00
	github.com/montanaflynn/stats => github.com/montanaflynn/stats v0.7.1
	github.com/munnerz/goautoneg => github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822
	github.com/niemeyer/pretty => github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e
	github.com/nxadm/tail => github.com/nxadm/tail v1.4.11
	github.com/onsi/ginkgo => github.com/onsi/ginkgo v1.16.5
	github.com/onsi/ginkgo/v2 => github.com/onsi/ginkgo/v2 v2.15.0
	github.com/onsi/gomega => github.com/onsi/gomega v1.31.1
	github.com/peterbourgon/diskv => github.com/peterbourgon/diskv v2.0.1+incompatible
	github.com/phpdave11/gofpdf => github.com/phpdave11/gofpdf v1.4.2
	github.com/phpdave11/gofpdi => github.com/phpdave11/gofpdi v1.0.13
	github.com/pierrec/cmdflag => github.com/pierrec/cmdflag v0.0.2
	github.com/pierrec/lz4/v3 => github.com/pierrec/lz4/v3 v3.3.5
	github.com/pkg/browser => github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c
	github.com/pkg/errors => github.com/pkg/errors v0.9.1
	github.com/pkg/sftp => github.com/pkg/sftp v1.13.6
	github.com/pmezard/go-difflib => github.com/pmezard/go-difflib v1.0.0
	github.com/prashantv/gostub => github.com/prashantv/gostub v1.1.0
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v1.18.0
	github.com/prometheus/client_model => github.com/prometheus/client_model v0.6.0
	github.com/prometheus/common => github.com/prometheus/common v0.47.0
	github.com/prometheus/procfs => github.com/prometheus/procfs v0.12.0
	github.com/prometheus/prometheus => github.com/prometheus/prometheus v1.99.0
	github.com/quasilyte/go-ruleguard => github.com/quasilyte/go-ruleguard v0.4.0
	github.com/quasilyte/go-ruleguard/dsl => github.com/quasilyte/go-ruleguard/dsl v0.3.22
	github.com/quasilyte/gogrep => github.com/quasilyte/gogrep v0.5.0
	github.com/quasilyte/stdinfo => github.com/quasilyte/stdinfo v0.0.0-20220114132959-f7386bf02567
	github.com/rogpeppe/fastuuid => github.com/rogpeppe/fastuuid v1.2.0
	github.com/rogpeppe/go-internal => github.com/rogpeppe/go-internal v1.12.0
	github.com/rs/xid => github.com/rs/xid v1.5.0
	github.com/rs/zerolog => github.com/rs/zerolog v1.32.0
	github.com/russross/blackfriday/v2 => github.com/russross/blackfriday/v2 v2.1.0
	github.com/ruudk/golang-pdf417 => github.com/ruudk/golang-pdf417 v0.0.0-20201230142125-a7e3863a1245
	github.com/schollz/progressbar/v2 => github.com/schollz/progressbar/v2 v2.15.0
	github.com/scylladb/go-reflectx => github.com/scylladb/go-reflectx v1.0.1
	github.com/scylladb/gocqlx => github.com/scylladb/gocqlx v1.5.0
	github.com/sergi/go-diff => github.com/sergi/go-diff v1.3.1
	github.com/shopspring/decimal => github.com/shopspring/decimal v1.3.1
	github.com/shurcooL/httpfs => github.com/shurcooL/httpfs v0.0.0-20230704072500-f1e31cf0ba5c
	github.com/shurcooL/vfsgen => github.com/shurcooL/vfsgen v0.0.0-20230704071429-0000e147ea92
	github.com/sirupsen/logrus => github.com/sirupsen/logrus v1.9.3
	github.com/spf13/afero => github.com/spf13/afero v1.11.0
	github.com/spf13/cobra => github.com/spf13/cobra v1.8.0
	github.com/spf13/pflag => github.com/spf13/pflag v1.0.5
	github.com/stoewer/go-strcase => github.com/stoewer/go-strcase v1.3.0
	github.com/stretchr/objx => github.com/stretchr/objx v0.5.1
	github.com/stretchr/testify => github.com/stretchr/testify v1.8.4
	github.com/ugorji/go/codec => github.com/ugorji/go/codec v1.2.12
	github.com/xeipuuv/gojsonpointer => github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb
	github.com/xeipuuv/gojsonreference => github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415
	github.com/xeipuuv/gojsonschema => github.com/xeipuuv/gojsonschema v1.2.0
	github.com/xlab/treeprint => github.com/xlab/treeprint v1.2.0
	github.com/zeebo/assert => github.com/zeebo/assert v1.3.1
	github.com/zeebo/xxh3 => github.com/zeebo/xxh3 v1.0.2
	go.etcd.io/bbolt => go.etcd.io/bbolt v1.3.8
	go.opencensus.io => go.opencensus.io v0.24.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc => go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.48.0
	go.opentelemetry.io/otel => go.opentelemetry.io/otel v1.23.1
	go.opentelemetry.io/otel/exporters/otlp/internal/retry => go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.17.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric => go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.43.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc => go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.23.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace => go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.23.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc => go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.23.1
	go.opentelemetry.io/otel/metric => go.opentelemetry.io/otel/metric v1.23.1
	go.opentelemetry.io/otel/sdk => go.opentelemetry.io/otel/sdk v1.23.1
	go.opentelemetry.io/otel/sdk/metric => go.opentelemetry.io/otel/sdk/metric v1.23.1
	go.opentelemetry.io/otel/trace => go.opentelemetry.io/otel/trace v1.23.1
	go.opentelemetry.io/proto/otlp => go.opentelemetry.io/proto/otlp v1.1.0
	go.starlark.net => go.starlark.net v0.0.0-20240123142251-f86470692795
	go.uber.org/atomic => go.uber.org/atomic v1.11.0
	go.uber.org/automaxprocs => go.uber.org/automaxprocs v1.5.3
	go.uber.org/goleak => go.uber.org/goleak v1.3.0
	go.uber.org/multierr => go.uber.org/multierr v1.11.0
	go.uber.org/zap => go.uber.org/zap v1.26.0
	gocloud.dev => gocloud.dev v0.36.0
	golang.org/x/crypto => golang.org/x/crypto v0.19.0
	golang.org/x/exp => golang.org/x/exp v0.0.0-20240213143201-ec583247a57a
	golang.org/x/exp/typeparams => golang.org/x/exp/typeparams v0.0.0-20240213143201-ec583247a57a
	golang.org/x/image => golang.org/x/image v0.15.0
	golang.org/x/lint => golang.org/x/lint v0.0.0-20210508222113-6edffad5e616
	golang.org/x/mobile => golang.org/x/mobile v0.0.0-20240213143359-d1f7d3436075
	golang.org/x/mod => golang.org/x/mod v0.15.0
	golang.org/x/net => golang.org/x/net v0.21.0
	golang.org/x/oauth2 => golang.org/x/oauth2 v0.17.0
	golang.org/x/sync => golang.org/x/sync v0.6.0
	golang.org/x/sys => golang.org/x/sys v0.17.0
	golang.org/x/term => golang.org/x/term v0.17.0
	golang.org/x/text => golang.org/x/text v0.14.0
	golang.org/x/time => golang.org/x/time v0.5.0
	golang.org/x/tools => golang.org/x/tools v0.18.0
	golang.org/x/xerrors => golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028
	gomodules.xyz/jsonpatch/v2 => gomodules.xyz/jsonpatch/v2 v2.4.0
	gonum.org/v1/gonum => gonum.org/v1/gonum v0.14.0
	gonum.org/v1/hdf5 => gonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc6946
	gonum.org/v1/plot => gonum.org/v1/plot v0.14.0
	google.golang.org/api => google.golang.org/api v0.165.0
	google.golang.org/appengine => google.golang.org/appengine v1.6.8
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20240213162025-012b6fc9bca9
	google.golang.org/genproto/googleapis/api => google.golang.org/genproto/googleapis/api v0.0.0-20240213162025-012b6fc9bca9
	google.golang.org/genproto/googleapis/rpc => google.golang.org/genproto/googleapis/rpc v0.0.0-20240213162025-012b6fc9bca9
	google.golang.org/grpc => google.golang.org/grpc v1.61.1
	google.golang.org/grpc/cmd/protoc-gen-go-grpc => google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.3.0
	google.golang.org/protobuf => google.golang.org/protobuf v1.32.0
	gopkg.in/check.v1 => gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
	gopkg.in/inconshreveable/log15.v2 => gopkg.in/inconshreveable/log15.v2 v2.16.0
	gopkg.in/inf.v0 => gopkg.in/inf.v0 v0.9.1
	gopkg.in/tomb.v1 => gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.1
	honnef.co/go/tools => honnef.co/go/tools v0.4.6
	k8s.io/api => k8s.io/api v0.29.1
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.29.1
	k8s.io/apimachinery => k8s.io/apimachinery v0.29.1
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.29.1
	k8s.io/client-go => k8s.io/client-go v0.29.1
	k8s.io/component-base => k8s.io/component-base v0.29.1
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.120.1
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20240209001042-7a0d5b415232
	k8s.io/kubernetes => k8s.io/kubernetes v0.29.1
	k8s.io/metrics => k8s.io/metrics v0.29.1
	nhooyr.io/websocket => nhooyr.io/websocket v1.8.10
	rsc.io/pdf => rsc.io/pdf v0.1.1
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.17.0
	sigs.k8s.io/json => sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd
	sigs.k8s.io/kustomize => sigs.k8s.io/kustomize v2.0.3+incompatible
	sigs.k8s.io/structured-merge-diff/v4 => sigs.k8s.io/structured-merge-diff/v4 v4.4.1
	sigs.k8s.io/yaml => sigs.k8s.io/yaml v1.4.0
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.32.0-20240212200630-3014d81c3a48.1
	cloud.google.com/go/storage v1.37.0
	code.cloudfoundry.org/bytefmt v0.0.0-20190710193110-1eb035ffe2b6
	github.com/akrylysov/pogreb v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go v1.49.0
	github.com/fsnotify/fsnotify v1.7.0
	github.com/go-redis/redis/v8 v8.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.7.1
	github.com/goccy/go-json v0.10.2
	github.com/gocql/gocql v0.0.0-20200131111108-92af2e088537
	github.com/gocraft/dbr/v2 v2.0.0-00010101000000-000000000000
	github.com/google/go-cmp v0.6.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v0.0.0-00010101000000-000000000000
	github.com/hashicorp/go-version v0.0.0-00010101000000-000000000000
	github.com/klauspost/compress v1.16.7
	github.com/kpango/fastime v1.1.9
	github.com/kpango/fuid v0.0.0-00010101000000-000000000000
	github.com/kpango/gache/v2 v2.0.0-00010101000000-000000000000
	github.com/kpango/glg v1.6.15
	github.com/kubernetes-csi/external-snapshotter/client/v6 v6.0.0-00010101000000-000000000000
	github.com/leanovate/gopter v0.0.0-00010101000000-000000000000
	github.com/lucasb-eyer/go-colorful v0.0.0-00010101000000-000000000000
	github.com/pierrec/lz4/v3 v3.0.0-00010101000000-000000000000
	github.com/planetscale/vtprotobuf v0.6.0
	github.com/quasilyte/go-ruleguard v0.0.0-00010101000000-000000000000
	github.com/quasilyte/go-ruleguard/dsl v0.3.22
	github.com/scylladb/gocqlx v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.4
	github.com/zeebo/xxh3 v1.0.2
	go.etcd.io/bbolt v1.3.8
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.47.0
	go.opentelemetry.io/otel v1.23.1
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.23.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.19.0
	go.opentelemetry.io/otel/metric v1.23.1
	go.opentelemetry.io/otel/sdk v1.23.1
	go.opentelemetry.io/otel/sdk/metric v1.23.1
	go.opentelemetry.io/otel/trace v1.23.1
	go.uber.org/automaxprocs v0.0.0-00010101000000-000000000000
	go.uber.org/goleak v1.3.0
	go.uber.org/ratelimit v0.3.0
	go.uber.org/zap v1.26.0
	gocloud.dev v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.21.0
	golang.org/x/oauth2 v0.17.0
	golang.org/x/sync v0.6.0
	golang.org/x/sys v0.17.0
	golang.org/x/text v0.14.0
	golang.org/x/time v0.5.0
	golang.org/x/tools v0.18.0
	gonum.org/v1/hdf5 v0.0.0-00010101000000-000000000000
	gonum.org/v1/plot v0.10.1
	google.golang.org/genproto/googleapis/api v0.0.0-20240205150955-31a09d347014
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240205150955-31a09d347014
	google.golang.org/grpc v1.61.0
	google.golang.org/protobuf v1.32.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.29.1
	k8s.io/apimachinery v0.29.1
	k8s.io/cli-runtime v0.0.0-00010101000000-000000000000
	k8s.io/client-go v0.29.1
	k8s.io/metrics v0.0.0-00010101000000-000000000000
	k8s.io/utils v0.0.0-20240102154912-e7106e64919e
	sigs.k8s.io/controller-runtime v0.0.0-00010101000000-000000000000
)

require (
	cloud.google.com/go v0.112.0 // indirect
	cloud.google.com/go/compute v1.24.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v1.1.6 // indirect
	git.sr.ht/~sbinet/gg v0.5.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/ajstarks/svgo v0.0.0-20211024235047-1546f124cd8b // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/campoy/embedmd v1.0.0 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/emicklei/go-restful/v3 v3.11.0 // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.8.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/go-fonts/liberation v0.3.2 // indirect
	github.com/go-latex/latex v0.0.0-20230307184459-12ec69307ad9 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.20.2 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.6 // indirect
	github.com/go-pdf/fpdf v0.9.0 // indirect
	github.com/go-toolsmith/astcopy v1.0.2 // indirect
	github.com/go-toolsmith/astequal v1.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-sql/sqlexp v0.0.0-00010101000000-000000000000 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/wire v0.5.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20180305231024-9cad4c3443a7 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.0 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/huandu/xstrings v1.4.0 // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/moby/term v0.0.0-20221205130635-1aeaba878587 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.18.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.45.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/quasilyte/gogrep v0.5.0 // indirect
	github.com/quasilyte/stdinfo v0.0.0-20220114132959-f7386bf02567 // indirect
	github.com/scylladb/go-reflectx v1.0.1 // indirect
	github.com/spf13/cobra v1.7.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/xlab/treeprint v1.2.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.47.0 // indirect
	go.opentelemetry.io/proto/otlp v1.1.0 // indirect
	go.starlark.net v0.0.0-20230525235612-a134d8f9ddca // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d // indirect
	golang.org/x/exp/typeparams v0.0.0-20230307190834-24139beb5833 // indirect
	golang.org/x/image v0.15.0 // indirect
	golang.org/x/mod v0.15.0 // indirect
	golang.org/x/term v0.17.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
	gomodules.xyz/jsonpatch/v2 v2.4.0 // indirect
	google.golang.org/api v0.164.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto v0.0.0-20240205150955-31a09d347014 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiextensions-apiserver v0.29.0 // indirect
	k8s.io/component-base v0.29.1 // indirect
	k8s.io/klog/v2 v2.110.1 // indirect
	k8s.io/kube-openapi v0.0.0-20231010175941-2dd684a91f00 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/kustomize/api v0.13.5-0.20230601165947-6ce0bf390ce3 // indirect
	sigs.k8s.io/kustomize/kyaml v0.14.3-0.20230601165947-6ce0bf390ce3 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
