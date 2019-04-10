BIN=bin

eurek8s:
	go build -o ${BIN}/eurek8s cmd/eurek8s/main.go

clean:
	[ -d bin ] && rm -r bin/
