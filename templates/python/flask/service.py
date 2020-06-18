#!target:../../src/service/{{.table.Name}}_service.py
{{$title := .table.Title}}{{$pkName := .table.Pk}}\
{{$comment := .table.Comment}}\
from ..model.{{.table.Name}} import {{$title}}Model
from ..repo.{{.table.Name}}_repo import {{$title}}Repo
{{$pkType := type "py" .table.PkType}} \
{{$pkTypeId := $.table.PkType}} \


# {{$comment}}服务
class {{$title}}Service:
    repo = {{$title}}Repo.create()

    def __init__(self):
        pass

    # 获取{{$comment}}
    def get(self, id):
        return self.repo.get(id)

    # 保存{{$comment}}
    def save(self, e) -> ({{$pkType}}, str):
        dst = None
        {{if equal_any $pkTypeId 3 4 5}}\
        if e.{{$pkName}} > 0: \
        {{else if eq $pkTypeId 1}}\
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
        return self.repo.save(dst), ""

    # 查询{{$comment}}
    def query_list(self, param):
        return []

    # 删除{{$comment}}
    def delete(self, id) -> (int, str):
        # 处理删除逻辑,如果有错误返回字符
        return self.repo.delete(id), ""

    # 批量删除{{$comment}}
    def batch_delete(self, arr: list) -> (int, str):
        i = self.repo.batch_delete(arr)
        return i, ""