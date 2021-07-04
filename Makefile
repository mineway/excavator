ALL_CMD= $(sort $(dir $(wildcard ./cmd/*/)))

# Make test
test:
ifdef bin
	@go test ./cmd/$(bin)
else
	@for cmd in ${ALL_CMD}; do \
		echo "Go Test on" $$cmd ":" ; \
		go test $$cmd ; \
	done
endif

# Make build on dist folder
build:
ifdef bin
	@go build -o ./dist/${bin} -v ./cmd/$(bin)
else
	@for cmd in ${ALL_CMD}; do \
		echo "Go Build on" $$cmd ":" ; \
        go build -o ./dist/$$(basename $$cmd) -v $$cmd ; \
	done
endif

# Run a command
run:
ifdef bin
	@go run -v ./cmd/$(bin)
else
	@echo "bin argument is required"
endif

# Clear one or multiple binaries
clear:
ifdef bin
	@rm -f ./dist/${bin}
else
	@for cmd in ${ALL_CMD}; do \
		@rm -f ./dist/$$(basename $$cmd)
	done
endif

# Run a command
worker-dev:
	@go run -v ./cmd/worker build
	@go run -v ./cmd/worker