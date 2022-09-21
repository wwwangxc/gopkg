test:
	cd ./config; go test -v -count=1 ./... -gcflags=-l;
	cd ./mysql; go test -v -count=1 ./... -gcflags=-l;
	cd ./orm; go test -v -count=1 ./... -gcflags=-l;
	cd ./redis; go test -v -count=1 ./... -gcflags=-l;
	cd ./etcd; go test -v -count=1 ./... -gcflags=-l;
