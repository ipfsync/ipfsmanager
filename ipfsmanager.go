package ipfsmanager

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/ipfsync/common"

	"github.com/ipfs/go-ipfs/repo"

	config "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/namesys"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	iface "github.com/ipfs/interface-go-ipfs-core"
)

var (
	// ErrIpfsDaemonLocked returns if another ipfs daemon is running and locking fsrepo.
	ErrIpfsDaemonLocked = errors.New("another IPFS daemon is running")
)

type IpfsManager struct {
	nd        *core.IpfsNode
	repo      repo.Repo
	API       iface.CoreAPI
	ndStarted bool
}

// NewIpfsManager creates a new IpfsManager. It will initialize IPFS if it's not initialized.
func NewIpfsManager(repoRoot string) (*IpfsManager, error) {
	daemonLocked, err := fsrepo.LockedByOtherProcess(repoRoot)
	if err != nil {
		return nil, err
	}
	if daemonLocked {
		return nil, ErrIpfsDaemonLocked
	}

	if err := common.CheckWritable(repoRoot); err != nil {
		return nil, err
	}

	_, err = loadPlugins(repoRoot)
	if err != nil {
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

	r, err := fsrepo.Open(repoRoot)
	if err != nil {
		return nil, err
	}

	return &IpfsManager{repo: r}, nil
}

func (im *IpfsManager) StartNode() error {

	ctx := context.Background()

	cfg := &core.BuildCfg{
		Repo:   im.repo,
		Online: true,
	}

	var err error
	im.nd, err = core.NewNode(ctx, cfg)
	if err != nil {
		return err
	}

	im.API, err = coreapi.NewCoreAPI(im.nd)
	if err != nil {
		return err
	}

	return nil
}

func (im *IpfsManager) StopNode() error {
	if err := im.nd.Close(); err != nil {
		return err
	}

	im.nd = nil
	im.API = nil

	return nil
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

func loadPlugins(repoPath string) (*loader.PluginLoader, error) {
	pluginpath := filepath.Join(repoPath, "plugins")

	// check if repo is accessible before loading plugins
	var plugins *loader.PluginLoader
	ok, err := common.CheckPermissions(repoPath)
	if err != nil {
		return nil, err
	}
	if !ok {
		pluginpath = ""
	}
	plugins, err = loader.NewPluginLoader(pluginpath)
	if err != nil {
		return nil, err
	}
	if err := plugins.Initialize(); err != nil {
		return nil, err
	}

	if err := plugins.Inject(); err != nil {
		return nil, err
	}

	return plugins, nil
}
