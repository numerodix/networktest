all: havenet


run: havenet
	./havenet

havenet: havenet.go
	go build -o havenet havenet.go


clean:
	rm havenet
