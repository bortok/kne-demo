// Stop Protocols for a running test

package tests

import (
	"keysight/athena/tests/pkg/api"
	"testing"
)

func TestClosPodHosts_StopProtocol(t *testing.T) {
	client, err := api.NewClient(opts)
	if err != nil {
		t.Fatal(err)
	}

	client.StopProtocol()
	client.Close()

	if err != nil {
		t.Fatal(err)
	}
}
