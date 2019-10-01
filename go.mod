module github.com/grpc-ecosystem/grpc-gateway

require (
	github.com/ghodss/yaml v1.0.0
	github.com/go-resty/resty v0.0.0-00010101000000-000000000000
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.2
	github.com/kr/pretty v0.1.0 // indirect
	github.com/rogpeppe/fastuuid v0.0.0-20150106093220-6724a57986af
	golang.org/x/net v0.0.0-20190930134127-c5a3c61f89f3
	golang.org/x/sys v0.0.0-20190927073244-c990c680b611 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20190927181202-20e1ac93f88c
	google.golang.org/grpc v1.24.0
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.0.0-20170812160011-eb3733d160e7 // indirect
)

replace github.com/go-resty/resty => gopkg.in/resty.v1 v1.12.0

go 1.13
