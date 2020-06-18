#!target:../../src/resource/{{.table.Name}}_res.py
{{$title := .table.Title}}
{{$pkName := .table.Pk}}\
{{$comment := .table.Comment}}\
{{$path := name_path .table.Name}}\
{{$pkType := type "py" .table.PkType}}\
from flask import request
from flask_restful import Resource, abort
from ..core import report_query
from ..util import dump_query
from ..model.{{.table.Name}} import {{$title}}Model
from ..service.{{.table.Name}}_service import {{$title}}Service
from ..util.response import ok, error

_s = {{$title}}Service()


class {{$title}}(Resource):
    s = _s

    # 获取{{$comment}}
    def get(self, id):
        e = self.s.get({{$pkType}}(id))
        if e is None:
            abort(404)
        return e.to_dict()

    # 删除{{$comment}}
    def delete(self, id):
        r = self.s.delete(int(id))
        if r[0] < 0:
            return error("删除失败"+r[1])
        else:
            return ok()

    # 更新{{$comment}}
    def put(self, id):
        if not request.is_json:
            return error("error data")
        d = {{$title}}Model().from_dict(request.get_json())
        d.{{$pkName}} = {{$pkType}}(id)
        r = self.s.save(d)
        if r[0] <= 0:
            return error("更新失败" + r[1])
        return ok()


class {{$title}}List(Resource):
    s = _s

    # 返回{{$comment}}列表
    def get(self):
        return self.s.query_list("")

    # 新增{{$comment}}
    def post(self):
        if not request.is_json:
            return error("error data")
        d = {{$title}}Model().from_dict(request.get_json())
        r = self.s.save(d)
        if r[0] <= 0:
            return error("新增失败" + r[1])
        return ok()

    # 批量删除{{$comment}}
    def delete(self):
        return error("警告：处于安全考虑,未启用批量操作")
        r = self.s.batch_delete(request.get_json())
        if r[0] < 0:
            return error("删除失败"+r[1])
        return ok()


class {{$title}}Paging(Resource):
    items = report_query.items

    # 返回{{$comment}}列表
    def get(self):
        query = request.args
        p = dump_query.parse_params(query.get("params"))
        p["page_size"] = query.get("rows")
        p["page_index"] = query.get("page")
        count, data = self.items.get("{{.table.Prefix}}/{{$title}}List").dumps_data(p)
        return {"total": count, "rows": data}


# please make sure register api in main.py. code likes:
# {{.table.Name}}_res.reg_api(api)
def route(api):
    api.add_resource({{$title}}List, "{{.global.url_prefix}}/{{$path}}")
    api.add_resource({{$title}}, "{{.global.url_prefix}}/{{$path}}/<id>")
    api.add_resource({{$title}}Paging, "{{.global.url_prefix}}/{{$path}}/paging")