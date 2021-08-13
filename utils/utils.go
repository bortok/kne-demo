package utils

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	pb "google.golang.org/protobuf/encoding/protojson"
	"io/ioutil"

	"keysight.com/otgclient"

	log "github.com/sirupsen/logrus"

	"github.com/golang/protobuf/proto"
	"github.com/openconfig/gnmi/client"

	gclient "github.com/openconfig/gnmi/client/gnmi"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
	"os"
	"strings"
)

type PortInfo struct {
	Tx uint64
	Rx uint64
}

var (
	Q                 = client.Query{TLS: &tls.Config{}}
	Context           = context.Background()
	GnmiServer        = flag.String("gnmi", "gnmi-service.default.svc.cluster.local:50051", "URI to gnmi server")
	Config            = flag.String("config", "", "config json file")
	GrpcServer        = flag.String("grpc", "grpc-service.default.svc.cluster.local:40051", "grpc server")
	PortMap           = make(map[string]*PortInfo)
	AteFlag           = &StringList{}
	AristaSshServices = &StringList{}
)

func init() {
	flag.Var(AteFlag, "ixia_loc", "Locators for 2 Ixia nodes.[DP1+CP1,DP2+CP2]")

	AteFlag.Set("service-athena-traffic-engine1.athena-dataplane.svc.cluster.local:5555+" +
		"service-athena-traffic-engine1.athena-dataplane.svc.cluster.local:50071," +
		"service-athena-traffic-engine2.athena-dataplane.svc.cluster.local:5555+" +
		"service-athena-traffic-engine2.athena-dataplane.svc.cluster.local:50071")
}

type StringList struct {
	Values []string
}

func (s StringList) String() string {
	return strings.Join(s.Values, ",")
}

func (s *StringList) Set(v string) error {
	s.Values = nil
	vals := strings.Split(v, ",")
	for _, elem := range vals {
		s.Values = append(s.Values, strings.TrimSpace(elem))
	}
	return nil
}

func ExecuteCapabilities(Context context.Context) error {
	Q.Addrs = []string{*GnmiServer}
	Q.Timeout = 30 * time.Second
	Q.Credentials = nil
	//Q.TLS.InsecureSkipVerify = true
	Q.TLS = nil

	r := &gpb.CapabilityRequest{}
	if err := proto.UnmarshalText("", r); err != nil {
		return fmt.Errorf("unable to parse gnmi.CapabilityRequest from %q : %v", "", err)
	}
	log.Infof("try get client...%v", Q)
	c, err := gclient.New(Context, client.Destination{
		Addrs:       Q.Addrs,
		Target:      Q.Target,
		Timeout:     Q.Timeout,
		Credentials: Q.Credentials,
		TLS:         Q.TLS,
	})
	log.Infof("check client status...")
	if err != nil {
		return fmt.Errorf("could not create a gNMI client: %v", err)
	}
	log.Infof("invoke capabilities on client...")
	response, err := c.(*gclient.Client).Capabilities(Context, r)
	if err != nil {
		return fmt.Errorf("target returned RPC error for Capabilities(%q) : %v", r.String(), err)
	}
	log.Infof("display capabilities...%s", response)
	return nil
}

func ExecutePoll(Context context.Context) error {
	Q.Addrs = []string{*GnmiServer}
	Q.Timeout = 10 * time.Second
	Q.Credentials = nil
	//Q.TLS.InsecureSkipVerify = true
	Q.TLS = nil
	Q.Type = client.Once

	Q.ProtoHandler = func(msg proto.Message) error {
		v := msg.(*gpb.SubscribeResponse)

		notification := v.GetUpdate()
		updates := notification.GetUpdate()
		for _, update := range updates {
			var data interface{}
			err := json.Unmarshal(update.Val.GetJsonVal(), &data)
			if err != nil {
				return nil
			}

			dict, _ := data.(map[string]interface{})
			portName, _ := dict["name"].(string)
			portTx, _ := dict["frames_tx"].(float64)
			portRx, _ := dict["frames_rx"].(float64)
			if _, ok := PortMap[portName]; ok {
				PortMap[portName].Tx = uint64(portTx)
				PortMap[portName].Rx = uint64(portRx)
			}
		}
		return nil
	}

	c := &client.BaseClient{}
	clientTypes := []string{}
	if err := c.Subscribe(Context, Q, clientTypes...); err != nil {
		return fmt.Errorf("client had error while displaying results:    %v", err)
	}
	return nil
}

func ReadStream(filename string) otgclient.Config {
	var config otgclient.Config

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("failed reading data from %s: %s", filename, err)
		return config
	}
	err = pb.Unmarshal([]byte(data), &config)
	if err != nil {
		jsonErr, _ := err.(*json.SyntaxError)
		problemPart := data[jsonErr.Offset-5 : jsonErr.Offset+5]
		log.Errorf("%v ~ error near '%s' (offset %d)", err, problemPart, jsonErr.Offset)
		log.Errorf("Unmarshal error: %v", err)
		os.Exit(1)
	}
	return config
}

func ValidateInput(no_ixia_nodes int) bool {

	inputErrors := false

	if *Config == "" {
		log.Errorf("Config json must be supplied.")
		inputErrors = true
	}
	if len(AteFlag.Values) != 2 {
		log.Errorf("Locators for %d ixia nodes were provided. Expected 2 locators.", no_ixia_nodes)
		inputErrors = true
	}

	return inputErrors

}

func GetAristaMacErrorMessage(node string, intf string) string {
	message := fmt.Sprintf("Could not fetch MAC for %s for %s node.\n", intf, node)
	message += fmt.Sprintf("Either script getAristaMac.sh is missing or 'no password' ssh config for admin not set on Arista.")
	return message
}

func LogAteFlag() {
	log.Infof("    Ixia locators for Athena controller:")
	no_ixia_nodes := len(AteFlag.Values)
	index := 0
	for index < no_ixia_nodes {
		vals := strings.Split(AteFlag.Values[index], "+")
		log.Infof("	Port%v:  DP: %s", index+1, vals[0])
		log.Infof("        	CP: %s", vals[1])
		index += 1
	}
}

func LogTestInput() {
	log.Infof("Test will be run with following parameters:")
	log.Infof("    Config JSON file: %s", *Config)
	log.Infof("    GNMI Server url: %s", *GnmiServer)
	log.Infof("    GRPC Server url: %s", *GrpcServer)

	LogAteFlag()
}

func HasErrorInSetConfigResponse(resp *otgclient.SetConfigResponse, err error) bool {
	if err != nil {
		log.Errorf("SetConfig returned error: %v", err)
		return true
	}

	configErr_500 := resp.GetStatusCode_500()
	if configErr_500 != nil {
		log.Errorf("SetConfig returned error: %v", configErr_500.InternalServerError.ResponseError.Errors)
		return true
	}

	configErr_400 := resp.GetStatusCode_400()
	if configErr_400 != nil {
		log.Errorf("SetConfig returned error: %v", configErr_400.BadRequest.ResponseError.Errors)
		return true
	}

	configErr_200 := resp.GetStatusCode_200()
	if len(configErr_200.Success.ResponseWarning.Warnings) > 0 {
		log.Infof("SetConfig returned warning: %v", configErr_200.Success.ResponseWarning.Warnings)
	}
	return false
}

func HasErrorInSetTransmitResponse(resp *otgclient.SetTransmitStateResponse, err error) bool {
	if err != nil {
		log.Errorf("Failed to start traffic: %v", err)
		return true
	}

	configErr_400 := resp.GetStatusCode_400()
	if configErr_400 != nil {
		log.Errorf("SetTransmitState returned error: %v", configErr_400.BadRequest.ResponseError.Errors)
		return true
	}

	configErr_200 := resp.GetStatusCode_200()
	if len(configErr_200.Success.ResponseWarning.Warnings) > 0 {
		log.Infof("SetTransmitState returned warning: %v", configErr_200.Success.ResponseWarning.Warnings)
	}
	return false
}
