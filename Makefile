.PHONY: test build install-hooks

test:
go test -v ./...

build:
go build -o bin/hyperagent main.go

install-hooks:
cp scripts/pre-push .git/hooks/pre-push
chmod +x .git/hooks/pre-push
echo "Git hooks installed successfully."
