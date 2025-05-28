build/client:
	GOOS=js GOARCH=wasm go build -o ./internal/assets/main.wasm ./cmd/client/main.go

build/server:
	go build -o ./tmp/server ./cmd/server/main.go

live/client:
	air --build.cmd="make build/client" \
		--build.bin=true \
		--build.delay=100 \
		--misc.clean_on_exit=true

live/server:
	air --build.cmd="make build/server" \
		--build.bin="./tmp/server" \
		--build.send_interrupt=true \
		--build.include_ext=html,css,js,wasm,svg \
		--build.include_dir="internal" \
		--build.delay=100 \
		--misc.clean_on_exit=true

.PHONY: build/client build/server live/client live/server