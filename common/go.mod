module common

go 1.22.1

require (
	github.com/hashicorp/consul/api v1.26.1
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/rabbitmq/amqp091-go v1.10.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.53.0
	go.opentelemetry.io/otel v1.28.0
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.2
)

replace google.golang.org/genproto => google.golang.org/genproto v0.0.0-20240711142825-46eb208f015d

replace google.golang.org/protobuf => google.golang.org/protobuf v1.34.2

replace github.com/armon/go-metrics => github.com/hashicorp/go-metrics v0.5.3

require (
	github.com/armon/go-metrics v0.5.3 // indirect
	github.com/confluentinc/confluent-kafka-go/v2 v2.5.0 // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-msgpack v1.1.5 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/miekg/dns v1.1.50 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	golang.org/x/exp v0.0.0-20240707233637-46b078467d37 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240711142825-46eb208f015d // indirect
)
