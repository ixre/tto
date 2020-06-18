#!target:../../src/repo/{{.table.Name}}_repo.py.gen
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
    def delete(self, pk: {{$pkType}}):
        with session_maker(self.sess) as db:
            db.delete(pk)

    # 保存{{$comment}}
    def save(self, entity: {{$Model}}) -> {{$pkType}}:
        db = self.sess
        pk = entity.{{$pkName}}
        {{if eq $pkTypeId 5}}\
        if pk <= 0:
            entity.{{$pkName}} = None
        {{else if eq $pkTypeId 24}}\
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
    def __update__(self,s,pk,entity):
        dst = s.query({{$title}}Model).filter({{$title}}Model.{{$pkName}} == pk).first()
        {{range $i,$c := exclude .columns $pkName}}\
        dst.{{$c.Name}} = entity.{{$c.Name}}
        {{end}}
