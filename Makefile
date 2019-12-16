OUT_DIR = _output
export OUT_DIR

.PHONY: all build

all build:
	go build -o ${OUT_DIR}/dynamicdns-go

run:
	@USERNAME=test PASSWORD=test DOMAIN=test.domain.com ./${OUT_DIR}/dynamicdns-go

clean:
	rm -rf ${OUT_DIR}
