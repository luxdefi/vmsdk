// Copyright (C) 2023-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

// "token-cli" implements tokenvm client operation interface.
package main

import (
	"os"

	"github.com/luxdefi/vmsdk/examples/tokenvm/cmd/token-cli/cmd"
	"github.com/luxdefi/vmsdk/utils"
)

func main() {
	if err := cmd.Execute(); err != nil {
		utils.Outf("{{red}}token-cli exited with error:{{/}} %+v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
