// protfile simply print the protofile use for rpcmanager

package main

import (
	"fmt"

	"github.com/gfanton/grpcutil/rpcmanager"
)

func main() { fmt.Print(rpcmanager.RPCManagerProtoFile) }
