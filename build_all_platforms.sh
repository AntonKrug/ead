go tool dist list

mkdir build

echo "Building Windows"
GOOS=windows GOARCH=amd64 go build -o build/ead-windows-x86-64.exe
GOOS=windows GOARCH=386 go build -o build/ead-windows-x86-32.exe
GOOS=windows GOARCH=arm go build -o build/ead-windows-arm.exe

echo "Building Linux"
GOOS=linux GOARCH=amd64 go build -o build/ead-linux-x86-64
GOOS=linux GOARCH=386 go build -o build/ead-linux-x86-32
GOOS=linux GOARCH=arm64 go build -o build/ead-linux-arm-64
GOOS=linux GOARCH=arm go build -o build/ead-linux-arm-32

echo "Building FreeBSD"
GOOS=freebsd GOARCH=amd64 go build -x -o build/ead-freebsd-x86-64
GOOS=freebsd GOARCH=386 go build -x -o build/ead-freebsd-x86-32

echo "Building MacOS"
GOOS=darwin GOARCH=amd64 go build -o build/ead-macos-x86-64
GOOS=darwin GOARCH=386 go build -o build/ead-macos-x86-32

