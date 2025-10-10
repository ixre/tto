package lang

import "testing"

func TestPkgStyleLikeGo(t *testing.T) {
	pkg := "github/com/ixre/go2o/core"
	//pkg = PkgStyleLikeGo(pkg)
	pkg = pkgRegex.ReplaceAllString(pkg, ".$1/")
	t.Log(pkg)
}
