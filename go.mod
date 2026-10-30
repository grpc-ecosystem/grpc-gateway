module github.com/grpc-ecosystem/grpc-gateway/v2

go 1.26.0

require (
	github.com/google/go-cmp v0.7.0
	go.yaml.in/yaml/v3 v3.0.5
	golang.org/x/text v0.43.0
	google.golang.org/genproto/googleapis/api v0.0.0-20261026165753-2ffab32bce7e
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261019180522-ce908d29299b
	google.golang.org/grpc v1.85.0
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/BurntSushi/toml v1.4.1-0.20240526193622-a339e1f7089c // indirect
	golang.org/x/exp v0.0.0-20261014143641-f0251eb4aa3e // indirect
	golang.org/x/exp/typeparams v0.0.0-20231108232855-2478ac86f678 // indirect
	golang.org/x/mod v0.42.0 // indirect
	golang.org/x/net v0.62.0 // indirect
	golang.org/x/sync v0.24.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/tools v0.52.0 // indirect
	honnef.co/go/tools v0.8.2 // indirect
)

tool (
	golang.org/x/exp/cmd/gorelease
	honnef.co/go/tools/cmd/staticcheck
)
