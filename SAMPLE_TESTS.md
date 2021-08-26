# Sample Athena tests for KNE

## Ixia Traffic Generator (Athena) back-2-back BGPv4 dataplane test

1. Create Ixia back-2-back topology

```Shell
./kne/kne_cli/kne_cli create kne-demo/topologies/kne_ixia-b2b_config.txt
./kne/kne_cli/kne_cli show kne-demo/topologies/kne_ixia-b2b_config.txt
kubectl get pods -n athena-dataplane
````

2. Copy and run a test package. This package would execute one BGPv4 test

```Shell
kubectl cp keysight/athena/sample-tests/otgclient test-client:/home/tests/
kubectl cp kne-demo/utils test-client:/home/tests/
kubectl cp kne-demo/kne-demo-tests test-client:/home/tests/
kubectl exec -it test-client -- /bin/bash -c "cd kne-demo-tests; go test"
````

[//]: # (TODO GAP utils.go has code specific to a topology)

4. Destroy the Ixia_TG + Arista topology once the testing is over

```Shell
./kne/kne_cli/kne_cli delete kne-demo/topologies/kne_ixia-b2b_config.txt
kubectl get pods -n athena-dataplane
````

## Arista dataplane test with Ixia Traffic Generator (Athena) - version with VLANs

1. Create Ixia_TG + Arista topology

[//]: # (TODO This relies on Arista CEOS images being present in gcr.io/kt-nts-athena-dev/ repository and access to it.)

```Shell
./kne/kne_cli/kne_cli create keysight/athena/kne/kne_config.txt
./kne/kne_cli/kne_cli show keysight/athena/kne/kne_config.txt
watch kubectl get pods -n athena-dataplane
````

  Once all PODs are running, terminate via ^C.

2. Once all the PODs the topology are running, push configuration files to Arista nodes

```Shell
./kne/kne_cli/kne_cli topology push keysight/athena/kne/kne_config.txt arista1 keysight/athena/kne/arista1_dual_config.txt
./kne/kne_cli/kne_cli topology push keysight/athena/kne/kne_config.txt arista2 keysight/athena/kne/arista2_dual_config.txt
````

[//]: # (TODO INFO[0000] Pushing config to athena-dataplane:arista1)
[//]: # (TODO Error: inappropriate ioctl for device - when running from Mac. No problem with Linux)

3. Copy a test package. This package contains two tests, one for BGPv4, with BGPv4 metrics used to pull status, and another for BGPv6, w/o use of the metrics

```Shell
kubectl exec -it test-client -- /bin/bash -c "rm -rf sample-tests"
kubectl cp keysight/athena/sample-tests test-client:/home/tests/sample-tests
````

4. Run non-raw traffic and BGPv4 metric test

```Shell
kubectl exec -it test-client -- /bin/bash -c 'cd sample-tests/tests; go test -v -run "^TestPacketForwardBgpv4$"'
````

5. Run raw traffic over BGPv6

```Shell
kubectl exec -it test-client -- /bin/bash -c 'cd sample-tests/tests; go test -v -run "^TestPacketForwardBgpv6$"'
````

4. Destroy the Ixia_TG + Arista topology once the testing is over

```Shell
./kne/kne_cli/kne_cli delete keysight/athena/kne/kne_config.txt
kubectl get pods -n athena-dataplane
````

## Arista dataplane test with Ixia Traffic Generator (Athena) - version without VLANs

1. Create Ixia_TG + Arista topology

[//]: # (TODO This relies on Arista CEOS images being present in gcr.io/kt-nts-athena-dev/ repository and access to it.)

```Shell
./kne/kne_cli/kne_cli create keysight/athena/kne/kne_config_no_vlan.txt
./kne/kne_cli/kne_cli show keysight/athena/kne/kne_config_no_vlan.txt
watch kubectl get pods -n athena-dataplane
````

  Once all PODs are running, terminate via ^C.

2. Once all the PODs the topology are running, push configuration files to Arista nodes

```Shell
./kne/kne_cli/kne_cli topology push keysight/athena/kne/kne_config_no_vlan.txt arista1 keysight/athena/kne/arista1_dual_config_no_vlan.txt
./kne/kne_cli/kne_cli topology push keysight/athena/kne/kne_config_no_vlan.txt arista2 keysight/athena/kne/arista2_dual_config_no_vlan.txt
````

[//]: # (TODO INFO[0000] Pushing config to athena-dataplane:arista1)
[//]: # (TODO Error: inappropriate ioctl for device - when running from Mac. No problem with Linux)

3. Copy a test package. This package contains two tests, one for BGPv4, with BGPv4 metrics used to pull status, and another for BGPv6, w/o use of the metrics

```Shell
kubectl exec -it test-client -- /bin/bash -c "rm -rf sample-tests"
kubectl cp keysight/athena/sample-tests test-client:/home/tests/sample-tests
kubectl exec -it test-client -- /bin/bash -c 'cd sample-tests/tests; go test -v -run "^TestPacketForwardBgpv4NoVlan$"'
kubectl exec -it test-client -- /bin/bash -c 'cd sample-tests/tests; go test -v -run "^TestPacketForwardBgpv6NoVlan$"'
````

4. Run non-raw traffic and BGPv4 metric test

```Shell
kubectl exec -it test-client -- /bin/bash -c 'cd sample-tests/tests; go test -v -run "^TestPacketForwardBgpv4NoVlan$"'
````

5. Run raw traffic test over BGPv6

```Shell
kubectl exec -it test-client -- /bin/bash -c 'cd sample-tests/tests; go test -v -run "^TestPacketForwardBgpv6NoVlan$"'
````

4. Destroy the Ixia_TG + Arista topology once the testing is over

```Shell
./kne/kne_cli/kne_cli delete keysight/athena/kne/kne_config_no_vlan.txt
kubectl get pods -n athena-dataplane
````
