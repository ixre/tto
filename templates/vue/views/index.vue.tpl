#!lang:ts＃!name:全功能界面
#!lang:ts#!target:vue/views/{{.table.Title}}/{{.table.Title}}Index.vue
{{$Class := .table.Title}}
{{$pkType := type "ts" .table.PkType}}
{{$pk := .table.Pk}}
<template>
<div class="app-container mod-grid-container">
    <!-- 查询和其他操作 -->
    <div class="filter-container mod-grid-header">
        <div class="it mod-grid-header-filter">
          <el-form :inline="true">
            <el-form-item label="状态" class="filter-item">
              <el-select v-model="queryParams.state" class="filter-select"
                         @change="queryPagingData">
                <el-option v-for="(it,i) in stateOptions" :key="i" :label="it.key" :value="it.value"/>
              </el-select>
            </el-form-item>
            <el-form-item label="关键词" class="filter-item">
              <el-input v-model="queryParams.keyword" clearable class="filter-input" placeholder="请输入关键词"/>
            </el-form-item>
            <el-form-item label="排序方式" class="filter-item">
              <el-select v-model="queryParams.order_by" class="filter-select">
                <el-option v-for="(it,i) in sortOptions" :key="i" :value="it.value" :label="it.key"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item class="filter-item">
              <span v-perm="{key: '', roles: ['admin'], visible: true }" @click="handleFilter">
                <el-button type="primary">搜索</el-button>
              </span>
            </el-form-item>
            <el-form-item class="filter-item">
              <span v-perm="{key: '[权限key]+1', roles: ['admin'], visible: true }" @click="handleCreate">
                <el-button>新增</el-button>
              </span>
            </el-form-item>
            <el-form-item class="filter-item">
              <span v-show="queryData.selectedRows?.length" @click="resetSelections">
                <el-button type="plain">清除选择项({{ "{{queryData.selectedRows?.length}}" }})</el-button>
              </span>
              <span v-perm="{key: '', roles: ['admin'], visible: true }" @click="handleDelete">
                <el-button v-show="queryData.selectedRows?.length" type="danger" :loading="queryData.requesting">删除</el-button>
              </span>
            </el-form-item>
          </el-form>
        </div>
    </div>
    <div class="mod-grid-body">
    <!-- 表单数据 -->
    <el-table ref="tableRef" v-loading="queryData.loading" :data="queryData.dataList" class="mod-grid-table"
              border :height="queryData.tableHeight" fit :highlight-current-row="false"
              row-key="{{$pk}}" empty-text="暂无数据"
              @select="handleSelectionChange" @select-all="handleSelectionChange">
        <el-table-column align="center" type="selection" width="40" fixed="left"/>
        <el-table-column type="index" width="50" label="序号"></el-table-column>
        {{range $i,$c := .columns}} \
        {{if ends_with $c.Name "_time"}} \
        <el-table-column width="140" align="left" label="{{$c.Comment}}" prop="{{$c.Name}}" :formatter="formatColTime"/>
        {{else if ends_with $c.Name "state"}} \
        <el-table-column width="140" align="left" label="{{$c.Comment}}">
            <template #default="{row}">
              <span v-for="(it,i) in stateOptions" v-if="it.value === row.state">{{"{{ it.key }}"}}</span>
            </template>
        </el-table-column>
        {{else if starts_with $c.Name "is_"}} \
        <el-table-column width="140" align="left" label="{{$c.Comment}}">
             <template #default="{row}">
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
            <template #default="scope">
              <span v-perm="{key: '[权限key]+4', roles: ['admin'], visible: true }" @click="handleEdit(scope.row)">
                <el-button type="primary" size="mini" icon="el-icon-edit-outline">编辑</el-button>
              </span>
              <span v-perm="{key: '[权限key]+2', roles: ['admin'], visible: true }" @click="handleDelete(scope.$index,scope.row)">
                <el-button type="danger" size="mini" icon="el-icon-delete" :loading="queryData.requesting && queryData.rowIndex === scope.$index">删除</el-button>
              </span>
            </template>
        </el-table-column>
    </el-table>
    </div>
    <div class="mod-grid-footer">
      <!-- 分页 -->
      <el-pagination v-show="queryData.total && queryData.total > 0" :total="queryData.total" :page-size.sync="queryData.size"
        :current-page="queryData.page" :page-sizes="[10, 20, 50, 100]"
        @current-change="queryPagingData" @size-change="(v)=>queryData.size = v"
        background layout="prev, pager, next,sizes,jumper,total"></el-pagination>
    </div>
</div>
</template>
<script lang="ts" setup>
import {onMounted, reactive, ref, nextTick} from "vue";
import {Paging{{$Class}},queryPaging{{$Class}},delete{{$Class}} } from '../../api';
import {{$Class}}Modal from './{{$Class}}Modal.vue';
import {Message,MessageBox,formatColTime} from "../../utils";
import {showModal,ListDataRef,queryDataList, deleteData,onSelectionChange,onResetSelection} from "../../components";

const tableRef = ref(null);

const queryData = reactive<ListDataRef<Paging{{$Class}}>
    & { tableHeight: number }>({
        page: 1, size: 20,selectedRows: [],tableRef: () => tableRef, 
        primary: (row) => row.{{$pk}},
        tableHeight: 0
    });

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
  order_by: sortOptions[0].value
});

onMounted(()=>{
    queryPagingData();
    nextTick(()=>{
      const rec = (tableRef.value as any).$el.getBoundingClientRect();
      queryData.tableHeight =  document.documentElement.clientHeight - rec.top - 80;
    })
})

// 读取分页数据
const queryPagingData = async (page?:number)=> {
  if (queryData.loading) return;
  queryData.page = page || 1;
  await queryDataList(queryData,queryPaging{{$Class}},queryParams);
};

const handleFilter = ()=>{
  queryData.page = 1;
  queryPagingData();
};

const handleSelectionChange = (rows: Array<Paging{{$Class}}>,row?:Paging{{$Class}})=>onSelectionChange(queryData, rows, row)
const resetSelections = () => onResetSelection(queryData)

// 新增数据
const handleCreate = ()=> openForm("新增{{.table.Comment}}")

// 编辑数据
const handleEdit = (row:Paging{{$Class}})=>openForm("编辑{{.table.Comment}}",row)

// 打开表单
const openForm = async (title:string,row?:Paging{{$Class}})=>{
   const data = await showModal({{$Class}}Modal,{
    modelValue:row?.{{.table.PkProp}}
  },{title});
  if(data)queryPagingData();
}


const handleDelete = (idx?:number,row?:Paging{{$Class}}) => {
    MessageBox.confirm('执行此操作数据无法恢复,是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(async () => {
        queryData.rowIndex = idx;
        if (queryData.requesting) return;
        const pkArr = row ? [row.{{$pk}}] : (queryData.selectedRows || []).map(r => r.{{$pk}});
        const { errCode, errMsg } = await deleteData(queryData,delete{{$Class}}, pkArr);
        if (errCode !== 0) {
          throw new Error(errMsg);
        }
        Message.success({message:'删除成功',duration:2000,onClose:queryPagingData});
    }).catch((ex) => {
        ex !=="cancel" && MessageBox.alert(ex.message,"错误")
    });
}
</script>
