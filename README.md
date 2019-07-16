# Build EAD

- Install and configure golang https://golang.org/dl/
- Install the following packates by typing the following: 
  ```
  go get github.com/logrusorgru/aurora
  go get github.com/gabriel-vasile/mimetype
  go get github.com/hoisie/mustache
  go get github.com/shurcooL/vfsgen
  go get github.com/dustin/go-humanize
  ```
- Download this project with `go get github.com/antonkrug/ead` ("undefined: Assets" error is expected as the `go generate` was not run yet)
- Go into the project:
  - On Windows: `cd %GOPATH%/src/github.com/antonkrug/ead`
  - On Linux: `cd $GOPATH/src/github.com/antonkrug/ead`

- Generate vfsdata `go generate`
- Build the project
  - Build the final native binary `go build`
  - To build all other platforms run `bash ./build_all_platforms.sh` (Tested on Linux and on Windows under GitBash command line)
