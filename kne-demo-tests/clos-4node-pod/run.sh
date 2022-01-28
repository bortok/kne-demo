#!/bin/bash

go test -run=TestClosPodHostsPacketForwardBgpV4_V4Flows -v
sleep 5
go test -run=TestClosPodHosts_ShutEth1 -v
sleep 5
go test -run=TestClosPodHostsPacketForwardBgpV4_V4Flows -v
sleep 5
go test -run=TestClosPodHosts_NoShutEth1 -v
sleep 5
go test -run=TestClosPodHostsPacketForwardBgpV4_V4Flows -v
