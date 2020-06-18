#!target:../../src/service/{{.table.Name}}_service.py.gen
{{$title := .table.Title}}{{$pkName := .table.Pk}}\
{{$comment := .table.Comment}}\
from time import time
from ..model.{{.table.Name}} import {{$title}}Model
from ..repo.{{.table.Name}}_repo import {{$title}}Repo
{{$pkTypeId := $.table.PkType}}


# {{$comment}}服务
class {{$title}}Service:
    repo = {{$title}}Repo.create()

    def __init__(self):
        pass

    # 获取{{$comment}}
    def get(self, id):
        return self.repo.get(id)

    # 删除{{$comment}}
    def delete(self, id):
        return self.repo.delete(id)

    # 保存{{$comment}}
    def save(self, e):
        dst = None
        {{if eq $pkTypeId 5}}\
        if e.{{$pkName}} > 0: \
        {{else if eq $pkTypeId 24}}\
        if e.{{$pkName}} is not None and e.{{$pkName}} != "":\
        {{else}}\
        if e.{{$pkName}} is not None: \
        {{end}}
            dst = self.repo.get(e.{{$pkName}})
        else:
            dst = {{$title}}Model()
            {{$c := try_get .columns "create_time"}}{{if $c }}dst.create_time = time(){{end}}
        # update other fields \
        {{range $i,$c := exclude .columns $pkName "create_time" "update_time"}}
        dst.{{$c.Name}} = e.{{$c.Name}}{{end}}\
        {{$c := try_get .columns "update_time"}}\
        {{if $c}}dst.updateTime = time(){{end}}
        return self.repo.save(dst)

    # 查询{{$comment}}
    def query_list(self, param):
        return []

