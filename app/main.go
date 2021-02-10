package main

import (
	"cluster-membership-go/config"
	"flag"
	"fmt"
)

func createFlags() (mode *string, configPath *string, ids *config.ArgStringList, addresses *config.ArgStringList, protocolPorts *config.ArgStringList, serverPorts *config.ArgStringList) {
	mode = flag.String(config.AppMode, "", "Possible values are: RUNNING or TEST")
	configPath = flag.String(config.ConfigPath, "", "The path to the config file")
	ids = &config.ArgStringList{}
	addresses = &config.ArgStringList{}
	protocolPorts = &config.ArgStringList{}
	serverPorts = &config.ArgStringList{}
	flag.Var(ids, config.NodeID, "The node id list, separated by comma")
	flag.Var(addresses, config.NodeAddress, "The node address list, separated by comma")
	flag.Var(protocolPorts, config.NodeProtocolPort, "The node protocol port list, separated by comma")
	flag.Var(serverPorts, config.NodeServerPort, "The node server port list, separated by comma")
	return
}

func main() {

	var mode, configPath, ids, addresses, protocolPorts, serverPorts = createFlags()

	flag.Parse()

	configuration, parsingError := config.Validate(*mode, *configPath, *ids, *addresses, *protocolPorts, *serverPorts)
	if parsingError != nil {
		fmt.Println(fmt.Errorf("Error produced at parsing  configuration parameters.\n%s", parsingError.Error()))
		return
	}
	configurationError := config.AddProgramConfig(&configuration)
	if configurationError != nil {
		fmt.Println(fmt.Errorf("Bad configuration format.\n%s", configurationError.Error()))
		return
	}

	//TODO open rest api ports

	//TODO start protocol engine

}
