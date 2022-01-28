// 4-node Clos POD eBGP v4 routes test

package tests

import (
	"keysight/athena/tests/pkg/api"
	"testing"
)

func TestClosPodHosts_RunTraffic(t *testing.T) {
	client, err := api.NewClient(opts)
	if err != nil {
		t.Fatal(err)
	}

	config := ClosPodHostsPacketForwardBgpV4_V4FlowsConfig(client)

	if err != nil {
		t.Fatal(err)
	}

	if err := client.StartTransmit(nil); err != nil {
		t.Fatal(err)
	}

	err = api.WaitFor(
		func() (bool, error) { return client.PortAndFlowMetricsOk(config) }, nil,
	)

	if err != nil {
		t.Fatal(err)
	}

}
