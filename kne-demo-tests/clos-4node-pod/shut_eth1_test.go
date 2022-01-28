// Shutdown Eth1 on DUTs

package tests

import (
	"keysight/athena/tests/pkg/api"
	"testing"
)

func TestClosPodHosts_ShutEth1(t *testing.T) {
	for _, location := range opts.DutPorts() {
		dut, err := api.NewSshClient(opts, location, "admin")
		if err != nil {
			t.Fatal(err)
		}
		defer dut.Close()

		if _, err := dut.PushDutConfigFile("./configs/shut_Eth1.txt"); err != nil {
			t.Fatal(err)
		}
	}
}
