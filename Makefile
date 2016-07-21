all: bin/havenet


test:
	make clean
	make run

unittest:
	(cd src && go test -v -cover)

run: bin/havenet
	bin/havenet
	bin/havenet -6

bin/havenet: src/*.go
	# CGO_ENABLED=0 to enable a static build
	CGO_ENABLED=0 go build -o bin/havenet `ls src/*.go | grep -v '_test.go'`


clean:
	-@rm bin/*
