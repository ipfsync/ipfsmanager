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

	err = im.StartNode()
	if err != nil {
		t.Fatalf("Unable to start Ipfs node. Error: %s", err)
	}

	keys, err := im.API.Key().List(context.TODO())
	if err != nil {
		t.Errorf("Unable to get key. Error: %s", err)
	}

	for _, key := range keys {
		t.Logf("Key ID: %s", key.ID())
	}

	err = im.StopNode()
	if err != nil {
		t.Fatalf("Unable to stop Ipfs node. Error: %s", err)
	}

}
