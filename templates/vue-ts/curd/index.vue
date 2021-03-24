#!lang:ts＃!name:全功能界面
#!lang:ts#!target:ts/feature/{{name_path .table.Name}}/index.vue
{{$Class := .table.Title}}
{{$PkType := .table.PkType}}
{{$Pk := .table.Pk}}
<template>
<div class="app-container mod-grid-container">
    <!-- 查询和其他操作 -->
    <div class="filter-container mod-grid-bar">
        <div class="it mod-grid-bar-back" v-show="false">
            <el-button class="filter-item" @click="handleBack">返回</el-button>
        </div>
        <div class="it mod-grid-bar-filter">
          <el-form :inline="true">
            <el-form-item label="状态:" class="filter-item">
              <el-select v-model="queryParams.state" class="filter-select"
                         @change="fetchData">
                <el-option v-for="(it,i) in stateOptions" :key="i" :label="it.key" :value="it.value"/>
              </el-select>
            </el-form-item>
            <el-form-item label="关键词:" class="filter-item">
              <el-input v-model="queryParams.keyword" clearable class="filter-input"
                        placeholder="请输入关键词"/>
            </el-form-item>
            <el-form-item label="排序方式:" class="filter-item">
              <el-select v-model="queryParams.order_by" class="filter-select">
                <el-option v-for="(it,i) in sortOptions" :key="i" :value="it.value" :label="it.key"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item class="filter-item">
              <el-button v-permission="['admin']" class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">搜索</el-button>
            </el-form-item>
          </el-form>
        </div>
        <div class="it mod-grid-bar-create">
            <el-button v-permission="['admin']" class="filter-item" icon="el-icon-plus" @click="handleCreate">新增</el-button>
        </div>
        <div class="it mod-grid-bar-del">
            <el-button v-show="selectedRows.length > 0" v-permission="['admin']" class="filter-item" type="danger" @click="batchDelete">删除</el-button>
        </div>
    </div>
    <!-- 表单数据　-->
    <el-table ref="table" v-loading="list.loading" :data="list.data" :height="tableHeight"
              fit :highlight-current-row="false" row-key="{{$Pk}}"
              @selection-change="handleSelectionChange" class="mod-grid-table">
        <el-table-column align="center" type="selection" width="45" />
        {{range $i,$c := .columns}} \
        {{if ends_with $c.Name "_time"}} \
        <el-table-column width="160" align="left" label="{{$c.Comment}}">
            <template slot-scope="scope">
                <span>{{`{{ scope.row.`}}{{$c.Name}}{{ `| parseTime}}`}}</span>
            </template>
        </el-table-column>
        {{else if ends_with $c.Name "state"}} \
        <el-table-column width="140" align="left" label="{{$c.Comment}}">
            <template slot-scope="scope">
                <span v-for="(value,attr) in states" v-if="value === scope.row.state">{{`{{attr}}`}}</span>
            </template>
        </el-table-column>
        {{else if starts_with $c.Name "is_"}} \
        <el-table-column width="140" align="left" label="{{$c.Comment}}">
             <template slot-scope="scope">
                <span v-if="scope.row.{{$c.Name}}===1" class="green">是</span>
                <span v-else class="red">否</span>
             </template>
        </el-table-column>
        {{else if $c.IsPk }} \
        <el-table-column width="60" align="left" label="{{$c.Comment}}" prop="{{$c.Name}}"/>
        {{else}} \
        <el-table-column align="left" label="{{$c.Comment}}" prop="{{$c.Name}}"/>
        {{end}} \
        {{end}} \

        <el-table-column align="center" width="160" label="操作">
            <template slot-scope="scope">
              <el-button type="text" icon="el-icon-edit-outline" title="编辑" @click="handleEdit(scope.row)">编辑</el-button>
              <el-button type="text" icon="el-icon-delete" title="删除" v-loading="requesting" @click="handleDelete(scope.row)">删除</el-button>
            </template>
        </el-table-column>
    </el-table>

    <!-- 分页　-->
    <pagination class="mod-grid-pagination" v-show="list.total>0" :total="list.total" :page.sync="list.page"
                :limit.sync="list.rows" @pagination="fetchData"/>

    <!-- 弹出操作框 -->
    <el-dialog width="35%" class="mod-dialog" :title="dialog.title" :visible.sync="dialog.open" @close="()=>dialog.modal=null">
        <component v-bind:is="dialog.modal" v-model="dialog.params" @callback="refresh"></component>
    </el-dialog>

</div>
</template>
<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Pagination from '@/components/Pagination/index.vue';
import {getPaging{{$Class}},delete{{$Class}},batchDelete{{$Class}} } from './api';
import {{$Class}}Form from './form.vue';
import {parseResult} from "@/fx";

// {{.table.Comment}}数据映射类
interface ListModel {  {{range $i,$c := .columns}}
    {{$c.Name}}:{{type "ts" $c.Type}} // {{$c.Comment}} \
    {{end}}
}

@Component({
    name: '{{$Class}}Index',
    components: {
        Pagination,
        {{$Class}}Form
    }
})
export default class extends Vue {
    private dialog:{title:string,open:boolean,params:any,modal:any}= {title:"Form",open:false,params:0,modal: null};
    private requesting = 0;
    private tableHeight = 0;
    private selectedIds :{{type "ts" $PkType}}[]= [];
    private selectedRows:ListModel[] = [];
    private list:{loading:boolean,total:number,page:number,rows:number,data: ListModel[]} = {loading:false,total:0, page: 1, rows: 10,data:[]};

    /** 定义排序条件 */
    private sortOptions = [
      {key: "默认排序", value: "{{.table.Pk}} DESC"},
      {key: "按创建时间先后顺序", value: "{{.table.Pk}} ASC"},
    ];

    /** 定义状态条件(自定义条件) */
    private stateOptions = [
      {key: "全部", value: -1},
      {key: "正常", value: 1},
      {key: "停用", value: 0}
    ];

    /** 定义查询参数 */
    private queryParams = {
      keyword: "",
      where: "0=0",
      state: this.stateOptions[0].value,
      order_by: this.sortOptions[0].value,
    };
    
    created() {
        this.fetchData();
    }

    mounted(){
        // 监听窗口大小变化
        this.$nextTick(()=>{
            const offset = ((this.$refs.table as Vue).$el as HTMLDivElement).offsetTop + 125;
            this.tableHeight = window.innerHeight - offset;
            window.onresize =()=> this.tableHeight = window.innerHeight - offset;
        });
    }

    // 读取分页数据
    private async fetchData(args:any|null = null) {
        const { data } = await getPaging{{$Class}}(this.list.page,this.list.rows,this.queryParams)
        this.list.data = data.rows;
        this.list.total = data.total;
        setTimeout(()=>this.list.loading = false,300);
    }

    private handleFilter() {
        this.list.page = 1;
        this.fetchData();
    }

    private handleSelectionChange(rows:ListModel[]) {
        this.selectedRows = rows;
        this.selectedIds = [];
        rows.map(row=> this.selectedIds.push(row.{{.table.Pk}}));
    }

    // 返回
    private handleBack(){
        this.$store.dispatch('delView', this.$route)
        this.$router.back();
    }

    // 新增数据
    private handleCreate() {
        this.modalForm(null,"新增{{.table.Comment}}");
    }

    // 编辑数据
    private handleEdit(row:ListModel){
        /** // 在新的tab页上打开临时页面
        const id = row.{{.table.Pk}}.toString();
        const path = addTempRoute(this.$router,
          `{{name_path .table.Name}}/details$${id}`,
          "{{.table.Comment}}详情:"+row.name,
          {{$Class}}Form);
        this.$router.push({path,query:{id}});
        */
        this.modalForm(row,"编辑{{.table.Comment}}");
    }

    // 打开模态表单
    private modalForm(row:ListModel|null,title:string){
      /** 关闭模态框时需要重置模态组件的内容 */
      this.dialog = {
        open:true,
        title:title,
        params:row?row.{{.table.Pk}}:0,
        modal:{{$Class}}Form
      };
    }

    // 参数state为true时,重置模态框并刷新数据,args接受传入的参数
    private refresh(s:{state:number,args:any}){
        this.dialog.open = false;
        if(s.state)this.fetchData(s.args);
    }

    private async notifyResult({errCode, errMsg}) {
      if (errCode === 0) {
        this.$notify.success({title: '提示', message: errMsg || '操作成功', duration: 2000});
      } else {
        await this.$alert(errMsg || "操作失败", "提示");
      }
    }

    private handleDelete(row:any){
        this.$confirm('执行此操作数据无法恢复,是否继续?', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
        }).then(async () => {
            if(this.requesting === 1)return;this.requesting = 1;
            let ret = await delete{{$Class}}(row.{{.table.Pk}}).finally(()=>this.requesting = 0);
            let result = parseResult(ret.data);
            await this.notifyResult(result);
            if (result.errCode === 0) {
              await this.fetchData();
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
            if(this.requesting === 1)return;this.requesting = 1;
            let ret = await batchDelete{{$Class}}(this.selectedIds).finally(()=>this.requesting = 0);
            let result = parseResult(ret.data);
            await this.notifyResult(result);
            if (result.errCode === 0) {
              await this.fetchData();
            }
        }).catch(() => {
            return false
        });
    }
}
</script>

<style lang="scss" scoped>
/** 页面内的样式定义在这里 */
</style>
