module github.com/guackamolly/zero-monitor

go 1.23.1

require (
	github.com/denisbrodbeck/machineid v1.0.1
	github.com/go-zeromq/zmq4 v0.17.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	github.com/jaypipes/ghw v0.13.0
	github.com/joho/godotenv v1.5.1
	github.com/labstack/echo/v4 v4.12.0
	github.com/mssola/useragent v1.0.0
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/showwin/speedtest-go v1.7.9
	github.com/wcharczuk/go-chart/v2 v2.1.2
	go.etcd.io/bbolt v1.3.11
)

require (
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-zeromq/goczmq/v4 v4.2.2 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/jaypipes/pcidb v1.0.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/tklauser/go-sysconf v0.3.14 // indirect
	github.com/tklauser/numcpus v0.8.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/image v0.18.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	howett.net/plist v1.0.0 // indirect
)

// Fork of go-chart to add tooltip to dots.
replace github.com/wcharczuk/go-chart/v2 => github.com/freitzzz/go-chart/v2 v2.0.0-20241111124638-827fb77786ad
