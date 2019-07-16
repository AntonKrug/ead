echo "Setting Bash strict mode"
# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -e
set -u
set -o pipefail
echo

echo "Go version and enviroment"
go version
go env
echo

echo "Supported distributions:"
go tool dist list
echo

echo "Installing depedencies"
go get github.com/logrusorgru/aurora
go get github.com/gabriel-vasile/mimetype
go get github.com/hoisie/mustache
go get github.com/shurcooL/vfsgen
go get github.com/dustin/go-humanize

echo "Generating  assets"
go generate

echo "Creating clean build folder"
mkdir -p build
rm -r ./build
mkdir -p build

TS=`date +"%Y%m%d-%H%M%S"`

echo "Building Windows"
GOOS=windows GOARCH=amd64 go build -o build/ead-windows-x86-64-$TS.exe
GOOS=windows GOARCH=386 go build -o build/ead-windows-x86-32-$TS.exe
GOOS=windows GOARCH=arm go build -o build/ead-windows-arm-$TS.exe

echo "Building Linux"
GOOS=linux GOARCH=amd64 go build -o build/ead-linux-x86-64-$TS
GOOS=linux GOARCH=386 go build -o build/ead-linux-x86-32-$TS
GOOS=linux GOARCH=arm64 go build -o build/ead-linux-arm-64-$TS
GOOS=linux GOARCH=arm go build -o build/ead-linux-arm-32-$TS

echo "Building FreeBSD"
GOOS=freebsd GOARCH=amd64 go build -o build/ead-freebsd-x86-64-$TS
GOOS=freebsd GOARCH=386 go build -o build/ead-freebsd-x86-32-$TS

echo "Building MacOS"
GOOS=darwin GOARCH=amd64 go build -o build/ead-macos-x86-64-$TS
GOOS=darwin GOARCH=386 go build -o build/ead-macos-x86-32-$TS

echo
echo "Finished artifacts:"
ls -la ./build
