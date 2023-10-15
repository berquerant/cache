GOMOD = go mod
GOTEST = go test -cover -race

.PHONY: test
test:
	$(GOTEST) ./...

.PHONY: init
init:
	$(GOMOD) tidy

.PHONY: vuln
vuln:
	go run golang.org/x/vuln/cmd/govulncheck ./...

.PHONY: vet
vet:
	go vet ./...
