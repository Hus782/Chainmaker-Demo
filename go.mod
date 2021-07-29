module myproject.com

go 1.15

require (
	chainmaker.org/chainmaker-go/common v0.0.0
	chainmaker.org/chainmaker-sdk-go v0.0.0
	github.com/ethereum/go-ethereum v1.10.5 // indirect
	github.com/tjfoc/gmtls v1.2.1 // indirect
	go.uber.org/zap v1.16.0
	google.golang.org/grpc v1.36.0
)

replace chainmaker.org/chainmaker-go/common => ./chainmaker-sdk-go/common

replace chainmaker.org/chainmaker-sdk-go => ./chainmaker-sdk-go
