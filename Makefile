test:
	cd ./config; go test -v -count=1 ./... -gcflags=-l;
