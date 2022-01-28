#!/bin/bash

go test -run=TestClosPodHosts_RemoveRouteMap -v
go test -run=TestClosPodHosts_StartProtocol -v

echo "Next is to validate end-to-end connectivity without route-maps."
read -p "Press Enter to continue..."
go test -run=TestClosPodHosts_RunTraffic -v

echo "Next is to apply a route-map for BGP connected redistribution and validate end-to-end connectivity."
read -p "Press Enter to continue..."

go test -run=TestClosPodHosts_ApplyRouteMap -v
go test -run=TestClosPodHosts_RunTraffic -v

echo "Next is to stop protocols."
read -p "Press Enter to continue..."

go test -run=TestClosPodHosts_StopProtocol -v
