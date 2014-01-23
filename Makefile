main: main.go
	go build main.go
datainput: datainput.go
	go build datainput.go
cachetest: cachetest.go datastore.go redis.go cache.go 
	go build $^
