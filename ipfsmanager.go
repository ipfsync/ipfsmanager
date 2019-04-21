package ipfsmanager

import (
	"context"
	"errors"
	"fmt"
	config "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/namesys"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"os"
	"path"
)

var (
	// ErrIpfsDaemonLocked returns if another ipfs daemon is running and locking fsrepo.
	ErrIpfsDaemonLocked = errors.New("another IPFS daemon is running")
)

type IpfsManaer struct {
	nd *core.IpfsNode
}

// NewIpfsManager creates a new IpfsManager. It will initialize IPFS if it's not initialized.
func NewIpfsManager(repoRoot string) (*IpfsManaer, error) {
	daemonLocked, err := fsrepo.LockedByOtherProcess(repoRoot)
	if err != nil {
		return nil, err
	}
	if daemonLocked {
		return nil, ErrIpfsDaemonLocked
	}

	if err := checkWritable(repoRoot); err != nil {
		return nil, err
	}

	if !fsrepo.IsInitialized(repoRoot) {
		conf, err := config.Init(os.Stdout, 2048)
		if err != nil {
			return nil, err
		}
		if err := fsrepo.Init(repoRoot, conf); err != nil {
			return nil, err
		}
		if err := initializeIpnsKeyspace(repoRoot); err != nil {
			return nil, err
		}
	}

}

func checkWritable(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		// dir exists, make sure we can write to it
		testfile := path.Join(dir, "test")
		fi, err := os.Create(testfile)
		if err != nil {
			if os.IsPermission(err) {
				return fmt.Errorf("%s is not writeable by the current user", dir)
			}
			return fmt.Errorf("unexpected error while checking writeablility of repo root: %s", err)
		}
		_ = fi.Close()
		return os.Remove(testfile)
	}

	if os.IsNotExist(err) {
		// dir doesn't exist, check that we can create it
		return os.Mkdir(dir, 0775)
	}

	if os.IsPermission(err) {
		return fmt.Errorf("cannot write to %s, incorrect permissions", err)
	}

	return err
}

func initializeIpnsKeyspace(repoRoot string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := fsrepo.Open(repoRoot)
	if err != nil { // NB: repo is owned by the node
		return err
	}

	nd, err := core.NewNode(ctx, &core.BuildCfg{Repo: r})
	if err != nil {
		return err
	}
	defer nd.Close()

	return namesys.InitializeKeyspace(ctx, nd.Namesys, nd.Pinning, nd.PrivateKey)
}
