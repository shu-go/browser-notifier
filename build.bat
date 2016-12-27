@pushd cmd\browser-notifier
go-bindata-assetfs assets/...
@popd
go build -o browser-notifier.exe -ldflags "-w -s" %* ./cmd/browser-notifier
