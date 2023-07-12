// protfile simply print the protofile use for rpcmanager
package rpcmanager

import _ "embed"

//go:embed rpcmanager.proto
var RPCManagerProtoFile string
