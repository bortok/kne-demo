// Single DUT eBGP v4 routes test

package tests

import (
	"keysight/athena/tests/pkg/api"
	"log"
	"testing"

	"github.com/open-traffic-generator/snappi/gosnappi"
)

func TestDUTEbgpv4RoutesGosnappi(t *testing.T) {
	client, err := api.NewClient(opts)
	if err != nil {
		t.Fatal(err)
	}

	defer client.Close()
	defer client.StopProtocol()

	config := Bgpv4RoutesConfig(client)

	if err != nil {
		t.Fatal(err)
	}

	log.Println("Logging into DUT ...")
	dut, err := api.NewSshClient(opts, opts.DutPorts()[0], "admin", "")
	if err != nil {
		t.Fatal(err)
	}
	defer dut.Close()

	if _, err := dut.PushDutConfigFile("./set_arista_ebgpv4.txt"); err != nil {
		t.Fatal(err)
	}
	defer dut.PushDutConfigFile("./unset_arista_ebgpv4.txt")

	if err := client.SetConfig(config); err != nil {
		t.Fatal(err)
	}

	if err := client.StartProtocol(); err != nil {
		t.Fatal(err)
	}

	err = api.WaitFor(
		func() (bool, error) { return client.AllBgp4SessionUp(config) }, nil,
	)

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

	gnmiClient, err := api.NewGnmiClient(opts, config)
	if err != nil {
		t.Fatal(err)
	}
	defer gnmiClient.Close()

	err = api.WaitFor(
		func() (bool, error) { return gnmiClient.AllBgp4SessionUp(config) }, nil,
	)

	if err != nil {
		t.Error(err)
		return
	}

	err = api.WaitFor(
		func() (bool, error) { return gnmiClient.PortAndFlowMetricsOk(config) }, nil,
	)

	if err != nil {
		t.Fatal(err)
	}

}

func Bgpv4RoutesConfig(client *api.ApiClient) gosnappi.Config {
	config := client.Api().NewConfig()

	// add ports
	p1 := config.Ports().Add().SetName("p1").SetLocation(opts.IxiaCPorts()[0])
	p2 := config.Ports().Add().SetName("p2").SetLocation(opts.IxiaCPorts()[1])

	// add devices
	d1 := config.Devices().Add().SetName("p1d1")
	d2 := config.Devices().Add().SetName("p2d1")

	// add flows and common properties
	for i := 1; i <= 2; i++ {
		flow := config.Flows().Add()
		flow.Metrics().SetEnable(true)
		flow.Duration().FixedPackets().SetPackets(1000)
		flow.Rate().SetPps(500)
	}

	// add protocol stacks for device d1
	d1Eth1 := d1.Ethernets().
		Add().
		SetName("p1d1eth").
		SetPortName(p1.Name()).
		SetMac("00:00:01:01:01:01").
		SetMtu(1500)

	d1Eth1.
		Ipv4Addresses().
		Add().
		SetName("p1d1ipv4").
		SetAddress("1.1.1.2").
		SetGateway("1.1.1.1").
		SetPrefix(24)

	d1Bgp := d1.Bgp().
		SetRouterId("1.1.1.2")

	d1BgpIpv4Interface1 := d1Bgp.
		Ipv4Interfaces().Add().
		SetIpv4Name("p1d1ipv4")

	d1BgpIpv4Interface1Peer1 := d1BgpIpv4Interface1.
		Peers().
		Add().
		SetAsNumber(2222).
		SetAsType(gosnappi.BgpV4PeerAsType.EBGP).
		SetPeerAddress("1.1.1.1").
		SetName("p1d1bgpv4")

	d1BgpIpv4Interface1Peer1V4Route1 := d1BgpIpv4Interface1Peer1.
		V4Routes().
		Add().
		SetNextHopIpv4Address("1.1.1.2").
		SetName("p1d1peer1rrv4").
		SetNextHopAddressType(gosnappi.BgpV4RouteRangeNextHopAddressType.IPV4).
		SetNextHopMode(gosnappi.BgpV4RouteRangeNextHopMode.MANUAL)

	d1BgpIpv4Interface1Peer1V4Route1.Addresses().Add().
		SetAddress("10.10.10.1").
		SetPrefix(32).
		SetCount(4).
		SetStep(1)

	d1BgpIpv4Interface1Peer1V4Route1.Advanced().
		SetMultiExitDiscriminator(50).
		SetOrigin(gosnappi.BgpRouteAdvancedOrigin.EGP)

	d1BgpIpv4Interface1Peer1V4Route1.Communities().Add().
		SetAsNumber(1).
		SetAsCustom(2).
		SetType(gosnappi.BgpCommunityType.MANUAL_AS_NUMBER)

	d1BgpIpv4Interface1Peer1V4Route1AsPath := d1BgpIpv4Interface1Peer1V4Route1.AsPath().
		SetAsSetMode(gosnappi.BgpAsPathAsSetMode.INCLUDE_AS_SET)

	d1BgpIpv4Interface1Peer1V4Route1AsPath.Segments().Add().
		SetAsNumbers([]int64{1112, 1113}).
		SetType(gosnappi.BgpAsPathSegmentType.AS_SEQ)

	// add protocol stacks for device d2
	d2Eth1 := d2.Ethernets().
		Add().
		SetName("p2d1eth").
		SetPortName(p2.Name()).
		SetMac("00:00:02:02:02:02").
		SetMtu(1500)

	d2Eth1.
		Ipv4Addresses().
		Add().
		SetName("p2d1ipv4").
		SetAddress("2.2.2.2").
		SetGateway("2.2.2.1").
		SetPrefix(32)

	d2Bgp := d2.Bgp().
		SetRouterId("2.2.2.2")

	d2BgpIpv4Interface1 := d2Bgp.
		Ipv4Interfaces().Add().
		SetIpv4Name("p2d1ipv4")

	d2BgpIpv4Interface1Peer1 := d2BgpIpv4Interface1.
		Peers().
		Add().
		SetAsNumber(3333).
		SetAsType(gosnappi.BgpV4PeerAsType.EBGP).
		SetPeerAddress("2.2.2.1").
		SetName("p2d1bgpv4")

	d2BgpIpv4Interface1Peer1V4Route1 := d2BgpIpv4Interface1Peer1.
		V4Routes().
		Add().
		SetNextHopIpv4Address("2.2.2.2").
		SetName("p2d1peer1rrv4").
		SetNextHopAddressType(gosnappi.BgpV4RouteRangeNextHopAddressType.IPV4).
		SetNextHopMode(gosnappi.BgpV4RouteRangeNextHopMode.MANUAL)

	d2BgpIpv4Interface1Peer1V4Route1.Addresses().Add().
		SetAddress("20.20.20.1").
		SetPrefix(32).
		SetCount(2).
		SetStep(2)

	d2BgpIpv4Interface1Peer1V4Route1.Advanced().
		SetMultiExitDiscriminator(40).
		SetOrigin(gosnappi.BgpRouteAdvancedOrigin.EGP)

	d2BgpIpv4Interface1Peer1V4Route1.Communities().Add().
		SetAsNumber(100).
		SetAsCustom(2).
		SetType(gosnappi.BgpCommunityType.MANUAL_AS_NUMBER)

	d2BgpIpv4Interface1Peer1V4Route1AsPath := d2BgpIpv4Interface1Peer1V4Route1.AsPath().
		SetAsSetMode(gosnappi.BgpAsPathAsSetMode.INCLUDE_AS_SET)

	d2BgpIpv4Interface1Peer1V4Route1AsPath.Segments().Add().
		SetAsNumbers([]int64{2223, 2224, 2225}).
		SetType(gosnappi.BgpAsPathSegmentType.AS_SEQ)

	// add endpoints and packet description flow 1
	f1 := config.Flows().Items()[0]
	f1.SetName(p1.Name() + " -> " + p2.Name()).
		TxRx().Device().
		SetTxNames([]string{d1BgpIpv4Interface1Peer1V4Route1.Name()}).
		SetRxNames([]string{d2BgpIpv4Interface1Peer1V4Route1.Name()})

	f1Eth := f1.Packet().Add().Ethernet()
	f1Eth.Src().SetValue(d1Eth1.Mac())
	f1Eth.Dst().SetValue("00:00:00:00:00:00")

	f1Ip := f1.Packet().Add().Ipv4()
	f1Ip.Src().SetValue("10.10.10.1")
	f1Ip.Dst().SetValue("20.20.20.1")

	// add endpoints and packet description flow 2
	f2 := config.Flows().Items()[1]
	f2.SetName(p2.Name() + " -> " + p1.Name()).
		TxRx().Device().
		SetTxNames([]string{d2BgpIpv4Interface1Peer1V4Route1.Name()}).
		SetRxNames([]string{d1BgpIpv4Interface1Peer1V4Route1.Name()})

	f2Eth := f2.Packet().Add().Ethernet()
	f2Eth.Src().SetValue(d2Eth1.Mac())
	f2Eth.Dst().SetValue("00:00:00:00:00:00")

	f2Ip := f2.Packet().Add().Ipv4()
	f2Ip.Src().SetValue("20.20.20.1")
	f2Ip.Dst().SetValue("10.10.10.1")

	return config
}
