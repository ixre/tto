package tto

import "github.com/ixre/tto/golang"

var (
	GoEntityRepTemplate     = NewTemplate(golang.TPL_ENTITY_REP, "", true)
	GoEntityRepIfceTemplate = NewTemplate(golang.TPL_ENTITY_REP_INTERFACE, "", true)
	GoRepoFactoryTemplate   = NewTemplate(golang.TPL_REPO_FACTORY, "", true)
)
