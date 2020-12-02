package lang

import "testing"

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : lang_test.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-12-02 13:44
 * description :
 * history :
 */

func TestPkgStyleLikeGo(t *testing.T) {
	pkg := "github/com/ixre/go2o/core"
	//pkg = PkgStyleLikeGo(pkg)
	pkg = pkgRegex.ReplaceAllString(pkg,".$1/")
	t.Log(pkg)
}