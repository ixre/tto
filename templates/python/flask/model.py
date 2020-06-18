#!target:../../src/model/{{.table.Name}}.py
from sqlalchemy import String, Column, BigInteger, Integer
from src.core.orm import Base
{{$pkName := .table.Pk}}

# {{.table.Comment}}
class {{.table.Title}}Model(Base):
    __tablename__ = '{{.table.Name}}'

    def __init__(self):
        return self

    # 创建新的{{.table.Comment}}
    {{$columns := exclude .columns "create_time" "update_time"}}\
    def __init__(self, {{range $i,$c := $columns}}\
    {{$c.Name}}={{default "py" $c.Type}}{{if not (is_last $i $columns)}}, {{end}}\
    {{end}}):
        {{range $i,$c := $columns}}\
        self.{{$c.Name}} = {{$c.Name}}
        {{end}}\

{{range $i,$c := .columns}}\
    {{$sqlType := sql_type "py" $c.Type $c.Length}}
    # {{$c.Comment}}\
{{if $c.IsPk}}
    {{$c.Name}} = Column("{{$c.Name}}", {{$sqlType}}, primary_key=True\
    {{if $c.IsAuto}}, autoincrement=True{{end}}, comment="{{$c.Comment}}")\
{{else}}
    {{$c.Name}} = Column("{{$c.Name}}", {{$sqlType}}, comment="{{$c.Comment}}")\
{{end}}{{end}}

    # 转为dict
    def to_dict(self)->dict:
        return {
            {{$columns := .columns}} \
            {{range $i,$c := $columns}}\
            "{{$c.Name}}": self.{{$c.Name}}{{if not (is_last $i $columns)}},{{else}}
        }{{end}}
            {{end}}

    # 转为dict
    def from_dict(self,d: dict):
        {{range $i,$c := $columns}}\
        self.{{$c.Name}} = d["{{$c.Name}}"]
        {{end}}\
        return self

    # 拷贝对象数据
    @staticmethod
    def copy(src):
        dst = {{.table.Title}}Model()
        {{range $i,$c := .columns}}\
        dst.{{$c.Name}} = src.{{$c.Name}}
        {{end}}\
        return dst

