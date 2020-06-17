#!lang:ts＃!name:全功能界面
#!lang:ts#!target:ts/feature/{{name_path .table.Name}}/all.vue
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
                    <el-input v-model="queryParams.keyword" clearable class="filter-item" style="width: 200px;" placeholder="输入关键词" />
                </el-form-item>
                <el-form-item label="状态：">
                    <el-select v-model="queryParams.state" clearable style="width: 200px" class="filter-item" @change="fetchData" >
                        <el-option v-for="(value,attr) in states" :key="value" :label="attr" :value="value" />
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
    <el-table ref="table" v-loading="list.loading" :data="list.data" :height="tableHeight" border fit highlight-current-row
              @selection-change="handleSelectionChange" class="mod-grid-table">
        <el-table-column align="center" type="selection" width="55" />
        {{range $i,$c := .columns}}
        {{if ends_with $c.Name "_time"}}
        <el-table-column width="160" align="center" label="{{$c.Comment}}">
            <template slot-scope="scope">
                <span>{{`{{ scope.row.`}}{{$c.Name}}{{ `| parseTime}}`}}</span>
            </template>
        </el-table-column>
        {{else if ends_with $c.Name "state"}}
        <el-table-column width="120" align="center" label="{{$c.Comment}}">
            <template slot-scope="scope">
                <span v-for="(value,attr) in states" v-if="value === scope.row.state">{{`{{attr}}`}}</span>
            </template>
        </el-table-column>
        {{else if starts_with $c.Name "is_"}}
        <el-table-column width="120" align="center" label="{{$c.Comment}}">
             <template slot-scope="scope">
                <span v-for="(it,i) in ['否','是']" v-if="i===scope.row.{{$c.Name}}-1">{{`{{it}}`}}</span>
             </template>
        </el-table-column>
        {{else}}
        <el-table-column width="120" align="center" label="{{$c.Comment}}" prop="{{$c.Name}}"/>
        {{end}}
        {{end}}

        <el-table-column align="center" label="操作">
            <template slot-scope="scope">
                <el-button type="primary" size="mini" icon="el-icon-edit" @click="handleEdit(scope.row)">编辑</el-button>
                <el-button type="danger" size="mini" icon="el-icon-delete" v-loading="requesting" @click="handleDelete(scope.row)">删除</el-button>
            </template>
        </el-table-column>
    </el-table>

    <!-- 分页　-->
    <pagination class="mod-grid-pagination" v-show="list.total>0" :total="list.total" :page.sync="list.page"
                :limit.sync="list.rows" @pagination="fetchData"/>

    <!-- 弹出操作框 -->
    <el-dialog class="mod-dialog" :title="dialog.title" :visible.sync="dialog.open" width="40%" @close="dialogClosed">
        <component v-bind:is="modal" :id="dialog.id" @refresh="refresh"></component>
    </el-dialog>

</div>
</template>
{{$Class := .table.Title}}
{{$PkType := .table.PkType}}
<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Pagination from '@/components/Pagination/index.vue';
import {getPaging{{$Class}},delete{{$Class}},batchDelete{{$Class}} } from './api';
import {{$Class}}Form from './form.vue';


// {{.table.Comment}}数据库映射类
interface IListModel {
    {{range $i,$c := .columns}}{{$c.Name}}:{{type "ts" $c.Type}} // {{$c.Comment}}
    {{end}}
}

@Component({
    name: '{{$Class}}List',
    components: {
        Pagination,
        {{$Class}}Form
    }
})
export default class extends Vue {
    private dialog = {title:"标题",open:false,id:0};
    private modal :any = null;
    private requesting = 0;
    private tableHeight = 0;
    private selectedIds :{{type "ts" $PkType}}[]= [];
    private selectedRows:IListModel[] = [];
    private list:{loading:boolean,total:number,page:number,rows:number,data: IListModel[]} = {loading:false,total:0, page: 1, rows: 10,data:[]};
    /** 定义查询参数 */
    private queryParams = {
       keyword:"",
       state:-1
    };
    /** 根据实际情况调整 */
    private states = {"全部":-1,"正常":1,"停用":0};

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
        this.list.page = 1
        this.fetchData();
    }

    private handleSelectionChange(rows:IListModel[]) {
        this.selectedRows = rows;
        this.selectedIds = [];
        rows.map(it=> this.selectedIds.push(it.{{.table.Pk}}));
    }

    // 新增数据
    private handleCreate() {
        this.modalForm(null,"新增{{.table.Comment}}");
    }

    // 编辑数据
    private handleEdit(row:IListModel){
        this.modalForm(row,"编辑{{.table.Comment}}");
    }

    // 打开模态表单
    private modalForm(row:IListModel|null,title:string){
      this.dialog.open = true;
      this.dialog.id = row?row.{{.table.Pk}}:0;
      this.modal = {{$Class}}Form;
      this.dialog.title = title;
    }

    // 关闭模态框时重置模态组件
    private dialogClosed(){
        this.modal = null;
    }

    // 参数state为true时,重置模态框并刷新数据,args接受传入的参数
    private refresh(s:{state:number,args:any}){
        if(!s.state)return;
        this.modal = null;
        this.dialog.open = false;
        this.fetchData(s.args);
    }

    private handleDelete(row:any){
        this.$confirm('执行此操作数据无法恢复,是否继续?', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
        }).then(async () => {
            if(this.requesting === 1)return;this.requesting = 1;
            let ret = await delete{{$Class}}(row.{{.table.Pk}}).finally(()=>this.requesting = 0);
            const {errCode,errMsg} = ret.data;
            if(errCode === 0){
                this.$notify.success({
                    title: '操作成功',
                    message: '操作成功',
                    duration:2000
                });
                await this.fetchData();
            }else{
                this.$notify.error({
                    title: '操作失败',
                    message: errMsg,
                    duration:2000
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
            if(this.requesting === 1)return;this.requesting = 1;
            let ret = await batchDelete{{$Class}}(this.selectedIds).finally(()=>this.requesting = 0);
            const {errCode,errMsg} = ret.data;
            if(errCode === 0){
                this.$notify.success({
                    title: '操作成功',
                    message: '操作成功',
                    duration:2000
                });
                await this.fetchData();
            }else{
                this.$notify.error({
                    title: '删除失败',
                    message: errMsg,
                    duration:2000
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
