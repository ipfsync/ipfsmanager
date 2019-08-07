module github.com/ipfsync/ipfsmanager

go 1.12

require (
	github.com/ipfs/go-ipfs v0.4.22
	github.com/ipfs/go-ipfs-config v0.0.3
	github.com/ipfs/interface-go-ipfs-core v0.0.8
	github.com/ipfsync/common v0.0.0
	google.golang.org/appengine v1.4.0 // indirect
)

replace github.com/ipfsync/appserver => ../appserver

replace github.com/ipfsync/ipfsync => ../ipfsync

replace github.com/ipfsync/ipfsmanager => ../ipfsmanager

replace github.com/ipfsync/common => ../common
