// Copyright 2026 irulast
// SPDX-License-Identifier: AGPL-3.0-or-later

package handlers

import (
	"net/http"
	"testing"
)

func TestConfiguredSourceAddressCompatibility(t *testing.T) {
	oldPolicy, oldProxies := sourceIPPolicy, trustedProxies
	t.Cleanup(func() { sourceIPPolicy, trustedProxies = oldPolicy, oldProxies })
	for _, tc := range []struct {
		name, proxies, peer, forwarded, expected string
	}{
		{"direct", "none", "192.0.2.20:30000", "", "192.0.2.20"},
		{"one proxy", "192.0.2.10", "192.0.2.10:30000", "198.51.100.20", "198.51.100.20"},
		{"two proxies", "192.0.2.10, 192.0.2.11", "192.0.2.10:30000", "198.51.100.20, 192.0.2.11", "198.51.100.20"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(EnvTrustedProxies, tc.proxies)
			if err := ConfigureSourceIPTrust(); err != nil {
				t.Fatal(err)
			}
			req := &http.Request{RemoteAddr: tc.peer, Header: make(http.Header)}
			if tc.forwarded != "" {
				req.Header.Set(xForwardedFor, tc.forwarded)
			}
			if got := GetSourceIP(req); got != tc.expected {
				t.Fatalf("client address = %s, want %s", got, tc.expected)
			}
		})
	}
}
