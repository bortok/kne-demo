// 4-node Clos POD eBGP v4 routes test

package tests

import (
	"keysight/athena/tests/pkg/api"
	"testing"
)

func TestClosPodHosts_StartProtocol(t *testing.T) {
	client, err := api.NewClient(opts)
	if err != nil {
		t.Fatal(err)
	}

	config := ClosPodHostsPacketForwardBgpV4_V4FlowsConfig(client)

	if err != nil {
		t.Fatal(err)
	}

	if err := client.SetConfig(config); err != nil {
		t.Fatal(err)
	}

	if err := client.StartProtocol(); err != nil {
		t.Fatal(err)
	}

	err = api.WaitFor(
		func() (bool, error) {
			return client.Bgp4MetricsAsExpected(config, []api.BgpMetric{
				{
					Name:             "BGPv4 Peer 1",
					Up:               true,
					SessionFlaps:     0,
					RoutesTx:         1,
					RoutesRx:         10,
					RouteWithdrawsTx: 0,
					RouteWithdrawsRx: 0,
				},
				{
					Name:             "BGPv4 Peer 2",
					Up:               true,
					SessionFlaps:     0,
					RoutesTx:         1,
					RoutesRx:         10,
					RouteWithdrawsTx: 0,
					RouteWithdrawsRx: 0,
				},
			})
		}, nil,
	)

	if err != nil {
		t.Fatal(err)
	}
}
