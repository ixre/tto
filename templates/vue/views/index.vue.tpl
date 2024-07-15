#!lang:ts＃!name:全功能界面
#!lang:ts#!target:vue/views/{{.table.Title}}/index.vue
{{$entityName := .table.Title}}
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
              <span v-perm="{key: '', roles: ['admin'], visible: true }" @click="queryPagingData">
                <el-button type="primary">搜索</el-button>
              </span>
            </el-form-item>
            <el-form-item class="filter-item">
              <span v-perm="{key: '[权限key]+1', roles: ['admin'], visible: true }" @click="openModal">
                <el-button>新增</el-button>
              </span>
            </el-form-item>
            <el-form-item class="filter-item">
              <span v-show="queryData.selectedRows?.length" @click="onResetSelection">
                <el-button plain>清除选择项({{ "{{queryData.selectedRows?.length}}" }})</el-button>
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
              @select="onSelectionChange" @select-all="onSelectionChange">
        <el-table-column align="center" type="selection" width="40" fixed="left"/>
        <el-table-column type="index" width="50" label="序号"/>
        {{range $i,$c := .columns}} \
        {{if ends_with $c.Name "Time"}} \
        <el-table-column width="140" align="left" label="{{$c.Comment}}" prop="{{lower_title $c.Prop}}" :formatter="formatColTime"/>
        {{else if ends_with $c.Name "State"}} \
        <el-table-column width="140" align="left" label="{{$c.Comment}}">
            <template #default="{row}">
              <span v-for="(it,i) in stateOptions" v-if="it.value === row.{{lower_title $c.Prop}}">{{"{{ it.key }}"}}</span>
            </template>
        </el-table-column>
        {{else if starts_with $c.Name "is"}} \
        <el-table-column width="140" align="left" label="{{$c.Comment}}">
             <template #default="{row}">
                <span v-if="row.{{lower_title $c.Prop}}===1" class="g-txt-green">是</span>
                <span v-else class="g-txt-red">否</span>
             </template>
        </el-table-column>
        {{else if $c.IsPk }} \
        <el-table-column width="60" align="left" label="{{$c.Comment}}" prop="{{lower_title $c.Prop}}"/>
        {{else}} \
        <el-table-column align="left" label="{{$c.Comment}}" prop="{{lower_title $c.Prop}}"/>
        {{end}} \
        {{end}} \

        <el-table-column align="center" width="160" label="操作" fixed="right">
            <template #default="scope">
              <span v-perm="{key: '[权限key]+4', roles: ['admin'], visible: true }" @click="openModal(scope.row)">
                <el-button type="primary" size="small" plain>编辑</el-button>
              </span>
              <span v-perm="{key: '[权限key]+2', roles: ['admin'], visible: true }" @click="handleDelete(scope.$index,scope.row)">
                <el-button type="danger" size="small" plain :loading="queryData.requesting && queryData.rowIndex === scope.$index">删除</el-button>
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
import {{"{"}}{{$entityName}},paging{{$entityName}},delete{{$entityName}} {{"}"}} from '@/api';
import {{$entityName}}Modal from './modal.vue';
import {Message,MessageBox,formatColTime} from "@/ext/utils";
import { showModal, ListRef, useDataTable } from "@/ext/compose"
// 定义排序条件
const sortOptions = [
  {key: "默认排序", value: "{{.table.Pk}} DESC"},
  {key: "按创建时间先后顺序", value: "{{.table.Pk}} ASC"},
];

// 定义状态条件(自定义条件)
const stateOptions = [
  {key: "全部", value: -1},
  {key: "正常", value: 1},
  {key: "停用", value: 0}
];

// 定义查询参数
const queryParams = reactive({
  keyword: "",
  where: "0=0",
  state: stateOptions[0].value,
  order_by: sortOptions[0].value
});

// {{.table.Comment}}数据映射类
interface R extends {{$entityName}}{}

const tableRef = ref();
const queryData = reactive<ListRef<R> & { tableHeight: number }>({
  page: 1,
  size: 20,
  selectedRows: [],
  tableRef: () => tableRef,
  primary: (row) => row.id,
  tableHeight: 0
})
const { onQueryData, onDelete, onSelectionChange, onResetSelection } = useDataTable(queryData)

onMounted(()=>{
    queryPagingData();
    nextTick(()=>{
      const rec = (tableRef.value as any).$el.getBoundingClientRect();
      queryData.tableHeight =  document.documentElement.clientHeight - rec.top - 80;
    })
})

// 读取分页数据
function queryPagingData(page?:number){
    onQueryData(paging{{$entityName}}, queryParams, page || 1)
}

// 打开表单
const showModal = async (row?:R)=>{
   const data = await showModal({{$entityName}}Modal,{
    modelValue:row?.{{.table.PkProp}}
  },{title: row ? "编辑{{.table.Comment}}" : "新增{{.table.Comment}}"});
  if(data)queryPagingData();
}


const handleDelete = (idx?:number,row?:R) => {
    MessageBox.confirm('执行此操作数据无法恢复,是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(async () => {
        await onDelete({
          index,
          row,
          fn: delete{{$entityName}},
          onClose: queryPagingData
        })
    })
}
</script>
