default: main datainput cachetest

main: main.go datastore.go redis.go cache.go
	go build $^
datainput: datainput.go
	go build datainput.go
cachetest: cachetest.go datastore.go redis.go cache.go 
	go build $^

clean:
	rm -f main datainput cachetest
