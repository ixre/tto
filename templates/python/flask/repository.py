#!target:../../src/repo/{{.table.Name}}_repo.py
{{$title := .table.Title}}{{$pkName := .table.Pk}}\
{{$comment := .table.Comment}}\
{{$Model := str_join "" .table.Title "Model"}}
{{$pkType := type "py" .table.PkType}}
{{$pkTypeId := $.table.PkType}}
from sqlalchemy.orm import sessionmaker

from ..core.orm import raw_session, session_query, session_maker
from ..model.{{.table.Name}} import {{$title}}Model


# {{$comment}}仓储
class {{$title}}Repo:
    sess: sessionmaker

    def __init__(self, session):
        self.sess = session

    # 创建{{$comment}}仓储
    @staticmethod
    def create():
        return {{$title}}Repo(raw_session)

    # 返回所有数据
    def all(self):
        with session_query(self.sess) as db:
            return db.query({{$Model}}).all()

    # 获取单个{{$comment}}
    def get(self, pk: {{$pkType}})-> {{$Model}}:
        with session_query(self.sess) as db:
            return db.query({{$title}}Model).get(pk)

    # 查找单个{{$comment}}
    def find_one(self,*criterion)-> {{$Model}}:
        with session_query(self.sess) as db:
            return db.query({{$Model}}).filter(criterion).one()

    # 筛选{{$comment}}
    def filter(self,*criterion):
        with session_query(self.sess) as db:
            return db.query({{$Model}}).filter(criterion)

    # 根据条件筛选{{$comment}}
    def filter_by(self,**kwargs):
        with session_query(self.sess) as db:
            return db.query({{$Model}}).filter_by(kwargs)

    # 删除{{$comment}}
    def delete(self, pk: {{$pkType}}) -> int:
        with session_maker(self.sess) as db:
            e = db.query({{$Model}}).filter({{$Model}}.id == pk)
            if e is None:
                return -1
            e.delete()
            return 1

    # 批量删除{{$comment}}
    def batch_delete(self, arr: list) -> int:
        with session_maker(self.sess) as db:
            e = db.query({{$Model}}).filter({{$Model}}.id in arr)
            if e is None:
                return -1
            e.delete()
            return 1

    # 保存{{$comment}}
    def save(self, entity: {{$Model}}) -> {{$pkType}}:
        db = self.sess
        pk = entity.{{$pkName}}
        {{if equal_any $pkTypeId 3 4 5}}\
        if pk <= 0:
            entity.{{$pkName}} = None
        {{else if eq $pkTypeId 1}}\
        if pk == "":
            entity.{{$pkName}} = None
        {{end}}\
        try:
            if entity.{{$pkName}} is not None:
                self.__update__(db, pk, entity)
            else:
                db.add(entity)
                db.flush()
                pk = entity.{{$pkName}}
            db.commit()
        except():
            db.rollback()
            raise
        finally:
            db.close()
        return pk

    # 更新数据{{$comment}}
    def __update__(self, s, pk, entity):
        dst = s.query({{$title}}Model).filter({{$title}}Model.{{$pkName}} == pk).first()
        {{range $i,$c := exclude .columns $pkName}}\
        dst.{{$c.Name}} = entity.{{$c.Name}}
        {{end}}

