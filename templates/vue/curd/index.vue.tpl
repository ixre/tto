#!lang:ts＃!name:全功能界面
#!lang:ts#!target:vue/{{name_path .table.Name}}/index.vue
{{$Class := .table.Title}}
{{$PkType := .table.PkType}}
{{$Pk := .table.Pk}}
<template>
<div class="app-container mod-grid-container">
    <!-- 查询和其他操作 -->
    <div class="filter-container mod-grid-header">
        <div class="it mod-grid-header-back" v-if="false">
            <el-button class="filter-item" @click="handleBack">返回</el-button>
        </div>
        <div class="it mod-grid-header-filter">
          <el-form :inline="true">
            <el-form-item label="状态:" class="filter-item">
              <el-select v-model="queryParams.state" class="filter-select"
                         @change="queryPagingData">
                <el-option v-for="(it,i) in stateOptions" :key="i" :label="it.key" :value="it.value"/>
              </el-select>
            </el-form-item>
            <el-form-item label="关键词:" class="filter-item">
              <el-input v-model="queryParams.keyword" clearable class="filter-input" placeholder="请输入关键词"/>
            </el-form-item>
            <el-form-item label="排序方式:" class="filter-item">
              <el-select v-model="queryParams.order_by" class="filter-select">
                <el-option v-for="(it,i) in sortOptions" :key="i" :value="it.value" :label="it.key"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item class="filter-item">
              <span v-perm="{key: '', roles: ['admin'], visible: true }" @click="handleFilter">
                <el-button type="primary" icon="el-icon-search">搜索</el-button>
              </span>
            </el-form-item>
            <el-form-item class="filter-item">
              <span v-perm="{key: '[权限key]+1', roles: ['admin'], visible: true }" @click="handleCreate">
                <el-button icon="el-icon-plus">新增</el-button>
              </span>
            </el-form-item>
            <el-form-item class="filter-item">
              <span v-perm="{key: '', roles: ['admin'], visible: true }" @click="handleDelete">
                <el-button v-show="selectedRows.length > 0" type="danger" icon="el-icon-delete" :loading="requesting">删除</el-button>
              </span>
            </el-form-item>
          </el-form>
        </div>
    </div>
    <div class="mod-grid-body">
    <!-- 表单数据 -->
    <el-table ref="table" v-loading="list.loading" :data="list.data" class="mod-grid-table"
              border :height="tableHeight" fit :highlight-current-row="false"
              row-key="{{$Pk}}" empty-text="暂无数据" @selection-change="handleSelectionChange">
        <el-table-column align="center" type="selection" width="40" fixed="left"/>
        <el-table-column type="index" width="50" label="序号"></el-table-column>
        {{range $i,$c := .columns}} \
        {{if ends_with $c.Name "_time"}} \
        <el-table-column width="140" align="left" label="{{$c.Comment}}" prop="{{$c.Name}}" :formatter="formatColTime"/>
        {{else if ends_with $c.Name "state"}} \
        <el-table-column width="140" align="left" label="{{$c.Comment}}">
            <template slot-scope="{row}">
              <span v-for="(it,i) in stateOptions" v-if="it.value === row.state">{{"{{ it.key }}"}}</span>
            </template>
        </el-table-column>
        {{else if starts_with $c.Name "is_"}} \
        <el-table-column width="140" align="left" label="{{$c.Comment}}">
             <template slot-scope="{row}">
                <span v-if="row.{{$c.Name}}===1" class="g-txt-green">是</span>
                <span v-else class="g-txt-red">否</span>
             </template>
        </el-table-column>
        {{else if $c.IsPk }} \
        <el-table-column width="60" align="left" label="{{$c.Comment}}" prop="{{$c.Name}}"/>
        {{else}} \
        <el-table-column align="left" label="{{$c.Comment}}" prop="{{$c.Name}}"/>
        {{end}} \
        {{end}} \

        <el-table-column align="center" width="160" label="操作" fixed="right">
            <template slot-scope="scope">
              <span v-perm="{key: '[权限key]+4', roles: ['admin'], visible: true }" @click="handleEdit(scope.row)">
                <el-button type="primary" size="mini" icon="el-icon-edit-outline">编辑</el-button>
              </span>
              <span v-perm="{key: '[权限key]+2', roles: ['admin'], visible: true }" @click="handleDelete(scope.$index,scope.row)">
                <el-button type="danger" size="mini" icon="el-icon-delete" :loading="requesting && $index === scope.$index">删除</el-button>
              </span>
            </template>
        </el-table-column>
    </el-table>
    </div>
    <div class="mod-grid-footer">
      <!-- 分页 -->
      <el-pagination v-show="list.total > 0" :total="list.total" :page-size.sync="list.rows"
        :current-page.sync="list.page" :page-sizes="[10, 20, 50, 100]"
        @current-change="queryPagingData" @size-change="(v)=>list.rows = v"
        background layout="prev, pager, next,sizes,jumper,total"></el-pagination>
    </div>
    <!-- 弹出操作框 -->
    <el-dialog v-if="dialog.modal != null" width="75%" class="mod-dialog" :title="dialog.title" visible :close-on-click-modal="false" @close="()=>dialog.modal=null">
        <component v-bind:is="dialog.modal" v-model="dialog.params" @callback="refresh"></component>
    </el-dialog>
</div>
</template>
<script lang="ts" setup>
import {onMounted, reactive, ref, nextTick} from "vue";
import {Paging{{$Class}},queryPaging{{$Class}},delete{{$Class}},batchDelete{{$Class}} } from './api';
import {{$Class}}Modal from './modal.vue';
import {Message,MessageBox,router,parseResult,formatColTime} from "@/utils";


const list = reactive<{ loading: boolean, total: number, page: number, rows: number, data: Array<PagingMmLevel> }>({loading:false,total:0, page: 1, rows: 20,data:[]});
const dialog = reactive({title:"Form",params:0,modal: null});
const $index = ref(-1);
const requesting = ref(false);
const selectedIds = ref([]);
const selectedRows = ref([]);
const table = ref(null);
const tableHeight = ref(0);

/** 定义排序条件 */
const sortOptions = [
  {key: "默认排序", value: "{{.table.Pk}} DESC"},
  {key: "按创建时间先后顺序", value: "{{.table.Pk}} ASC"},
];

/** 定义状态条件(自定义条件) */
const stateOptions = [
  {key: "全部", value: -1},
  {key: "正常", value: 1},
  {key: "停用", value: 0}
];

/** 定义查询参数 */
const queryParams = reactive({
  keyword: "",
  where: "0=0",
  state: stateOptions[0].value,
  order_by: sortOptions[0].value,
});

onMounted(()=>{
    queryPagingData();
    nextTick(()=>{
      const rec = table.value.$el.getBoundingClientRect();
      tableHeight.value = document.body.clientHeight - rec.top - 80;
    })
})

// 读取分页数据
const queryPagingData = async ()=> {
  if(list.loading) return;
  list.loading=true;
    const { data } = await queryPaging{{$Class}}(list.page,list.rows,queryParams)
      .finally(()=>list.loading=false);
    list.data = data.rows;
    if(data.total > 0)list.total = data.total;
}

const handleFilter = ()=>{
  list.page = 1;
  queryPagingData();
}

const handleSelectionChange = (rows: Array<Paging{{$Class}}>)=> {
    selectedRows.value = rows;
    selectedIds.value = [];
    rows.map(row=> selectedIds.value.push(row.{{.table.Pk}}));
}

// 返回
const handleBack = ()=>{
   // this.$store.dispatch('delView', this.$route)
   router.go(-1);
}

// 新增数据
const handleCreate = ()=> openForm(null,"新增{{.table.Comment}}")

// 编辑数据
const handleEdit = (row:Paging{{$Class}})=>openForm(row,"编辑{{.table.Comment}}")

// 打开表单
const openForm = (row:Paging{{$Class}},title:string)=>{
  /** #! 在新的tab页上打开临时页面 */
  // router.push({path:"../{{name_path .table.Name}}/detail",query:{title,{{.table.PkProp}}}});
  dialog.title = title;
  dialog.params = row.{{.table.PkProp}};
  dialog.modal = {{$Class}}Modal;
}

// 参数state为true时,重置模态框并刷新数据,args接受传入的参数
const refresh = ({state = 0,close = true,data = {}})=>{
    if(close)dialog.modal = null;
    if(state)queryPagingData(data);
}

const handleDelete = (idx:number,row:Paging{{$Class}}) => {
    MessageBox.confirm('执行此操作数据无法恢复,是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(async () => {
        if(requesting.value)return;requesting.value=true;$index.value = idx;
        let ret = await (row && row.{{.table.Pk}}?
          delete{{$Class}}(row.{{.table.Pk}}):
          batchDelete{{$Class}}(selectedIds.value))
          .finally(()=>requesting.value=false;$index.value = -1;);
        let {errCode,errMsg} = parseResult(ret.data);
        if (errCode === 0) {
          Message.success({message:'删除成功',duration:2000,onClose:queryPagingData});
        }else{
          MessageBox.alert(errMsg,"提示")
        }
    }).catch((ex) => {
      ex !=="cancel" && MessageBox.alert(ex.message,"错误")
    });
}
</script>
