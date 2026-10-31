// Copyright 2026 irulast
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"context"
	"strings"
	"testing"
)

func TestMaintenanceTopologyValidation(t *testing.T) {
	for _, pools := range []int{0, 2, 3} {
		obj, err := newErasureServerPools(context.Background(), make(EndpointServerPools, pools))
		if obj != nil || err == nil || !strings.Contains(err.Error(), "requires exactly one erasure pool") {
			t.Fatalf("pool count %d: expected startup rejection before storage initialization, got object=%v error=%v", pools, obj, err)
		}
	}
}
