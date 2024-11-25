[private]
default:
  @just --list --unsorted

# Source
source:
  @go run cmd/data/main.go

# Run
run *macs:
  @go run cmd/yeah/main.go {{macs}}

# Serve
serve:
  @go run cmd/yeah/main.go -l -v debug

# Deploy
deploy:
  flyctl deploy --remote-only
