# Sample Ixia-c tests for KNE

## Ixia-c Traffic Generator back-2-back BGPv4 test

1. Create Ixia_TG back-2-back topology

```Shell
./kne/kne_cli/kne_cli create kne-demo/topologies/kne_ixia-b2b_config.txt
./kne/kne_cli/kne_cli show kne-demo/topologies/kne_ixia-b2b_config.txt
watch kubectl get pods -n ixia-c-b2b
````

2. Copy and run a test package. This package would execute one BGPv4 test

```Shell
kubectl cp kne-demo/kne-demo-tests/b2b-tests gosnappi:/go/sample-tests/
kubectl exec -it gosnappi -- /bin/bash -c "cd /go/sample-tests/b2b-tests; go test -run=TestB2BEbgpv4RoutesGosnappi -v"
````

4. Destroy Ixia_TG back-2-back topology once the testing is over

```Shell
./kne/kne_cli/kne_cli delete kne-demo/topologies/kne_ixia-b2b_config.txt
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
./kne/kne_cli/kne_cli create kne-demo/topologies/kne_ixia-dut_config.txt
./kne/kne_cli/kne_cli show kne-demo/topologies/kne_ixia-dut_config.txt
watch kubectl get pods -n ixia-c-dut
````

3. Copy and run a test package. This package would execute one BGPv4 test

```Shell
kubectl cp kne-demo/kne-demo-tests/dut-tests gosnappi:/go/sample-tests/
kubectl exec -it gosnappi -- /bin/bash -c "cd /go/sample-tests/dut-tests; go test -run=TestDUTEbgpv4RoutesGosnappi -v"
````

4. Destroy the topology once the testing is over

```Shell
./kne/kne_cli/kne_cli delete kne-demo/topologies/kne_ixia-dut_config.txt
kubectl get pods -n ixia-c-dut
````

##  Ixia-c 3-node Traffic Generator tests for 2-node Arista setup

1. Create Ixia_TG + Arista topology

[//]: # (TODO This relies on Arista CEOS images being present in gcr.io/kt-nts-athena-dev/ repository and access to it.)

```Shell
cd keysight/athena/kne/
../../../kne/kne_cli/kne_cli create kne_config.txt
../../../kne/kne_cli/kne_cli show kne_config.txt
watch kubectl get pods -n ixia-c
````

  Once all PODs are running, terminate via ^C.

2. Run test BGP test package with IPv4 and IPv6 routes and traffic flows

```Shell
kubectl exec -it gosnappi -- /bin/bash -c 'cd sample-tests/tests; go test -run=TestPacketForwardBgpV4_V4V6Flows -tags=arista -v'
````

6. Destroy the Ixia_TG + Arista topology once the testing is over

```Shell
../../../kne/kne_cli/kne_cli delete kne_config.txt
kubectl get pods -n ixia-c
cd ../../../
````

