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

		// Apply to POD switches only (index 0 and 2)
		if i == 0 || i == 1 {
			if _, err := dut.PushDutConfigFile("./configs/apply_route-map_pod.txt"); err != nil {
				t.Fatal(err)
			}
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

		// Apply to POD switches only (index 0 and 2)
		if i == 0 || i == 1 {
			if _, err := dut.PushDutConfigFile("./configs/remove_route-map_pod.txt"); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestClosPodHosts_FixRouteMap(t *testing.T) {
	for i, location := range opts.DutPorts() {
		dut, err := api.NewSshClient(opts, location, "admin")
		if err != nil {
			t.Fatal(err)
		}
		defer dut.Close()

		// Apply to POD switches only (index 0 and 2)
		if i == 0 || i == 1 {
			if _, err := dut.PushDutConfigFile("./configs/apply_route-map_pod_fix.txt"); err != nil {
				t.Fatal(err)
			}
		}
	}
}
