module github.com/nging-plugins/caddymanager

go 1.25.3

// replace github.com/coscms/webcore => ../../coscms/webcore

require (
	github.com/GehirnInc/crypt v0.0.0-20230320061759-8cc1b52080c5
	github.com/admpub/caddy v1.3.1-0.20251019083840-3f976d1025c7
	github.com/admpub/go-download/v2 v2.2.0
	github.com/admpub/goth v0.0.4
	github.com/admpub/log v1.4.0
	github.com/admpub/null v8.0.5+incompatible
	github.com/admpub/once v0.0.1
	github.com/admpub/resty/v2 v2.7.3
	github.com/admpub/tail v1.1.1
	github.com/admpub/useragent v0.0.2
	github.com/caddy-plugins/caddy-expires v1.1.3
	github.com/caddy-plugins/caddy-filter v0.15.5
	github.com/caddy-plugins/caddy-jwt/v3 v3.8.3
	github.com/caddy-plugins/caddy-locale v0.0.2
	github.com/caddy-plugins/caddy-prometheus v0.1.0
	github.com/caddy-plugins/caddy-rate-limit v1.7.0
	github.com/caddy-plugins/caddy-s3browser v0.2.2
	github.com/caddy-plugins/cors v0.0.3
	github.com/caddy-plugins/ipfilter v1.1.8
	github.com/caddy-plugins/loginsrv v0.2.6
	github.com/caddy-plugins/nobots v0.2.1
	github.com/caddy-plugins/webdav v1.3.4
	github.com/caddyserver/certmagic v0.25.1
	github.com/coscms/webcore v0.13.3-0.20251105084756-f06c33ff80c6
	github.com/spf13/cobra v1.10.2
	github.com/stretchr/testify v1.11.1
	github.com/webx-top/com v1.4.1
	github.com/webx-top/db v1.29.1
	github.com/webx-top/echo v1.22.25
	github.com/webx-top/restyclient v0.0.6
	go.uber.org/zap v1.27.1
	golang.org/x/crypto v0.47.0
	golang.org/x/net v0.48.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/abbot/go-http-auth v0.4.0 // indirect
	github.com/abh/errorutil v1.0.0 // indirect
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d // indirect
	github.com/admpub/9t v0.0.1 // indirect
	github.com/admpub/captcha-go v0.0.1 // indirect
	github.com/admpub/ccs-gm v0.0.5 // indirect
	github.com/admpub/checksum v1.1.0 // indirect
	github.com/admpub/color v1.8.1 // indirect
	github.com/admpub/confl v0.2.4 // indirect
	github.com/admpub/copier v0.1.1 // indirect
	github.com/admpub/cron v0.1.1 // indirect
	github.com/admpub/dateparse v0.0.0-20250903020633-d86d3f2a4cfd // indirect
	github.com/admpub/decimal v1.3.2 // indirect
	github.com/admpub/email v2.4.1+incompatible // indirect
	github.com/admpub/errors v0.8.2 // indirect
	github.com/admpub/events v1.3.6 // indirect
	github.com/admpub/fasthttp v0.0.7 // indirect
	github.com/admpub/fsnotify v1.7.1 // indirect
	github.com/admpub/gifresize v1.0.2 // indirect
	github.com/admpub/go-bindata-assetfs v0.0.1 // indirect
	github.com/admpub/go-captcha-assets v0.0.0-20250122071745-baa7da4bda0d // indirect
	github.com/admpub/go-captcha/v2 v2.0.7 // indirect
	github.com/admpub/go-figure v0.0.2 // indirect
	github.com/admpub/go-isatty v0.0.11 // indirect
	github.com/admpub/go-password v0.1.3 // indirect
	github.com/admpub/go-pretty/v6 v6.0.4 // indirect
	github.com/admpub/go-ps v0.0.1 // indirect
	github.com/admpub/go-reuseport v0.5.0 // indirect
	github.com/admpub/go-utility v0.0.1 // indirect
	github.com/admpub/go-zbar-wasm v0.0.0-20251219053451-e85456633a16 // indirect
	github.com/admpub/godotenv v1.4.4 // indirect
	github.com/admpub/httpscerts v0.0.0-20180907121630-a2990e2af45c // indirect
	github.com/admpub/humanize v0.0.0-20190501023926-5f826e92c8ca // indirect
	github.com/admpub/i18n v0.5.3 // indirect
	github.com/admpub/identicon v1.0.2 // indirect
	github.com/admpub/imageproxy v0.10.1 // indirect
	github.com/admpub/imaging v1.6.3 // indirect
	github.com/admpub/ip2region/v3 v3.0.4 // indirect
	github.com/admpub/json5 v0.0.1 // indirect
	github.com/admpub/license_gen v0.1.2 // indirect
	github.com/admpub/machineid v1.0.2 // indirect
	github.com/admpub/mahonia v0.0.0-20151019004008-c528b747d92d // indirect
	github.com/admpub/mail v0.0.0-20170408110349-d63147b0317b // indirect
	github.com/admpub/map2struct v0.1.3 // indirect
	github.com/admpub/mysql-schema-sync v0.2.6 // indirect
	github.com/admpub/oauth2/v4 v4.0.3 // indirect
	github.com/admpub/pp v0.0.7 // indirect
	github.com/admpub/qrcode v0.0.5 // indirect
	github.com/admpub/realip v0.2.7 // indirect
	github.com/admpub/safesvg v0.0.8 // indirect
	github.com/admpub/securecookie v1.3.0 // indirect
	github.com/admpub/service v0.0.8 // indirect
	github.com/admpub/sessions v0.3.0 // indirect
	github.com/admpub/sonyflake v0.0.1 // indirect
	github.com/admpub/sse v0.0.0-20160126180136-ee05b128a739 // indirect
	github.com/admpub/timeago v1.3.0 // indirect
	github.com/admpub/webdav/v4 v4.1.10 // indirect
	github.com/admpub/xencoding v0.0.3 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.41.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.4 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.32.6 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.19.6 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.16 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.20.18 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.95.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.41.5 // indirect
	github.com/aws/smithy-go v1.24.0 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/boombuler/barcode v1.1.0 // indirect
	github.com/caddyserver/zerossl v0.1.4 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clipperhouse/stringish v0.1.1 // indirect
	github.com/clipperhouse/uax29/v2 v2.3.0 // indirect
	github.com/coscms/captcha v0.2.3 // indirect
	github.com/coscms/forms v1.16.9 // indirect
	github.com/coscms/oauth2s v0.5.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dsoprea/go-logging v0.0.0-20200710184922-b02d349568dd // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/ebitengine/purego v0.9.1 // indirect
	github.com/fcjr/aia-transport-go v1.3.0 // indirect
	github.com/flynn/go-shlex v0.0.0-20150515145356-3f9db97f8568 // indirect
	github.com/francoispqt/gojay v1.2.13 // indirect
	github.com/fynelabs/selfupdate v0.2.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.12 // indirect
	github.com/go-acme/lego/v4 v4.31.0 // indirect
	github.com/go-errors/errors v1.5.1 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-jose/go-jose/v4 v4.1.3 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.30.1 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/sessions v1.4.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/gosimple/slug v1.15.0 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/h2non/filetype v1.1.3 // indirect
	github.com/h2non/go-is-svg v0.0.0-20160927212452-35e8c4b0612c // indirect
	github.com/hashicorp/go-syslog v1.0.0 // indirect
	github.com/hashicorp/go-version v1.8.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jimstudt/http-authentication v0.0.0-20140401203705-3eca13d6893a // indirect
	github.com/json-iterator/go v1.1.13-0.20220915233716-71ac16282d12 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/klauspost/cpuid v1.3.1 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/klauspost/crc32 v1.3.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/libdns/acmedns v0.5.0 // indirect
	github.com/libdns/cloudflare v0.2.2 // indirect
	github.com/libdns/edgeone v1.0.1 // indirect
	github.com/libdns/he v1.2.1 // indirect
	github.com/libdns/libdns v1.1.1 // indirect
	github.com/libdns/rfc2136 v1.0.1 // indirect
	github.com/libdns/tencentcloud v1.4.3 // indirect
	github.com/lufia/plan9stats v0.0.0-20251013123823-9fd1530e3ec3 // indirect
	github.com/maruel/rs v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.19 // indirect
	github.com/mholt/acmez/v3 v3.1.4 // indirect
	github.com/microcosm-cc/bluemonday v1.0.27 // indirect
	github.com/miekg/dns v1.1.69 // indirect
	github.com/minio/crc64nvme v1.1.1 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/minio-go/v7 v7.0.97 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/muesli/smartcrop v0.3.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/naoina/go-stringutil v0.1.0 // indirect
	github.com/naoina/toml v0.1.1 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/oschwald/maxminddb-golang v1.13.1 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/prometheus/client_golang v1.23.2 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.67.5 // indirect
	github.com/prometheus/procfs v0.19.2 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	github.com/quic-go/quic-go v0.58.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/russross/blackfriday v1.6.0 // indirect
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd // indirect
	github.com/segmentio/fasthash v1.0.3 // indirect
	github.com/shirou/gopsutil/v4 v4.25.12 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/tarent/lib-compose/v2 v2.0.1 // indirect
	github.com/tetratelabs/wazero v1.11.0 // indirect
	github.com/tinylib/msgp v1.6.3 // indirect
	github.com/tklauser/go-sysconf v0.3.16 // indirect
	github.com/tklauser/numcpus v0.11.0 // indirect
	github.com/tuotoo/qrcode v0.0.0-20220425170535-52ccc2bebf5d // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/vbauerster/mpb/v6 v6.0.4 // indirect
	github.com/webx-top/captcha v0.1.0 // indirect
	github.com/webx-top/chardet v0.0.2 // indirect
	github.com/webx-top/client v0.9.6 // indirect
	github.com/webx-top/codec v0.3.0 // indirect
	github.com/webx-top/image v0.1.2 // indirect
	github.com/webx-top/pagination v0.3.2 // indirect
	github.com/webx-top/poolx v0.0.0-20210912044716-5cfa2d58e380 // indirect
	github.com/webx-top/tagfast v0.0.1 // indirect
	github.com/webx-top/validation v0.0.3 // indirect
	github.com/webx-top/validator v0.3.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	github.com/zeebo/blake3 v0.2.4 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap/exp v0.3.0 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	golang.org/x/exp v0.0.0-20251219203646-944ab1f22d93 // indirect
	golang.org/x/image v0.34.0 // indirect
	golang.org/x/mod v0.31.0 // indirect
	golang.org/x/oauth2 v0.34.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	golang.org/x/time v0.14.0 // indirect
	golang.org/x/tools v0.40.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	rsc.io/qr v0.2.0 // indirect
)
