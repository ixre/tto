package tto

import "github.com/ixre/tto/golang"

var (
	GoEntityRepTemplate     = NewTemplate(golang.TPL_ENTITY_REP, "")
	GoEntityRepIfceTemplate = NewTemplate(golang.TPL_ENTITY_REP_INTERFACE, "")
	GoRepoFactoryTemplate   = NewTemplate(golang.TPL_REPO_FACTORY, "")
)
