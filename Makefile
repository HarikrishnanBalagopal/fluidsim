.PHONY: clean
clean:
	rm -rf bin/

.PHONY: build
build:
	go build -o bin/fluidsim

.PHONY: build-wasm
build-wasm:
	GOOS=js GOARCH=wasm go build -o bin/fluidsim.wasm

.PHONY: build-tiny-wasm
build-tiny-wasm:
	tinygo build -o bin/fluidsim.wasm -target wasm

.PHONY: copy
copy:
	cp bin/fluidsim.wasm docs/assets/wasm/

.PHONY: serve
serve:
	cd docs/ && python3 -m http.server 8080

.PHONY: benchmark
benchmark:
	go test -cpuprofile profiles/cpu.prof -memprofile profiles/mem.prof -bench . ./utils

.PHONY: visualize-benchmark
visualize-benchmark:
	go tool pprof -http=: profiles/cpu.prof
