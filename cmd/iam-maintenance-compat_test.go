// Copyright 2026 irulast
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
)

func TestIAMMaintenanceLifecycle(t *testing.T) {
	suite := newTestSuiteIAM(TestSuiteCommon{serverType: "ErasureSD", signer: signerV4}, false)
	c := &check{t, suite.serverType}
	suite.SetUpSuite(c)
	defer suite.TearDownSuite(c)
	suite.TestSTSForRoot(c)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	accessKey, secretKey := mustGenerateCredentials(c)
	if err := suite.adm.AddUser(ctx, accessKey, secretKey); err != nil {
		t.Fatal(err)
	}
	if _, err := suite.adm.AttachPolicy(ctx, madmin.PolicyAssociationReq{Policies: []string{"readwrite"}, User: accessKey}); err != nil {
		t.Fatal(err)
	}
	bucket := getRandomBucketName()
	client := suite.getUserClient(c, accessKey, secretKey, "")
	if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
		t.Fatal(err)
	}
	data := []byte("IAM restart compatibility fixture")
	if _, err := client.PutObject(ctx, bucket, "fixture", bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{}); err != nil {
		t.Fatal(err)
	}
	suite.RestartIAMSuite(c)
	client = suite.getUserClient(c, accessKey, secretKey, "")
	if info, err := client.StatObject(ctx, bucket, "fixture", minio.StatObjectOptions{}); err != nil || info.Size != int64(len(data)) {
		t.Fatalf("ordinary read after IAM restart: size=%d error=%v", info.Size, err)
	}

	if err := suite.adm.RemoveUser(ctx, accessKey); err != nil {
		t.Fatal(err)
	}
	suite.RestartIAMSuite(c)
	users, err := suite.adm.ListUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if _, exists := users[accessKey]; exists {
		t.Fatal("removed user remained in the administrative inventory after restart")
	}
	_, replacementSecret := mustGenerateCredentials(c)
	if err := suite.adm.AddUser(ctx, accessKey, replacementSecret); err != nil {
		t.Fatal(err)
	}
	if _, err := suite.adm.AttachPolicy(ctx, madmin.PolicyAssociationReq{Policies: []string{"readwrite"}, User: accessKey}); err != nil {
		t.Fatal(err)
	}
	client = suite.getUserClient(c, accessKey, replacementSecret, "")
	if _, err := client.StatObject(ctx, bucket, "fixture", minio.StatObjectOptions{}); err != nil {
		t.Fatal(err)
	}
}
