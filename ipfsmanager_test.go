package ipfsmanager

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

var testdataDir = filepath.Join(".", "testdata")
var ipfsPath = filepath.Join(testdataDir, "ipfs")

func TestMain(m *testing.M) {
	// Ensure testdata dir exists
	_ = os.MkdirAll(testdataDir, os.ModePerm)
	// Remove old testing datastore
	_ = os.RemoveAll(ipfsPath)

	retCode := m.Run()

	// Cleanup
	_ = os.RemoveAll(ipfsPath)

	os.Exit(retCode)
}

func TestIpfsManager(t *testing.T) {
	im, err := NewIpfsManager(ipfsPath)
	if err != nil {
		t.Fatalf("Unable to create IpfsManager. Error: %s", err)
	}
	defer im.Close()

	keys, err := im.API.Key().List(context.TODO())

	for _, key := range keys {
		t.Logf("Key ID: %s", key.ID())
	}
}
