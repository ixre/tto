package tto

import (
	golang2 "github.com/ixre/tto/lang/golang"
)

var (
	GoEntityRepTemplate     = NewTemplate(golang2.TPL_ENTITY_REP, "", true)
	GoEntityRepIfceTemplate = NewTemplate(golang2.TPL_ENTITY_REP_INTERFACE, "", true)
	GoRepoFactoryTemplate   = NewTemplate(golang2.TPL_REPO_FACTORY, "", true)
)
