#!/bin/bash

go test -run=TestClosPodHosts_RemoveRouteMap -v
go test -run=TestClosPodHostsPacketForwardBgpV4_V4Flows -v

echo "Next is to apply a route-map for BGP connected redistribution and validate end-to-end connectivity."
read -p "Press Enter to continue..."

go test -run=TestClosPodHosts_ApplyRouteMap -v
go test -run=TestClosPodHostsPacketForwardBgpV4_V4Flows -v

