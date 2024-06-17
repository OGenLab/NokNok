module lightspeed.2dao3.com/nokserver

go 1.19

replace (
	//run `git submodule update --init --recursive` to get module 'usercenter.sdk'
	lightspeed.2dao3.com/wallet/usercenter/sdk/app v1.0.0 => ./thirdpart/usercenter.sdk/gopkg/app
	lightspeed.2dao3.com/wallet/usercenter/sdk/jwt v1.0.0 => ./thirdpart/usercenter.sdk/gopkg/jwt
	lightspeed.2dao3.com/wallet/usercenter/sdk/kms v1.0.0 => ./thirdpart/usercenter.sdk/gopkg/kms
	lightspeed.2dao3.com/wallet/usercenter/sdk/user v1.0.0 => ./thirdpart/usercenter.sdk/gopkg/user
)

require (
	github.com/ethereum/go-ethereum v1.10.25
	github.com/gomodule/redigo v1.8.9
	github.com/lonng/nano v0.5.1
	github.com/mymmrac/telego v0.30.2
	github.com/pelletier/go-toml/v2 v2.0.1
	github.com/sirupsen/logrus v1.4.2
	github.com/speps/go-hashids v2.0.0+incompatible
	gorm.io/driver/postgres v1.4.4
	gorm.io/gorm v1.25.0
	lightspeed.2dao3.com/wallet/usercenter/sdk/app v1.0.0
	qoobing.com/gomod/api v1.1.2
	qoobing.com/gomod/cache v1.3.2
	qoobing.com/gomod/log v1.2.8
	qoobing.com/gomod/redis v1.2.0
	qoobing.com/gomod/str v1.0.5
)

require (
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/bwmarrin/snowflake v0.3.0 // indirect
	github.com/bytedance/sonic v1.11.8 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/fasthttp/router v1.5.1 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.8.1 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.11.1 // indirect
	github.com/goccy/go-json v0.9.7 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grbit/go-json v0.11.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.13.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.12.0 // indirect
	github.com/jackc/pgx/v4 v4.17.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/klauspost/cpuid/v2 v2.2.6 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/prometheus/client_golang v1.19.1 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/savsgio/gotils v0.0.0-20240303185622-093b76447511 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/tylerb/gls v0.0.0-20150407001822-e606233f194d // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.54.0 // indirect
	github.com/valyala/fastjson v1.6.4 // indirect
	golang.org/x/arch v0.6.0 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/net v0.23.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto v0.0.0-20210630183607-d20f26d13c79 // indirect
	google.golang.org/grpc v1.39.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
