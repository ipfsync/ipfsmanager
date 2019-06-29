module github.com/ipfsync/ipfsmanager

go 1.12

require (
	github.com/ipfs/go-ipfs v0.4.21
	github.com/ipfsync/common v0.0.0
	google.golang.org/appengine v1.4.0 // indirect
)

replace github.com/ipfsync/appserver => ../appserver

replace github.com/ipfsync/ipfsync => ../ipfsync

replace github.com/ipfsync/ipfsmanager => ../ipfsmanager

replace github.com/ipfsync/common => ../common
