// Shutdown Eth1 on TORs

package tests

import (
	"keysight/athena/tests/pkg/api"
	"testing"
)

func TestClosPodHosts_ShutEth1(t *testing.T) {
	for i, location := range opts.DutPorts() {
		dut, err := api.NewSshClient(opts, location, "admin")
		if err != nil {
			t.Fatal(err)
		}
		defer dut.Close()

		// Apply to TOR switches only (index 2 and 3)
		if i == 2 || i == 3 {
			if _, err := dut.PushDutConfigFile("./configs/shut_Eth1.txt"); err != nil {
				t.Fatal(err)
			}
		}
	}
}
