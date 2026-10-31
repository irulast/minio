// Copyright 2026 irulast
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import "testing"

func TestCredentialRedactionPreservesStoredConfiguration(t *testing.T) {
	// A legacy configuration can outlive the help schema that originally wrote it.
	// Redacted diagnostics must remain safe without mutating the recovery export.
	const subsystem = "legacy-notification-fixture"
	cfg := Config{subsystem: map[string]KVS{Default: {
		{Key: "password", Value: "fixture-password"},
		{Key: "connection_string", Value: "fixture-dsn"},
		{Key: "MINIO_NOTIFY_NATS_USER_CREDENTIALS", Value: "/fixture/credentials"},
		{Key: "format", Value: "namespace"},
	}}}
	redacted := cfg.RedactSensitiveInfo()
	for _, key := range []string{"password", "connection_string", "MINIO_NOTIFY_NATS_USER_CREDENTIALS"} {
		if got := redacted[subsystem][Default].Get(key); got != "*redacted*" {
			t.Fatalf("credential field %s not redacted", key)
		}
		if cfg[subsystem][Default].Get(key) == "*redacted*" {
			t.Fatal("redaction mutated stored configuration")
		}
	}
	if redacted[subsystem][Default].Get("format") != "namespace" {
		t.Fatal("ordinary configuration was changed")
	}
	if got := cfg.getTargetKVS(subsystem, Default, true); len(got) != 1 || got[0].Key != "format" {
		t.Fatal("redacted config response retained credentials")
	}
	if got := cfg.getTargetKVS(subsystem, Default, false); len(got) != 4 {
		t.Fatal("authorized recovery export lost configuration")
	}
}
