// Copyright (C) 2023-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package gossiper

import (
	"context"

	"github.com/luxdefi/node/ids"
	"github.com/luxdefi/node/snow/engine/common"
)

type Gossiper interface {
	Run(common.AppSender)
	TriggerGossip(context.Context) error // may be triggered by run already
	HandleAppGossip(ctx context.Context, nodeID ids.NodeID, msg []byte) error
	Done() // wait after stop
}
