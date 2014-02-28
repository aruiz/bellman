default: main

tools: datainput

tests: cachetest

main: main.go datastore.go redis.go cache.go config.go
	go build $^
datainput: datainput.go
	go build datainput.go
cachetest: cachetest.go datastore.go redis.go cache.go  config.go
	go build $^
clean:
	rm -f main datainput cachetest
