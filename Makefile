run:
	go run cmd/sso/main.go --config=./config/local.yaml

lint:
	golangci-lint --config golangci.yaml run ./... --deadline=2m --timeout=2m