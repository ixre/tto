package impl

#!type:0#!lang:go
#!target:{{.global.pkg}}/mongo/repo/impl/{{.table.Name}}_repo_impl.go
{{$title := .table.Title}}
{{$pkType := type "go" .table.PkType}}
{{$pkProp := .table.PkProp}}
{{$tableName := .table.Name}}
{{$shortTitle := .table.ShortTitle}}
{{$p := substr .table.Name 0 1 }}
{{$structName := join (lower_title $title) "RepositoryImpl"}}
import(
    "context"
    "{{pkg "go" .global.pkg}}/mongo/repo/model"
    "{{pkg "go" .global.pkg}}/mongo/repo"
    "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// 获取自增主键ID
func getIncrementPkId(d *mongo.Database, collname string) int {
	table := d.Collection("ids") // 指定表名为ids表
	var result struct {
		Name   string `json:"name" bson:"name"`
		NextID int    `json:"next_id" bson:"next_id"`
	}
	// vs调试会缓存上次的结果,导致看上去没有自增
	table.FindOneAndUpdate(
		context.TODO(),
		bson.M{"name": collname},
		bson.M{"$inc": bson.M{"next_id": 1}}).Decode(&result)
	if result.NextID == 0 {
		_, _ = table.InsertOne(context.TODO(), bson.M{
			"name": collname, "next_id": 1})
		return 1
	}
	return result.NextID
}

var _ repo.I{{$title}}Repository = new({{$structName}})
type {{$structName}} struct{
    d *mongo.Database
}


// New{{$title}}Repository Create new {{$title}}Repository
func New{{$title}}Repository(d *mongo.Database) repo.I{{$title}}Repository{
    return &{{$structName}}{
        d:d,
    }
}
// Get{{$shortTitle}} 获取{{.table.Comment}}
func ({{$p}} *{{$structName}}) Get{{$shortTitle}}(primary {{$pkType}})*model.{{$shortTitle}}{
    dst := model.{{$shortTitle}}{}
	ret := {{$p}}.d.Collection("{{$tableName}}").FindOne(context.TODO(), bson.M{"_id": primary})
	ret.Decode(&dst)
	return &dst
}

/*
// Get{{$shortTitle}}By GetBy {{.table.Comment}}
func ({{$p}} *{{$structName}}) Get{{$shortTitle}}By(where string,v ...interface{})*model.{{$shortTitle}}{
    e := model.{{$shortTitle}}{}
    err := {{$p}}.d.GetBy(&e,where,v...)
    if err == nil{
        return &e
    }
    if err != sql.ErrNoRows{
      log.Printf("[ Orm][ Error]: %s; Entity:{{$shortTitle}}\n",err.Error())
    }
    return nil
}
*/

/*
// Count{{$shortTitle}} Count {{.table.Comment}} by condition
func ({{$p}} *{{$structName}}) Count{{$shortTitle}}(where string,v ...interface{})(int,error){
   return {{$p}}.d.Count(model.{{$shortTitle}}{},where,v...)
}
*/

/*
// Select{{$shortTitle}} Select {{.table.Comment}}
func ({{$p}} *{{$structName}}) Select{{$shortTitle}}(where string,v ...interface{})[]*model.{{$shortTitle}} {
    list := make([]*model.{{$shortTitle}},0)
    err := {{$p}}.d.Select(&list,where,v...)
    if err != nil && err != sql.ErrNoRows{
      log.Printf("[ Orm][ Error]: %s; Entity:{{$shortTitle}}\n",err.Error())
    }
    return list
}
*/

// Save{{$shortTitle}} 保存{{.table.Comment}}
func ({{$p}} *{{$structName}}) Save{{$shortTitle}}(v *model.{{$shortTitle}})(err error){
	if v.{{$pkProp}} <= 0 {
        v.{{$pkProp}} = getIncrementPkId({{$p}}.d,"{{$tableName}}")
	    _, err = {{$p}}.d.Collection("{{$tableName}}").InsertOne(context.TODO(), v)
    }else{
	    _, err = {{$p}}.d.Collection("{{$tableName}}").UpdateOne(context.TODO(), bson.M{"_id": v.Id}, bson.M{"$set": v})
    }
    return err
}

// Delete{{$shortTitle}} 删除{{.table.Comment}}
func ({{$p}} *{{$structName}}) Delete{{$shortTitle}}(primary {{$pkType}}) error {
    _, err := {{$p}}.d.Collection("{{$tableName}}").DeleteOne(context.TODO(), bson.M{"_id": primary})
	return err
}

/*
// BatchDelete{{$shortTitle}} Batch Delete {{.table.Comment}}
func ({{$p}} *{{$structName}}) BatchDelete{{$shortTitle}}(where string,v ...interface{})(int64,error) {
    r,err := {{$p}}.d.Delete(model.{{$shortTitle}}{},where,v...)
    if err != nil && err != sql.ErrNoRows{
      log.Printf("[ Orm][ Error]: %s; Entity:{{$shortTitle}}\n",err.Error())
    }
    return r,err
}
*/

// PagingQuery{{$shortTitle}} {{.table.Comment}}分页数据
func ({{$p}} *{{$structName}}) Paging{{$shortTitle}}(filter bson.M, begin int, size int,
	opts ...*options.FindOptions) (int, []*model.{{$shortTitle}}, error) {
	
	opts = append(opts, options.Find().SetSkip(int64(begin)))
	opts = append(opts, options.Find().SetLimit(int64(size)))
	opts = append(opts, options.Find().SetSort(bson.M{"_id": 1}))
	cursor, err := {{$p}}.d.Collection("data_source").Find(context.TODO(),
		filter, opts...)
	if err != nil {
		return 0, nil, err
	}
	var results []*model.{{$shortTitle}}
	err = cursor.All(context.TODO(), &results)
	_ = cursor.Close(context.TODO())
	if err != nil {
		return 0, results, err
	}
	return 0, results, nil
}
