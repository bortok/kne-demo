module keysight.com/athena

go 1.16

require (
	github.com/golang/protobuf v1.4.3
	github.com/openconfig/gnmi v0.0.0-20210525213403-320426956c8a
	github.com/openconfig/grpctunnel v0.0.0-20210316133510-53c78bc4915b // indirect
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/grpc v1.34.0 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	keysight.com/otgclient v0.0.0-00010101000000-000000000000
	keysight.com/utils v0.0.0-00010101000000-000000000000 // indirect
)

replace keysight.com/otgclient => ../otgclient

replace keysight.com/utils => ../utils
