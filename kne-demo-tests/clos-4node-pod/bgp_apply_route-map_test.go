// Apply a route-map to BGP redistribute connected statement

package tests

import (
	"keysight/athena/tests/pkg/api"
	"testing"
)

func TestClosPodHosts_ApplyRouteMap(t *testing.T) {
	for i, location := range opts.DutPorts() {
		dut, err := api.NewSshClient(opts, location, "admin")
		if err != nil {
			t.Fatal(err)
		}
		defer dut.Close()

		if i == 0 {
			if _, err := dut.PushDutConfigFile("./configs/apply_route-map_tor1-1.txt"); err != nil {
				t.Fatal(err)
			}
		} else if i == 1 {
			if _, err := dut.PushDutConfigFile("./configs/apply_route-map_tor1-2.txt"); err != nil {
				t.Fatal(err)
			}
		} else {
			break
		}
	}
}

func TestClosPodHosts_RemoveRouteMap(t *testing.T) {
	for i, location := range opts.DutPorts() {
		dut, err := api.NewSshClient(opts, location, "admin")
		if err != nil {
			t.Fatal(err)
		}
		defer dut.Close()

		if i == 0 {
			if _, err := dut.PushDutConfigFile("./configs/remove_route-map_tor1-1.txt"); err != nil {
				t.Fatal(err)
			}
		} else if i == 1 {
			if _, err := dut.PushDutConfigFile("./configs/remove_route-map_tor1-2.txt"); err != nil {
				t.Fatal(err)
			}
		} else {
			break
		}
	}
}
