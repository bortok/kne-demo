package tests

import (
	"github.com/open-traffic-generator/snappi/gosnappi"
	"keysight/athena/tests/pkg/api"
)

var opts *api.Opts

func init() {
	opts = api.NewOpts()
	// will print warning if config doesn't exist or is bad
	opts.PatchFromFile("opts.json")
	opts.BindFlags()
}

func ClosPodHostsPacketForwardBgpV4_V4FlowsConfig(client *api.ApiClient) gosnappi.Config {
	config := client.Api().NewConfig()

	// add ports
	p1 := config.Ports().Add().SetName("p1").SetLocation(opts.IxiaCPorts()[0])
	p2 := config.Ports().Add().SetName("p2").SetLocation(opts.IxiaCPorts()[1])

	// add devices
	d1 := config.Devices().Add().SetName("d1")
	d2 := config.Devices().Add().SetName("d2")

	// add flows and common properties
	for i := 1; i <= 2; i++ {
		flow := config.Flows().Add()
		flow.Metrics().SetEnable(true)
		flow.Duration().FixedPackets().SetPackets(1000)
		flow.Rate().SetPps(250)
	}

	// add protocol stacks for device d1
	d1Eth1 := d1.Ethernets().
		Add().
		SetName("d1Eth").
		SetPortName(p1.Name()).
		SetMac("00:00:01:01:01:01").
		SetMtu(1500)

	d1Eth1.
		Ipv4Addresses().
		Add().
		SetName("p1d1ipv4").
		SetAddress("169.254.1.11").
		SetGateway("169.254.1.1").
		SetPrefix(24)

	d1Bgp := d1.Bgp().
		SetRouterId("10.0.1.101")

	d1BgpIpv4Interface1 := d1Bgp.
		Ipv4Interfaces().Add().
		SetIpv4Name("p1d1ipv4")

	d1BgpIpv4Interface1Peer1 := d1BgpIpv4Interface1.
		Peers().
		Add().
		SetAsNumber(65001).
		SetAsType(gosnappi.BgpV4PeerAsType.EBGP).
		SetPeerAddress("169.254.1.1").
		SetName("BGPv4 Peer 1")

	d1BgpIpv4Interface1Peer1V4Route1 := d1BgpIpv4Interface1Peer1.
		V4Routes().
		Add().
		SetName("p1d1peer1rrv4")

	d1BgpIpv4Interface1Peer1V4Route1.Addresses().Add().
		SetAddress("10.0.1.101").
		SetPrefix(32).
		SetCount(1).
		SetStep(1)

	d1BgpIpv4Interface1Peer1V4Route1.Advanced().
		SetMultiExitDiscriminator(10).
		SetOrigin(gosnappi.BgpRouteAdvancedOrigin.IGP)

	// add protocol stacks for device d2
	d2Eth1 := d2.Ethernets().
		Add().
		SetName("d2Eth").
		SetPortName(p2.Name()).
		SetMac("00:00:02:02:02:02").
		SetMtu(1500)

	d2Eth1.
		Ipv4Addresses().
		Add().
		SetName("p2d1ipv4").
		SetAddress("169.254.1.11").
		SetGateway("169.254.1.1").
		SetPrefix(24)

	d2Bgp := d2.Bgp().
		SetRouterId("10.0.1.102")

	d2BgpIpv4Interface1 := d2Bgp.
		Ipv4Interfaces().Add().
		SetIpv4Name("p2d1ipv4")

	d2BgpIpv4Interface1Peer1 := d2BgpIpv4Interface1.
		Peers().
		Add().
		SetAsNumber(65002).
		SetAsType(gosnappi.BgpV4PeerAsType.EBGP).
		SetPeerAddress("169.254.1.1").
		SetName("BGPv4 Peer 2")

	d2BgpIpv4Interface1Peer1V4Route1 := d2BgpIpv4Interface1Peer1.
		V4Routes().
		Add().
		SetName("p2d1peer1rrv4")

	d2BgpIpv4Interface1Peer1V4Route1.Addresses().Add().
		SetAddress("10.0.1.102").
		SetPrefix(32).
		SetCount(1).
		SetStep(1)

	d2BgpIpv4Interface1Peer1V4Route1.Advanced().
		SetMultiExitDiscriminator(20).
		SetOrigin(gosnappi.BgpRouteAdvancedOrigin.IGP)

	// add endpoints and packet description flow 1
	f1 := config.Flows().Items()[0]
	f1.SetName(p1.Name() + " -> " + p2.Name() + "-IPv4").
		TxRx().Device().
		SetTxNames([]string{d1BgpIpv4Interface1Peer1V4Route1.Name()}).
		SetRxNames([]string{d2BgpIpv4Interface1Peer1V4Route1.Name()})

	f1Eth := f1.Packet().Add().Ethernet()
	f1Eth.Src().SetValue(d1Eth1.Mac())
	f1Eth.Dst().SetValue("00:00:00:00:00:00")

	f1Ip := f1.Packet().Add().Ipv4()
	f1Ip.Src().SetValue("10.0.1.101")
	f1Ip.Dst().SetValue("10.0.1.102")

	// add endpoints and packet description flow 2
	f2 := config.Flows().Items()[1]
	f2.SetName(p2.Name() + " -> " + p1.Name() + "-IPv4").
		TxRx().Device().
		SetTxNames([]string{d2BgpIpv4Interface1Peer1V4Route1.Name()}).
		SetRxNames([]string{d1BgpIpv4Interface1Peer1V4Route1.Name()})

	f2Eth := f2.Packet().Add().Ethernet()
	f2Eth.Src().SetValue(d2Eth1.Mac())
	f2Eth.Dst().SetValue("00:00:00:00:00:00")

	f2Ip := f2.Packet().Add().Ipv4()
	f2Ip.Src().SetValue("10.0.1.102")
	f2Ip.Dst().SetValue("10.0.1.101")

	return config
}
