module io.matrinos/docker/runner

go 1.16

require (
	github.com/Matrinos/iot-cadence-go-core v0.0.0-20210619125056-5cb4de703787
	github.com/containerd/containerd v1.5.2 // indirect
	github.com/docker/docker v20.10.7+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/go-resty/resty/v2 v2.6.0
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/pborman/uuid v1.2.1
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/teris-io/shortid v0.0.0-20201117134242-e59966efd125
	go.uber.org/cadence v0.17.0
	go.uber.org/zap v1.17.0
	google.golang.org/grpc v1.38.0 // indirect
)

replace github.com/apache/thrift => github.com/apache/thrift v0.0.0-20190309152529-a9b748bb0e02
