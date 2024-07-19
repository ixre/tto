#!lang:ts＃!name:全功能界面
#!lang:ts#!target:vue-pro/views/{{.table.Title}}/index.vue
{{$entityName := .table.Title}}
{{$pkType := type "ts" .table.PkType}}
{{$pk := .table.Pk}}
<template>

<ExtDataView :data="ext" :query="queryPagingData" @open-modal="openModal">
  <template #filter="{ params }">
    <ExtFormElement v-for="(ele, i) in searchColumns" :key="i" v-model="params[ele.prop]" :config="ele" />
  </template>
  
  <ExtDataTable :columns="columns">
    <template #op="scope">
      <div class="flex gap align-center justify-end">
        <span v-perm="{ key: `${permKey}+4`, roles: ['admin'], visible: true }" @click="openModal(scope.row)">
          <el-button type="primary" size="small" plain>编辑</el-button>
        </span>
        <span v-perm="{ key: `${permKey}+2`, roles: ['admin'], visible: true }"
          @click="handleDelete(scope.$index, scope.row)">
          <el-button type="danger" size="small" plain
            :loading="tableData.requesting && tableData.rowIndex === scope.$index"
            >删除</el-button
          >
        </span>
      </div>
    </template>
  </ExtDataTable>
</ExtDataView>
</template>
<script lang="ts" setup>
/**
 * {{.table.Comment}}列表页
 * @author: {{.global.user}}
 * @description: {{.global.time}}
 */
import {onMounted} from "vue";
import {{"{"}}{{$entityName}},paging{{$entityName}},delete{{$entityName}} {{"}"}} from '@/api';
import {{$entityName}}Modal from './modal.vue';
import {MessageBox} from "@/ext/utils";
import { ExtDataView, ExtDataTable,ExtFormElement, Columns, showModal, useDataTable, ColumnProps,ElementProps } from "@/ext/compose"

// 定义数据列, 如果不指定,则默认使用ElementPlus列的定义方式
const columns: Array<ColumnProps> = [
  Columns.selection(),
  Columns.index(),
  {{range $i,$c := .columns}} \
  {{if or (contain $c.Name "image") (contain $c.Name "img")}} \
    Columns.image("{{lower_title $c.Prop}}", { label: "{{$c.Comment}}", width: 75 }),
  {{else if ends_with $c.Name "Time"}} \
    Columns.timestamp("{{lower_title $c.Prop}}", { label: "{{$c.Comment}}",width:140 }),
  {{else if or (ends_with $c.Name "status") (ends_with $c.Name "state")}} \
    Columns.option("{{lower_title $c.Prop}}", { label: "{{$c.Comment}}"},{0:"停用",1:"正常"}),
  {{else if starts_with $c.Name "is"}} \
    Columns.option("{{lower_title $c.Prop}}", { label: "{{$c.Comment}}"},{0:"否",1:"是"}),
  {{else if $c.IsPk }} \
    Columns.text("{{lower_title $c.Prop}}", { label: "{{$c.Comment}}", width: 60, align: "left"}),
  {{else}} \
    Columns.text("{{lower_title $c.Prop}}", { label: "{{$c.Comment}}"}),
  {{end}} \
  {{end}} \
  Columns.custom({ label: "操作", slot: "op", align: "right" })
]

// 搜索列配置
const searchColumns: Array<ElementProps> = [
  /** 定义状态条件(自定义条件) */
  { type: "select", label: "状态", prop: "state", options: [
    {label: "全部状态", value: -1},
    {label: "正常", value: 1},
    {label: "停用", value: 0}
  ]},
  { type: "input", label: "关键词", prop: "keyword", clearable: true },
  /** 定义排序条件 */
  { type: "select", label: "排序", prop: "orderBy", options: [
    {label: "默认排序", value: "{{.table.Pk}} DESC"},
    {label: "按创建时间先后顺序", value: "{{.table.Pk}} ASC"},
  ]},
]

// {{.table.Comment}}数据映射类
interface R extends {{$entityName}}{}

// 创建数据表对象
const ext = useDataTable<R>({
  // 指定主键
  primary: (row) => row.{{lower_title .table.PkProp}},
  // 权限Key
  permKey: "{{.table.Name}}",
  // 默认查询参数
  defaultParams: {keyword: "",state: -1,orderBy: ""}
})
const { tableData, permKey, onDelete } = ext
// 分页数据查询函数
const queryPagingData = ext.buildQueryPagingData(paging{{$entityName}}, {})

onMounted(()=>{
  queryPagingData();
})

/**
 * 打开模态窗口
 * @param row 数据行,当不为空时更新数据,反之新增
 */
const openModal = async (row?:R)=>{
   const data = await showModal({{$entityName}}Modal,{
    modelValue:row?.{{lower_title .table.PkProp}}
  },{title: row ? "编辑{{.table.Comment}}" : "新增{{.table.Comment}}"});
  if(data)queryPagingData();
}

/**
 * 删除数据行
 * @param index 删除指定索引的行,用于单行删除,批量删除传空值
 * @param row 单行删除时传递
 */
const handleDelete = (index?:number,row?:R) => {
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
