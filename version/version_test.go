package version_test

import (
	"fmt"
	"github.com/playmood/restful-demo/version"
	"testing"
)

func TestVersion(t *testing.T) {
	fmt.Println(version.FullVersion())
}
