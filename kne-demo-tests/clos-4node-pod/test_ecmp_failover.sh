#!/bin/bash

go test -run=TestClosPodHostsPacketForwardBgpV4_V4Flows -v

echo "Next is to shutdown TOR links to POD1-1 and check if traffic would still flow."
read -p "Press Enter to continue..."

go test -run=TestClosPodHosts_ShutEth1 -v
go test -run=TestClosPodHostsPacketForwardBgpV4_V4Flows -v

echo "Next is to bring TOR links to POD1-1 back up and check if traffic would still flow."
read -p "Press Enter to continue..."

go test -run=TestClosPodHosts_NoShutEth1 -v
go test -run=TestClosPodHostsPacketForwardBgpV4_V4Flows -v
