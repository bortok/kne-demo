#!/bin/bash

go test -run=TestClosPodHosts_RemoveRouteMap
go test -run=TestClosPodHosts_StartProtocol

echo "Next is to validate end-to-end connectivity without route-maps."
read -p "Press Enter to continue..."
go test -run=TestClosPodHosts_RunTraffic

echo "Next is to apply an ingress route-map for BGP on POD switches and validate end-to-end connectivity."
read -p "Press Enter to continue..."

go test -run=TestClosPodHosts_ApplyRouteMap
go test -run=TestClosPodHosts_RunTraffic

echo "Next is to apply a FIX to the ingress route-map and validate end-to-end connectivity."
read -p "Press Enter to continue..."

go test -run=TestClosPodHosts_FixRouteMap
go test -run=TestClosPodHosts_RunTraffic

echo "Next is to stop protocols."
read -p "Press Enter to continue..."

go test -run=TestClosPodHosts_StopProtocol
