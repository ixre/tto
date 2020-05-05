#!lang:ts＃!name:列表界面
#!target:ts/feature/{{.table.Prefix}}/{{.table.Name}}/list.vue
<template>
  <div class="app-container mod-grid-container">
    <!-- 查询和其他操作 -->
    <div class="filter-container mod-grid-bar">
        <div class="mod-grid-bar-create">
            <el-button v-permission="['admin']" class="filter-item" type="primary" icon="el-icon-edit" @click="handleCreate">新增</el-button>
        </div>
        <div class="mod-grid-bar-filter">
            <el-form :inline="true">
                <el-form-item label="关键词：">
                    <el-input v-model="listQuery.keyword" clearable class="filter-item" style="width: 200px;" placeholder="输入关键词" />
                </el-form-item>
                <el-form-item label="状态：">
                    <el-select v-model="listQuery.state" clearable style="width: 200px" class="filter-item" placeholder="选择状态" >
                        <el-option v-for="(item,index) in stateOptions" :key="index" :label="item.key" :value="item.value" />
                    </el-select>
                </el-form-item>
                <el-form-item>
                    <el-button v-permission="['admin']" class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">查找</el-button>
                </el-form-item>
            </el-form>
        </div>
        <div class="mod-grid-bar-del">
            <el-button v-show="selectedRows.length > 0" v-permission="['admin']" class="filter-item" type="danger" icon="el-icon-edit" @click="batchDelete">删除</el-button>
        </div>
    </div>
    <!-- 表单数据　-->
    <el-table ref="table" v-loading="listLoading" :data="list" :height="tableHeight" border fit highlight-current-row
      element-loading-text="正在查询中..." @selection-change="handleSelectionChange"
      class="mod-grid-table">
      <el-table-column align="center" type="selection" width="55" />
{{range $i,$c := .columns}}
  {{if ends_with $c.Name "_time"}}
      <el-table-column width="100px" align="center" label="{{$c.Comment}}">
        <template slot-scope="scope">
          <span>{{`{{ scope.row.`}}{{$c.Name}}{{ `| parseTime}}`}}</span>
        </template>
      </el-table-column>
  {{else}}
      <el-table-column  width="80" align="center" label="{{$c.Comment}}" prop="{{$c.Name}}"/>
  {{end}}
{{end}}

      <el-table-column align="center" label="操作" width="120">
        <template slot-scope="scope">
            <el-button type="primary" size="mini" icon="el-icon-edit" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button type="danger" size="mini" icon="el-icon-delete" @click="handleDelete(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页　-->
    <pagination class="mod-grid-pagination" v-show="total>0" :total="total" :page.sync="listQuery.page"
      :limit.sync="listQuery.rows" @pagination="getList"/>
  </div>
</template>
{{$Class := .table.Title}}
<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Pagination from '@/components/Pagination/index.vue';
import { I{{$Class}}Dbo,getPaging{{$Class}},delete{{$Class}},batchDelete{{$Class}} } from './api';

@Component({
  name: '{{$Class}}List',
  components: {
    Pagination
  }
})
export default class extends Vue {
  private tableHeight = 0;
  private total = 0;
  private list: I{{$Class}}Dbo[] = [];
  private listLoading = true;
  private stateOptions = [{key:"全部",value:-1},{key:"正常",value:1},{key:"停用",value:0}];
  private selectedIds :any[]= [];
  private selectedRows:any[] = [];
  private listQuery = {
      page: 1,
      rows: 20,
      keyword:"",
      state:0,
      params: "",
  };
  created() {
    this.getList();
  }

  mounted(){
    // 监听窗口大小变化
    this.$nextTick(()=>{
        const offset = ((this.$refs.table as Vue).$el as HTMLDivElement).offsetTop + 125;
        this.tableHeight = window.innerHeight - offset;
        window.onresize =()=> this.tableHeight = window.innerHeight - offset;
    });
  }

  private async getList() {
    this.listLoading = true;
    const { data } = await getPaging{{$Class}}(this.listQuery)
    this.list = data.rows;
    this.total = data.total;
    this.listLoading = false;
  }
  private handleFilter() {
    this.listQuery.page = 1
    this.getList();
  }
  private handleSelectionChange(rows:any[]) {
    this.selectedRows = rows;
    this.selectedIds = [];
    rows.map(it=> this.selectedIds.push(it.id));
  }
  private handleCreate() {
    this.$router.push({path:"../{{.table.Name}}/create"});
  }
  private handleEdit(row:any){
     this.$router.push({path:"../{{.table.Name}}/edit/"+row.{{.table.Pk}}});
  }
  private handleDelete(row:any){
     this.$confirm('执行此操作数据无法恢复,是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
     }).then(async () => {
        let ret = await deleteProdVersion(row.{{.table.Pk}});
        const {errCode,errMsg} = ret.data;
         if(errCode === 0){
          this.$notify.success({
            title: '操作成功',
            message: '操作成功'
          });
          this.getList();
        }else{
          this.$notify.error({
            title: '操作失败',
            message: errMsg
          });
        }
     }).catch(() => {
       return false
     });
  }
  private batchDelete(){
    this.$confirm('执行此操作数据无法恢复,是否继续?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
       }).then(async () => {
          let ret = await batchDeleteProdVersion(this.selectedIds);
          const {errCode,errMsg} = ret.data;
           if(errCode === 0){
            this.$notify.success({
              title: '操作成功',
              message: '操作成功'
            });
            this.getList();
          }else{
            this.$notify.error({
              title: '删除失败',
              message: errMsg
            });
          }
       }).catch(() => {
         return false
       });
  }
}
</script>

<style lang="scss" scoped>
.edit-input {
  padding-right: 100px;
}

.cancel-btn {
  position: absolute;
  right: 15px;
  top: 10px;
}
</style>
