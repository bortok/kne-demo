// 3 DUT eBGP v4/v6 routes test

package tests

import (
	"keysight/athena/tests/pkg/api"
	"log"
	"testing"

	"github.com/open-traffic-generator/snappi/gosnappi"
)

func Test3DUTPacketForwardBgpV4_V4V6Flows(t *testing.T) {
	client, err := api.NewClient(opts)
	if err != nil {
		t.Fatal(err)
	}

	defer client.Close()
	defer client.StopProtocol()

	config := PacketForwardBgpv4_v4v6_RoutesConfig(client)

	if err != nil {
		t.Fatal(err)
	}

	log.Println("Logging into DUTs ...")
	for i, location := range opts.DutPorts() {
		dut, err := api.NewSshClient(opts, location, "admin")
		if err != nil {
			t.Fatal(err)
		}
		defer dut.Close()

		if i == 0 {
			if _, err := dut.PushDutConfigFile("./configs/set_arista1.txt"); err != nil {
				t.Fatal(err)
			}
			defer dut.PushDutConfigFile("./configs/unset_arista.txt")
		} else if i == 1 {
			if _, err := dut.PushDutConfigFile("./configs/set_arista2.txt"); err != nil {
				t.Fatal(err)
			}
			defer dut.PushDutConfigFile("./configs/unset_arista.txt")
		} else if i == 2 {
			if _, err := dut.PushDutConfigFile("./configs/set_arista3.txt"); err != nil {
				t.Fatal(err)
			}
			defer dut.PushDutConfigFile("./configs/unset_arista.txt")
		} else {
			break
		}
	}

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
}

func PacketForwardBgpv4_v4v6_RoutesConfig(client *api.ApiClient) gosnappi.Config {
	config := client.Api().NewConfig()

	// add ports
	p1 := config.Ports().Add().SetName("p1").SetLocation(opts.IxiaCPorts()[0])
	p2 := config.Ports().Add().SetName("p2").SetLocation(opts.IxiaCPorts()[1])
	p3 := config.Ports().Add().SetName("p3").SetLocation(opts.IxiaCPorts()[2])

	// add devices
	d1 := config.Devices().Add().SetName("d1")
	d2 := config.Devices().Add().SetName("d2")
	d3 := config.Devices().Add().SetName("d3")

	// add flows and common properties
	for i := 1; i <= 8; i++ {
		flow := config.Flows().Add()
		flow.Metrics().SetEnable(true)
		flow.Duration().FixedPackets().SetPackets(1000)
		flow.Rate().SetPps(500)
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
		SetAddress("1.1.1.1").
		SetGateway("1.1.1.2").
		SetPrefix(24)
	//The IPv6 address is needed for Arista to forward IPv6 data-packets
	//based on the IPv6 route injected into Arista using BGPv4 session.
	d1Eth1.
		Ipv6Addresses().
		Add().
		SetName("p1d1ipv6").
		SetAddress("0:1:1:1::1").
		SetGateway("0:1:1:1::2").
		SetPrefix(64)

	d1Bgp := d1.Bgp().
		SetRouterId("1.1.1.1")

	d1BgpIpv4Interface1 := d1Bgp.
		Ipv4Interfaces().Add().
		SetIpv4Name("p1d1ipv4")

	d1BgpIpv4Interface1Peer1 := d1BgpIpv4Interface1.
		Peers().
		Add().
		SetAsNumber(1111).
		SetAsType(gosnappi.BgpV4PeerAsType.EBGP).
		SetPeerAddress("1.1.1.2").
		SetName("BGPv4 Peer 1")

	d1BgpIpv4Interface1Peer1V4Route1 := d1BgpIpv4Interface1Peer1.
		V4Routes().
		Add().
		SetName("p1d1peer1rrv4")

	d1BgpIpv4Interface1Peer1V4Route1.Addresses().Add().
		SetAddress("10.10.10.0").
		SetPrefix(24).
		SetCount(2).
		SetStep(2)

	d1BgpIpv4Interface1Peer1V4Route1.Advanced().
		SetMultiExitDiscriminator(10).
		SetOrigin(gosnappi.BgpRouteAdvancedOrigin.IGP)

	d1BgpIpv4Interface1Peer1V6Route1 := d1BgpIpv4Interface1Peer1.
		V6Routes().
		Add().
		SetNextHopIpv6Address("0:1:1:1::1").
		SetName("p1d1peer1rrv6").
		SetNextHopAddressType(gosnappi.BgpV6RouteRangeNextHopAddressType.IPV6).
		SetNextHopMode(gosnappi.BgpV6RouteRangeNextHopMode.MANUAL)

	d1BgpIpv4Interface1Peer1V6Route1.Addresses().Add().
		SetAddress("0:10:10:10::0").
		SetPrefix(64).
		SetCount(2).
		SetStep(2)

	d1BgpIpv4Interface1Peer1V6Route1.Advanced().
		SetMultiExitDiscriminator(50).
		SetOrigin(gosnappi.BgpRouteAdvancedOrigin.EGP)

	d1BgpIpv4Interface1Peer1V6Route1.Communities().Add().
		SetAsNumber(1).
		SetAsCustom(2).
		SetType(gosnappi.BgpCommunityType.MANUAL_AS_NUMBER)

	d1BgpIpv4Interface1Peer1V6Route1AsPath := d1BgpIpv4Interface1Peer1V6Route1.AsPath().
		SetAsSetMode(gosnappi.BgpAsPathAsSetMode.INCLUDE_AS_SET)

	d1BgpIpv4Interface1Peer1V6Route1AsPath.Segments().Add().
		SetAsNumbers([]int64{1112, 1113}).
		SetType(gosnappi.BgpAsPathSegmentType.AS_SEQ)

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
		SetAddress("2.2.2.2").
		SetGateway("2.2.2.1").
		SetPrefix(24)

	d2Eth1.
		Ipv6Addresses().
		Add().
		SetName("p2d1ipv6").
		SetAddress("0:2:2:2::2").
		SetGateway("0:2:2:2::1").
		SetPrefix(64)

	d2Bgp := d2.Bgp().
		SetRouterId("2.2.2.2")

	d2BgpIpv4Interface1 := d2Bgp.
		Ipv4Interfaces().Add().
		SetIpv4Name("p2d1ipv4")

	d2BgpIpv4Interface1Peer1 := d2BgpIpv4Interface1.
		Peers().
		Add().
		SetAsNumber(2222).
		SetAsType(gosnappi.BgpV4PeerAsType.EBGP).
		SetPeerAddress("2.2.2.1").
		SetName("BGPv4 Peer 2")

	d2BgpIpv4Interface1Peer1V4Route1 := d2BgpIpv4Interface1Peer1.
		V4Routes().
		Add().
		SetName("p2d1peer1rrv4")

	d2BgpIpv4Interface1Peer1V4Route1.Addresses().Add().
		SetAddress("20.20.20.0").
		SetPrefix(24).
		SetCount(2).
		SetStep(2)

	d2BgpIpv4Interface1Peer1V4Route1.Advanced().
		SetMultiExitDiscriminator(20).
		SetOrigin(gosnappi.BgpRouteAdvancedOrigin.IGP)

	d2BgpIpv4Interface1Peer1V6Route1 := d2BgpIpv4Interface1Peer1.
		V6Routes().
		Add().
		SetNextHopIpv6Address("0:2:2:2::2").
		SetName("p2d1peer1rrv6").
		SetNextHopAddressType(gosnappi.BgpV6RouteRangeNextHopAddressType.IPV6).
		SetNextHopMode(gosnappi.BgpV6RouteRangeNextHopMode.MANUAL)

	d2BgpIpv4Interface1Peer1V6Route1.Addresses().Add().
		SetAddress("0:20:20:20::0").
		SetPrefix(64).
		SetCount(2).
		SetStep(2)

	d2BgpIpv4Interface1Peer1V6Route1.Advanced().
		SetMultiExitDiscriminator(40).
		SetOrigin(gosnappi.BgpRouteAdvancedOrigin.EGP)

	d2BgpIpv4Interface1Peer1V6Route1.Communities().Add().
		SetAsNumber(100).
		SetAsCustom(2).
		SetType(gosnappi.BgpCommunityType.MANUAL_AS_NUMBER)

	d2BgpIpv4Interface1Peer1V6Route1AsPath := d2BgpIpv4Interface1Peer1V6Route1.AsPath().
		SetAsSetMode(gosnappi.BgpAsPathAsSetMode.INCLUDE_AS_SET)

	d2BgpIpv4Interface1Peer1V6Route1AsPath.Segments().Add().
		SetAsNumbers([]int64{2223, 2224, 2225}).
		SetType(gosnappi.BgpAsPathSegmentType.AS_SEQ)

	// add protocol stacks for device d3

	d3Eth1 := d3.Ethernets().
		Add().
		SetName("d3Eth").
		SetPortName(p3.Name()).
		SetMac("00:00:03:03:03:02").
		SetMtu(1500)

	d3Eth1.
		Ipv4Addresses().
		Add().
		SetName("p3d1ipv4").
		SetAddress("3.3.3.2").
		SetGateway("3.3.3.1").
		SetPrefix(24)

	d3Eth1.
		Ipv6Addresses().
		Add().
		SetName("p3d1ipv6").
		SetAddress("0:3:3:3::2").
		SetGateway("0:3:3:3::1").
		SetPrefix(64)

	d3Bgp := d3.Bgp().
		SetRouterId("3.3.3.2")

	d3BgpIpv4Interface1 := d3Bgp.
		Ipv4Interfaces().Add().
		SetIpv4Name("p3d1ipv4")

	d3BgpIpv4Interface1Peer1 := d3BgpIpv4Interface1.
		Peers().
		Add().
		SetAsNumber(3332).
		SetAsType(gosnappi.BgpV4PeerAsType.EBGP).
		SetPeerAddress("3.3.3.1").
		SetName("BGPv4 Peer 3")

	d3BgpIpv4Interface1Peer1V4Route1 := d3BgpIpv4Interface1Peer1.
		V4Routes().
		Add().
		SetName("p3d1peer1rrv4")

	d3BgpIpv4Interface1Peer1V4Route1.Addresses().Add().
		SetAddress("30.30.30.0").
		SetPrefix(24).
		SetCount(2).
		SetStep(2)

	d3BgpIpv4Interface1Peer1V4Route1.Advanced().
		SetMultiExitDiscriminator(30).
		SetOrigin(gosnappi.BgpRouteAdvancedOrigin.IGP)

	d3BgpIpv4Interface1Peer1V6Route1 := d3BgpIpv4Interface1Peer1.
		V6Routes().
		Add().
		SetNextHopIpv6Address("0:3:3:3::2").
		SetName("p3d1peer1rrv6").
		SetNextHopAddressType(gosnappi.BgpV6RouteRangeNextHopAddressType.IPV6).
		SetNextHopMode(gosnappi.BgpV6RouteRangeNextHopMode.MANUAL)

	d3BgpIpv4Interface1Peer1V6Route1.Addresses().Add().
		SetAddress("0:30:30:30::0").
		SetPrefix(64).
		SetCount(2).
		SetStep(2)

	d3BgpIpv4Interface1Peer1V6Route1.Advanced().
		SetMultiExitDiscriminator(33).
		SetOrigin(gosnappi.BgpRouteAdvancedOrigin.EGP)

	d3BgpIpv4Interface1Peer1V6Route1.Communities().Add().
		SetAsNumber(1).
		SetAsCustom(2).
		SetType(gosnappi.BgpCommunityType.MANUAL_AS_NUMBER)

	d3BgpIpv4Interface1Peer1V6Route1AsPath := d3BgpIpv4Interface1Peer1V6Route1.AsPath().
		SetAsSetMode(gosnappi.BgpAsPathAsSetMode.INCLUDE_AS_SET)

	d3BgpIpv4Interface1Peer1V6Route1AsPath.Segments().Add().
		SetAsNumbers([]int64{3333, 3334}).
		SetType(gosnappi.BgpAsPathSegmentType.AS_SEQ)

	// add endpoints and packet description flow 1
	f1 := config.Flows().Items()[0]
	f1.SetName(p1.Name() + " -> " + p2.Name() + "-IPv6").
		TxRx().Device().
		SetTxNames([]string{d1BgpIpv4Interface1Peer1V6Route1.Name()}).
		SetRxNames([]string{d2BgpIpv4Interface1Peer1V6Route1.Name()})

	f1Eth := f1.Packet().Add().Ethernet()
	f1Eth.Src().SetValue(d1Eth1.Mac())
	f1Eth.Dst().SetValue("00:00:00:00:00:00")

	f1Ip := f1.Packet().Add().Ipv6()
	f1Ip.Src().SetValue("0:10:10:10::1")
	f1Ip.Dst().SetValue("0:20:20:20::1")

	// add endpoints and packet description flow 2
	f2 := config.Flows().Items()[1]
	f2.SetName(p1.Name() + " -> " + p3.Name() + "-IPv6").
		TxRx().Device().
		SetTxNames([]string{d1BgpIpv4Interface1Peer1V6Route1.Name()}).
		SetRxNames([]string{d3BgpIpv4Interface1Peer1V6Route1.Name()})

	f2Eth := f2.Packet().Add().Ethernet()
	f2Eth.Src().SetValue(d1Eth1.Mac())
	f2Eth.Dst().SetValue("00:00:00:00:00:00")

	f2Ip := f2.Packet().Add().Ipv6()
	f2Ip.Src().SetValue("0:10:10:10::1")
	f2Ip.Dst().SetValue("0:30:30:30::1")

	// add endpoints and packet description flow 3
	f3 := config.Flows().Items()[2]
	f3.SetName(p2.Name() + " -> " + p1.Name() + "-IPv6").
		TxRx().Device().
		SetTxNames([]string{d2BgpIpv4Interface1Peer1V6Route1.Name()}).
		SetRxNames([]string{d1BgpIpv4Interface1Peer1V6Route1.Name()})

	f3Eth := f3.Packet().Add().Ethernet()
	f3Eth.Src().SetValue(d2Eth1.Mac())
	f3Eth.Dst().SetValue("00:00:00:00:00:00")

	f3Ip := f3.Packet().Add().Ipv6()
	f3Ip.Src().SetValue("0:20:20:20::1")
	f3Ip.Dst().SetValue("0:10:10:10::1")

	// add endpoints and packet description flow 4
	f4 := config.Flows().Items()[3]
	f4.SetName(p3.Name() + " -> " + p1.Name() + "-IPv6").
		TxRx().Device().
		SetTxNames([]string{d3BgpIpv4Interface1Peer1V6Route1.Name()}).
		SetRxNames([]string{d1BgpIpv4Interface1Peer1V6Route1.Name()})

	f4Eth := f4.Packet().Add().Ethernet()
	f4Eth.Src().SetValue(d3Eth1.Mac())
	f4Eth.Dst().SetValue("00:00:00:00:00:00")

	f4Ip := f4.Packet().Add().Ipv6()
	f4Ip.Src().SetValue("0:30:30:30::1")
	f4Ip.Dst().SetValue("0:10:10:10::1")

	// add endpoints and packet description flow 5
	f5 := config.Flows().Items()[4]
	f5.SetName(p1.Name() + " -> " + p2.Name() + "-IPv4").
		TxRx().Device().
		SetTxNames([]string{d1BgpIpv4Interface1Peer1V4Route1.Name()}).
		SetRxNames([]string{d2BgpIpv4Interface1Peer1V4Route1.Name()})

	f5Eth := f5.Packet().Add().Ethernet()
	f5Eth.Src().SetValue(d1Eth1.Mac())
	f5Eth.Dst().SetValue("00:00:00:00:00:00")

	f5Ip := f5.Packet().Add().Ipv4()
	f5Ip.Src().SetValue("10.10.10.1")
	f5Ip.Dst().SetValue("20.20.20.1")

	// add endpoints and packet description flow 6
	f6 := config.Flows().Items()[5]
	f6.SetName(p1.Name() + " -> " + p3.Name() + "-IPv4").
		TxRx().Device().
		SetTxNames([]string{d1BgpIpv4Interface1Peer1V4Route1.Name()}).
		SetRxNames([]string{d3BgpIpv4Interface1Peer1V4Route1.Name()})

	f6Eth := f6.Packet().Add().Ethernet()
	f6Eth.Src().SetValue(d1Eth1.Mac())
	f6Eth.Dst().SetValue("00:00:00:00:00:00")

	f6Ip := f6.Packet().Add().Ipv4()
	f6Ip.Src().SetValue("10.10.10.1")
	f6Ip.Dst().SetValue("30.30.30.1")

	// add endpoints and packet description flow 7
	f7 := config.Flows().Items()[6]
	f7.SetName(p2.Name() + " -> " + p1.Name() + "-IPv4").
		TxRx().Device().
		SetTxNames([]string{d2BgpIpv4Interface1Peer1V4Route1.Name()}).
		SetRxNames([]string{d1BgpIpv4Interface1Peer1V4Route1.Name()})

	f7Eth := f7.Packet().Add().Ethernet()
	f7Eth.Src().SetValue(d2Eth1.Mac())
	f7Eth.Dst().SetValue("00:00:00:00:00:00")

	f7Ip := f7.Packet().Add().Ipv4()
	f7Ip.Src().SetValue("20.20.20.1")
	f7Ip.Dst().SetValue("10.10.10.1")

	// add endpoints and packet description flow 8
	f8 := config.Flows().Items()[7]
	f8.SetName(p3.Name() + " -> " + p1.Name() + "-IPv4").
		TxRx().Device().
		SetTxNames([]string{d3BgpIpv4Interface1Peer1V4Route1.Name()}).
		SetRxNames([]string{d1BgpIpv4Interface1Peer1V4Route1.Name()})

	f8Eth := f8.Packet().Add().Ethernet()
	f8Eth.Src().SetValue(d3Eth1.Mac())
	f8Eth.Dst().SetValue("00:00:00:00:00:00")

	f8Ip := f8.Packet().Add().Ipv4()
	f8Ip.Src().SetValue("30.30.30.1")
	f8Ip.Dst().SetValue("10.10.10.1")

	return config
}
