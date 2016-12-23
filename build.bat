@pushd cmd
go-bindata-assetfs assets/...
go build -o ../browser-notifier.exe -ldflags "-w -s" %*
@popd
