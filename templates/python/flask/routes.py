#!kind:1
#!target:../../src/resource/routes.py.gen
{{range $i,$table := .tables}}\
from . import {{$table.Name}}_res
{{end}}

def map_route(api):
    {{range $i,$table := .tables}}\
    {{$table.Name}}_res.reg_api(api)
    {{end}}