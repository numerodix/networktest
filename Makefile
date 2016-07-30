all: build

build: bin/havenet


bin/havenet: src/*.go
	python build.py  # fill in the version
	# CGO_ENABLED=0 to enable a static build
	CGO_ENABLED=0 go build -o bin/havenet `ls src/*.go | grep -v '_test.go'`

clean:
	-@rm -f bin/*

test-standard: bin/havenet
	bin/havenet
	bin/havenet -6

test-monochrome: bin/havenet
	bin/havenet -nc

test-version: bin/havenet
	bin/havenet -V

unittest: bin/havenet
	(cd src && go test -v -cover)


all-archs:
	-@rm -f dist/*
	# Darwin
	GOOS=darwin GOARCH=386 make clean build
	@mv bin/havenet dist/havenet-darwin32
	GOOS=darwin GOARCH=amd64 make clean build
	@mv bin/havenet dist/havenet-darwin64
	# Freebsd
	GOOS=freebsd GOARCH=386 make clean build
	@mv bin/havenet dist/havenet-freebsd32
	GOOS=freebsd GOARCH=amd64 make clean build
	@mv bin/havenet dist/havenet-freebsd64
	# Linux
	GOOS=linux GOARCH=386 make clean build
	@mv bin/havenet dist/havenet-linux32
	GOOS=linux GOARCH=amd64 make clean build
	@mv bin/havenet dist/havenet-linux64
	# Windows
	GOOS=windows GOARCH=386 make clean build
	@mv bin/havenet dist/havenet-win32.exe
	GOOS=windows GOARCH=amd64 make clean build
	@mv bin/havenet dist/havenet-win64.exe
