// +build !deploy_build

package data

import (
	"net/http"
)

var Assets http.FileSystem = http.Dir("data/assets")
