package athena

//package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"keysight.com/otgclient"
	"keysight.com/utils"
	"testing"
	"time"
)

func TestPacketForwardBgpv4(t *testing.T) {
	*utils.Config = "ixia_2node_config_bgp4_attributes.json"
	packetForward(t)
}

//func TestPacketForwardBgpv6(t *testing.T) {
//	*utils.Config = "config_bgp6_attributes.json"
//	packetForward(t)
//}

func packetForward(t *testing.T) {

	flag.Parse()

	utils.LogTestInput()

	inputErrors := utils.ValidateInput(2)
	if inputErrors == true {
		t.Errorf("Aborting test run due to input errors.")
		return
	}

	var err error

	configuration := utils.ReadStream(*utils.Config)
	*(configuration.Ports)[0].Location = utils.AteFlag.Values[0]
	*(configuration.Ports)[1].Location = utils.AteFlag.Values[1]

	url := *utils.GrpcServer
	log.Infof("Connecting to grpc server.. %s ", url)
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	log.Infof("Connected to grpc server.. %s ", url)
	if err != nil {
		t.Errorf("Failed to connect to grpc server on %s", *utils.GrpcServer)
		return
	}
	defer conn.Close()

	openApiClient := otgclient.NewOpenapiClient(conn)
	if openApiClient == nil {
		t.Errorf("Failed to get open api client")
		return
	}

	configResp, err := openApiClient.SetConfig(utils.Context, &otgclient.SetConfigRequest{Config: &configuration})
	if utils.HasErrorInSetConfigResponse(configResp, err) {
		t.Errorf("Failed to SetConfig. Aborting test.")
		return
	}

	tx_port_names := []string{"p1", "p2"}
	flow_names := []string{"p1-p2", "p2-p1"}

	for i := 0; i < len(tx_port_names); i++ {
		utils.Q.Queries = append(utils.Q.Queries, []string{"port_metrics[name=" + tx_port_names[i] + "]"})
		utils.PortMap[tx_port_names[i]] = &utils.PortInfo{Tx: 0, Rx: 0}
	}
	for i := 0; i < len(flow_names); i++ {
		utils.Q.Queries = append(utils.Q.Queries, []string{"flow_metrics[name=" + flow_names[i] + "]"})
		utils.PortMap[flow_names[i]] = &utils.PortInfo{Tx: 0, Rx: 0}
	}

	log.Infof("Ixia ports have been configured.")
	log.Infof("Waiting for 20 secs for sessions to be configured and started")
	time.Sleep(20 * time.Second)

	state := otgclient.TransmitState{FlowNames: nil, State: otgclient.TransmitState_State_start}
	transResp, err := openApiClient.SetTransmitState(utils.Context, &otgclient.SetTransmitStateRequest{TransmitState: &state})
	if utils.HasErrorInSetTransmitResponse(transResp, err) {
		t.Errorf("Could not start traffic.Aborting test.")
		return
	}
	log.Infof("PacketForward config is set and traffic started.")
	log.Infof("Lets verify stats for the 2 bi-directional flows.")

	expected_rx_packets := 1000

	iter := 0
	for iter < 10 {
		err = utils.ExecutePoll(utils.Context)
		if err != nil {
			t.Errorf("Failed to get ExecutePoll: %v", err)
			return
		}

		if utils.PortMap[flow_names[0]].Rx == uint64(expected_rx_packets) && utils.PortMap[flow_names[1]].Rx == uint64(expected_rx_packets) {
			log.Infof("Expected Rx packets recieved by reciever ports for all flows - final stats:")
			for k, v := range utils.PortMap {
				log.Infof("%s: %v", k, v)
			}

			configResp1, err1 := openApiClient.SetConfig(utils.Context, &otgclient.SetConfigRequest{})
			if utils.HasErrorInSetConfigResponse(configResp1, err1) {
				return
			}
			return
		}
		for k, v := range utils.PortMap {
			log.Infof("[%v] %s: %v", iter, k, v)
		}
		if utils.PortMap[flow_names[0]].Rx > uint64(expected_rx_packets) || utils.PortMap[flow_names[1]].Rx > uint64(expected_rx_packets) {
			break
		}

		iter += 1
	}

	t.Errorf("Failed to get correct stats - final stats:")
	for k, v := range utils.PortMap {
		t.Errorf("%s: %v", k, v)
	}

	configResp1, err1 := openApiClient.SetConfig(utils.Context, &otgclient.SetConfigRequest{})
	if utils.HasErrorInSetConfigResponse(configResp1, err1) {
		return
	}
}
