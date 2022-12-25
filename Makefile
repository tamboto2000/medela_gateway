build:
	go build -o medela-gateway ./cmd

run:
	./medela-gateway -c config.json