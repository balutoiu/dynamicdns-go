OUT_DIR = _output
export OUT_DIR

.PHONY: all build

all build:
	go build -o ${OUT_DIR}/dynamicdns-go

run:
	@./${OUT_DIR}/dynamicdns-go

docker:
	docker build -t alinbalutoiu/dynamicdns-go .

clean:
	rm -rf ${OUT_DIR}
