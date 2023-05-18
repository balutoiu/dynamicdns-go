OUT_DIR = _output
export OUT_DIR

.PHONY: all build

all build:
	CGO_ENABLED=0 go build -o ${OUT_DIR}/dynamicdns-go cmd/dynamicdns-go/main.go

run:
	@./${OUT_DIR}/dynamicdns-go

clean:
	rm -rf ${OUT_DIR}
