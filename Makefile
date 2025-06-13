.PHONY: clean
clean:
	rm -rf bin/

.PHONY: build
build:
	go build -o bin/fluidsim

.PHONY: build-wasm
build-wasm:
	GOOS=js GOARCH=wasm go build -o bin/fluidsim.wasm

.PHONY: decomp
decomp: build-wasm
	@echo decompile wasm
	wasm2wat bin/fluidsim.wasm > bin/fluidsim.wat
	cp bin/fluidsim.wat bin/copyfluidsim.wat

.PHONY: recomp
recomp:
	@echo recompile wat
	wat2wasm bin/copyfluidsim.wat -o bin/copyfluidsim.wasm

.PHONY: recompcopy
recompcopy: recomp
	@echo copy recompile output
	cp bin/copyfluidsim.wasm docs/assets/wasm/fluidsim.wasm

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
