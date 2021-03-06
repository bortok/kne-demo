# Sample Ixia-c tests for KNE

## Ixia-c Traffic Generator back-2-back BGPv4 test

1. Create Ixia_TG back-2-back topology

```Shell
kne_cli create kne-demo/topologies/kne_ixia-b2b_config.txt
kubectl get pods -n ixia-c-b2b
````

2. Copy and run a test package. This package would execute one BGPv4 test

```Shell
kubectl cp kne-demo/kne-demo-tests/b2b-tests gosnappi:/go/sample-tests/
kubectl exec -it gosnappi -- /bin/bash -c "cd /go/sample-tests/b2b-tests; go test -run=TestB2BEbgpv4RoutesGosnappi -v"
````

4. Destroy Ixia_TG back-2-back topology once the testing is over

```Shell
kne_cli delete kne-demo/topologies/kne_ixia-b2b_config.txt
kubectl get pods -n ixia-c-b2b
````

## Ixia-c Traffic Generator Single DUT BGPv4 test

1. This test topology has a single DUT and two IXIA_TG ports surrounding it. The test runs eBGPv4 protocol

````
      ate1 <---------> dut <---------> ate2     
     IXIA_TG       ARISTA_CEOS        IXIA_TG   
     1.1.1.2    1.1.1.1  2.2.2.1      2.2.2.2   
     AS2222          AS1111           AS3333    
````


2. Create KNE topology with a single DUT and two Ixia_TG nodes

```Shell
kne_cli create kne-demo/topologies/kne_ixia-dut_config.txt
kubectl get pods -n ixia-c-dut
````

3. Copy and run a test package. This package would execute one BGPv4 test

```Shell
kubectl cp kne-demo/kne-demo-tests/dut-tests gosnappi:/go/sample-tests/
kubectl exec -it gosnappi -- /bin/bash -c "cd /go/sample-tests/dut-tests; go test -run=TestDUTEbgpv4RoutesGosnappi -v"
````

4. Destroy the topology once the testing is over

```Shell
kne_cli delete kne-demo/topologies/kne_ixia-dut_config.txt
kubectl get pods -n ixia-c-dut
````

##  Ixia-c 3-node Traffic Generator tests for 2-node Arista setup

1. Create Ixia_TG + Arista topology

[//]: # (TODO This relies on Arista CEOS images being present in gcr.io/kt-nts-athena-dev/ repository and access to it.)

```Shell
cd kne-demo/topologies
kne_cli create kne_ixia3_arista2_config.txt
kubectl get pods -n ixia-c
````

2. Run test BGP test package with IPv4 and IPv6 routes and traffic flows

```Shell
kubectl exec -it gosnappi -- /bin/bash -c 'cd sample-tests/tests; go test -run=TestPacketForwardBgpV4_V4V6Flows -tags=arista -v'
````

6. Destroy the Ixia_TG + Arista topology once the testing is over

```Shell
kne_cli delete kne_ixia3_arista2_config.txt
kubectl get pods -n ixia-c
cd ../../../
````

##  Ixia-c 3-node Traffic Generator tests for 3-node Arista setup

1. Create Ixia_TG + Arista topology

[//]: # (TODO This relies on Arista CEOS images being present in gcr.io/kt-nts-athena-dev/ repository and access to it.)

```Shell
cd kne-demo/topologies
kne_cli create kne_ixia-c-ceos-3node_config.txt
kubectl get pods -n ixia-c-ceos-3node
cd ../..
````

2. Copy and run a test package. This package would execute BGP test package with IPv4 and IPv6 routes and traffic flows

```Shell
kubectl cp kne-demo/kne-demo-tests/ceos-3node-tests gosnappi:/go/sample-tests/
kubectl exec -it gosnappi -- /bin/bash -c "cd /go/sample-tests/ceos-3node-tests; go test -run=Test3DUTPacketForwardBgpV4_V4V6Flows -v"
````

3. Destroy the Ixia_TG + Arista topology once the testing is over

```Shell
cd kne-demo/topologies
kne_cli delete kne_ixia-c-ceos-3node_config.txt
kubectl get pods -n ixia-c-ceos-3node
cd ../..
````

##  Ixia-c 4-node Clos POD

![Ixia-c 4-node Clos POD Diagram](images/topology_clos-4node-pod.png)

1. Create Ixia_TG + 4-node Clos POD topology

[//]: # (TODO This relies on Arista CEOS images being present in gcr.io/kt-nts-athena-dev/ repository and access to it.)

```Shell
cd kne-demo/topologies/clos-4node-pod
kne_cli create kne_clos-4node-pod-ceos.txt
kubectl get pods -n clos-4node-pod-ceos
cd ../../..
````

To inspect BGP status on the devices, use

````
kubectl exec pod1-1 -c pod1-1 -n clos-4node-pod-ceos -- Cli -e -c "sh ip bgp summary"
kubectl exec pod1-2 -c pod1-2 -n clos-4node-pod-ceos -- Cli -e -c "sh ip bgp summary"
kubectl exec tor1-1 -c tor1-1 -n clos-4node-pod-ceos -- Cli -e -c "sh ip bgp summary"
kubectl exec tor1-2 -c tor1-2 -n clos-4node-pod-ceos -- Cli -e -c "sh ip bgp summary"
````

To inspect IP routes installed from BGP, use

````
kubectl exec pod1-1 -c pod1-1 -n clos-4node-pod-ceos -- Cli -e -c "sh ip route bgp"
kubectl exec pod1-2 -c pod1-2 -n clos-4node-pod-ceos -- Cli -e -c "sh ip route bgp"
kubectl exec tor1-1 -c tor1-1 -n clos-4node-pod-ceos -- Cli -e -c "sh ip route bgp"
kubectl exec tor1-2 -c tor1-2 -n clos-4node-pod-ceos -- Cli -e -c "sh ip route bgp"
````


2. Copy and run a test package. This package would execute BGP test package with IPv4 and IPv6 routes and traffic flows, with ECMP failover

```Shell
kubectl exec gosnappi -- rm -rf /go/sample-tests/clos-4node-pod
kubectl cp kne-demo/kne-demo-tests/clos-4node-pod gosnappi:/go/sample-tests/
kubectl exec -it gosnappi -- /bin/bash -c "cd /go/sample-tests/clos-4node-pod; bash test_ecmp_failover.sh"
````

3. Copy and run a test package. This package would execute BGP test package with IPv4 and IPv6 routes and traffic flows, with BGP route-maps

```Shell
kubectl exec gosnappi -- rm -rf /go/sample-tests/clos-4node-pod
kubectl cp kne-demo/kne-demo-tests/clos-4node-pod gosnappi:/go/sample-tests/
kubectl exec -it gosnappi -- /bin/bash -c "cd /go/sample-tests/clos-4node-pod; bash test_bgp_route-map.sh"
````

3. Destroy the Ixia_TG + 4-node Clos POD topology once the testing is over

```Shell
cd kne-demo/topologies/clos-4node-pod
kne_cli delete kne_clos-4node-pod-ceos.txt
kubectl get pods -n clos-4node-pod-ceos
cd ../../..
````

