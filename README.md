# Build EAD

- Install and configure golang https://golang.org/dl/
- Install the following packates by typing the following: 
  ```
  go get github.com/logrusorgru/aurora
  go get github.com/gabriel-vasile/mimetype
  go get github.com/hoisie/mustache
  ```
- Download this project with `go get github.com/antonkrug/ead`
- Go into the project:
-   On Windows: `cd %GOPATH%/src/github.com/antonkrug/ead`
-   On Linux: `cd $GOPATH/src/github.com/antonkrug/ead`

- Generate vfsdata `go generate`
- Build the final binary `go build`
