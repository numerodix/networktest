all: bin/havenet


test:
	make clean
	make run

run: bin/havenet
	bin/havenet

bin/havenet: src/*.go
	go build -o bin/havenet src/*.go


clean:
	-@rm bin/*
