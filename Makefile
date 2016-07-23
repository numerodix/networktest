all: build

build: bin/havenet


bin/havenet: src/*.go
	# CGO_ENABLED=0 to enable a static build
	CGO_ENABLED=0 go build -o bin/havenet `ls src/*.go | grep -v '_test.go'`

run: bin/havenet
	bin/havenet
	bin/havenet -6

test:
	make clean
	make run

unittest:
	(cd src && go test -v -cover)


all-archs:
	rm dist/*
	# Darwin
	GOOS=darwin GOARCH=386 make clean build
	mv bin/havenet dist/havenet-darwin32
	GOOS=darwin GOARCH=amd64 make clean build
	mv bin/havenet dist/havenet-darwin64
	# Dragonfly
	GOOS=dragonfly GOARCH=amd64 make clean build
	mv bin/havenet dist/havenet-dragonfly64
	# Freebsd
	GOOS=freebsd GOARCH=386 make clean build
	mv bin/havenet dist/havenet-freebsd32
	GOOS=freebsd GOARCH=amd64 make clean build
	mv bin/havenet dist/havenet-freebsd64
	# Linux
	GOOS=linux GOARCH=386 make clean build
	mv bin/havenet dist/havenet-linux32
	GOOS=linux GOARCH=amd64 make clean build
	mv bin/havenet dist/havenet-linux64
	# Openbsd
	GOOS=openbsd GOARCH=386 make clean build
	mv bin/havenet dist/havenet-openbsd32
	GOOS=openbsd GOARCH=amd64 make clean build
	mv bin/havenet dist/havenet-openbsd64
	# Netbsd
	GOOS=netbsd GOARCH=386 make clean build
	mv bin/havenet dist/havenet-netbsd32
	GOOS=netbsd GOARCH=amd64 make clean build
	mv bin/havenet dist/havenet-netbsd64
	# Windows
	GOOS=windows GOARCH=386 make clean build
	mv bin/havenet dist/havenet-win32.exe
	GOOS=windows GOARCH=amd64 make clean build
	mv bin/havenet dist/havenet-win64.exe


clean:
	-@rm bin/*
