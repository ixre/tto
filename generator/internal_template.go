package generator

import "github.com/ixre/tto/generator/internal"

var(
	TPL_ENTITY_REP *CodeTemplate
	TPL_ENTITY_REP_INTERFACE *CodeTemplate
	TPL_REPO_FACTORY *CodeTemplate
)

func resolveRepTag(s string)*CodeTemplate {
	t := NewTemplate(s)
	return t.Replace("<Ptr>", "{{.Ptr}}", -1).
		Replace("<E>", "{{.E}}", -1).
		Replace("<E2>", "{{.E2}}", -1).
		Replace("<R>", "{{.R}}", -1).
		Replace("<R2>", "{{.R2}}", -1).
		Replace("<IsPK>", "{{.IsPK}}", -1)
}

func init() {
	TPL_ENTITY_REP = resolveRepTag(internal.TPL_ENTITY_REP)
	TPL_ENTITY_REP_INTERFACE = resolveRepTag(internal.TPL_ENTITY_REP_INTERFACE)
	TPL_REPO_FACTORY = resolveRepTag(internal.TPL_REPO_FACTORY)
}


