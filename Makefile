build:
	go run github.com/goreleaser/goreleaser@latest build --snapshot --rm-dist --skip-validate
	mv dist/auto-devops_linux_amd64/auto-devops auto-devops
	docker build --pull . -t paskalmaksim/auto-devops:dev

push:
	docker push paskalmaksim/auto-devops:dev

run:
	rm -rf /tmp/auto-devops
	mkdir -p /tmp/auto-devops
	go run ./cmd -log.level=debug
	ls -la /tmp/auto-devops

install:
	goreleaser build --single-target --snapshot --rm-dist --skip-validate
	sudo mv dist/auto-devops_darwin_amd64/auto-devops /usr/local/bin/auto-devops

test:
	./scripts/validate-license.sh
	go mod tidy
	go fmt ./cmd/... ./pkg/...
	go vet ./cmd/... ./pkg/...
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run -v
	AUTO_DEVOPS_BOOTSTRAP=testdata/auto-devops.zip go test -race -coverprofile coverage.out ./cmd/... ./pkg/...

cover:
	go tool cover -html=coverage.out

zip:
	rm -rf auto-devops.zip
	cd testdata; zip -r ../auto-devops.zip .
	unzip -l auto-devops.zip