build:
	goreleaser build --snapshot --rm-dist --skip-validate
	mv dist/auto-devops_linux_amd64/auto-devops auto-devops
	docker build --pull . -t paskalmaksim/auto-devops:dev

push:
	docker push paskalmaksim/auto-devops:dev

run:
	rm -rf /tmp/auto-devops
	mkdir -p /tmp/auto-devops
	go run ./cmd -log.level=debug
	ls -la /tmp/auto-devops

test:
	./scripts/validate-license.sh
	go fmt ./cmd/... ./pkg/...
	go mod tidy
	golangci-lint run -v
	AUTO_DEVOPS_BOOTSTRAP=testdata/auto-devops.zip go test -coverprofile coverage.out ./cmd/... ./pkg/...

cover:
	go tool cover -html=coverage.out

zip:
	rm -rf auto-devops.zip
	cd testdata; zip -r ../auto-devops.zip .
	unzip -l auto-devops.zip