package config

import (
	"cluster-membership-go/src/common"
	"cluster-membership-go/src/model"
	"fmt"
	"strconv"
	"strings"
)

//Configuration type
type Configuration struct {
	Mode  string
	Nodes []model.Node
}

//NilConfig for return when needed
var NilConfig = Configuration{}

//NodeID is the command line argument name for seeds node id
var NodeID = "id"

//NodeAddress is the command line argument name for seeds node address
var NodeAddress = "address"

//NodeProtocolPort is the command line argument name for seeds node protocol port
var NodeProtocolPort = "protocol.port"

//NodeServerPort is the command line argument name for seeds node server port
var NodeServerPort = "server.port"

//AppMode is the command line argument that tells us if it is executing in test or production mode
var AppMode = "mode"
var runningMode = "RUNNING"
var testMode = "TEST"

//ArgListSep is the string used to parse the command line parameters that contains list of elements
var ArgListSep = ","

//ArgStringList implements
type ArgStringList struct {
	Values []string
}

//IsBoolFlag implements the Interface Value, from package flag
func (arg *ArgStringList) IsBoolFlag() bool {
	return false
}

//Set implements the Interface Value, from package flag
func (arg *ArgStringList) Set(argS string) error {
	var array = strings.Split(argS, ",")
	if len(array) == 0 {
		return &common.BasicError{ErrorS: "Field is not valid, no value provided"}
	}
	arg.Values = array
	return nil
}

//String implements the Interface Value, from package flag
func (arg *ArgStringList) String() string {
	return common.ArrayToString(arg.Values, ArgListSep)
}

//Validate validates the command line input
func Validate(mode string, ids ArgStringList, addresses ArgStringList,
	protocolPorts ArgStringList, serverPorts ArgStringList) (Configuration, error) {
	if mode != strings.ToUpper(runningMode) && mode != strings.ToUpper(testMode) {
		return Configuration{}, &common.ArgumentParsingError{ErrorS: fmt.Sprintf("The 'mode' variable must be either %s or %s", runningMode, testMode)}
	}
	if len(ids.Values) != len(addresses.Values) || len(addresses.Values) != len(protocolPorts.Values) ||
		len(protocolPorts.Values) != len(serverPorts.Values) {
		return NilConfig, &common.ArgumentParsingError{ErrorS: fmt.Sprintf("Arguments %s, %s, %s and %s must contains the same number of element", NodeID, NodeAddress, NodeProtocolPort, NodeServerPort)}
	}
	if chk := common.CheckNonEmpty(ids.Values, NodeID); chk != nil {
		return NilConfig, chk
	} else if chk := common.CheckNonEmpty(addresses.Values, NodeAddress); chk != nil {
		return NilConfig, chk
	} else if chk := common.CheckNonEmpty(protocolPorts.Values, NodeProtocolPort); chk != nil {
		return NilConfig, chk
	} else if chk := common.CheckNonEmpty(serverPorts.Values, NodeServerPort); chk != nil {
		return NilConfig, chk
	} else {
		return buildConf(mode, ids, addresses, protocolPorts, serverPorts)
	}

}

func buildConf(mode string, ids ArgStringList, addresses ArgStringList,
	protocolPorts ArgStringList, serverPorts ArgStringList) (Configuration, error) {
	var nodes = make([]model.Node, len(ids.Values))
	for i := 0; i < len(ids.Values); i++ {
		var pPort, pError = strconv.Atoi(protocolPorts.Values[i])
		if pError != nil {
			return Configuration{}, pError
		}
		var sPort, sError = strconv.Atoi(serverPorts.Values[i])
		if pError != nil {
			return Configuration{}, sError
		}
		nodes[i] = model.Node{ID: ids.Values[i], Address: addresses.Values[i], ProtocolPort: pPort, ServerPort: sPort}
	}
	return Configuration{Nodes: nodes, Mode: mode}, nil
}
