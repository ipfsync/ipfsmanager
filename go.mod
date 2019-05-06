module github.com/ipfsync/ipfsmanager

go 1.12

require (
	github.com/ipfs/go-ipfs v0.4.20
	github.com/ipfs/go-ipfs-config v0.0.1
	github.com/ipfs/interface-go-ipfs-core v0.0.6
)

replace github.com/ipfsync/appserver => ../appserver

replace github.com/ipfsync/ipfsync => ../ipfsync

replace github.com/ipfsync/ipfsmanager => ../ipfsmanager
