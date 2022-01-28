// No Shutdown Eth1 on DUTs

package tests

import (
	"keysight/athena/tests/pkg/api"
	"log"
	"testing"
)

func TestClosPodHosts_NoShutEth1(t *testing.T) {
	log.Println("Next is to bring TOR links to POD1-1 back up and check if taffic would still flow...")

	for _, location := range opts.DutPorts() {
		dut, err := api.NewSshClient(opts, location, "admin")
		if err != nil {
			t.Fatal(err)
		}
		defer dut.Close()

		if _, err := dut.PushDutConfigFile("./configs/noshut_Eth1.txt"); err != nil {
			t.Fatal(err)
		}
	}
}
