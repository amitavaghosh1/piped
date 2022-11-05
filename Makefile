run:
	go run -race main.go

bench:
	go test -cpu 1,2,4,8,16 -benchmem -run=^$ -bench . piped/piper/v2