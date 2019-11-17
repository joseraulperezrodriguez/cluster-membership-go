package main

import (
	"cluster-membership-go/src/config"
	"flag"
)

func main() {

	var mode = flag.String(config.AppMode, "", "Possible values are: RUNNING or TEST")
	var ids = config.ArgStringList{}
	var addresses = config.ArgStringList{}
	var protocolPorts = config.ArgStringList{}
	var serverPorts = config.ArgStringList{}
	flag.Var(&ids, config.NodeID, "The node id list, separated by comma")
	flag.Var(&addresses, config.NodeAddress, "The node address list, separated by comma")
	flag.Var(&protocolPorts, config.NodeProtocolPort, "The node protocol port list, separated by comma")
	flag.Var(&serverPorts, config.NodeServerPort, "The node server port list, separated by comma")

	flag.Parse()

	config.Validate(*mode, ids, addresses, protocolPorts, serverPorts)

}
