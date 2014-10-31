all: tclient tserver

rebuild: clean all

tclient:
	mkdir -p ./bin
	go build -o ./bin/client client.go common.go marzullo.go
	
tserver:
	mkdir -p ./bin
	go build -o ./bin/server server.go common.go

clean:
	rm -rf ./bin

test:
	go run test_marzullo.go marzullo.go