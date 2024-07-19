#!kind:2#!target:vue-pro/ext/compose/components/DataTableLayout.vue

<template>
  <!-- 表单数据 -->
  <div class="tableview__table">
    <el-table
      ref="tableRef"
      v-loading="tableData.loading"
      :data="tableData.rows"
      :highlight-current-row="false"
      row-key="id"
      empty-text="暂无数据"
      height="100%"
      @select="onSelectionChange"
      @select-all="onSelectionChange"
    >
      <slot v-if="!props.columns">
        <!-- "当未指定columns参数时,使用兼容模式, 可以套用ElementUI定义列表的方式" /*
        <template #default>
          <el-table-column align="center" type="selection" width="40" fixed="left" />
          <el-table-column type="index" width="50" label="序号" />
          <el-table-column align="left" label="认证姓名" prop="certifiedName" />
          <el-table-column align="left" label="性别" prop="gender" />
          <el-table-column align="left" label="昵称" prop="nickname" />
          <el-table-column align="left" label="工作状态: 1: 离线 2:在线空闲 3: 工作中" prop="workStatus" />
          <el-table-column align="left" label="评分" prop="grade" />
          <el-table-column align="left" label="状态: 1: 正常  2: 锁定" prop="status" />
          <el-table-column width="140" align="left" label="是否认证 0:否 1:是">
            <template #default="{ row }">
              <span v-if="row.isCertified === 1" class="g-txt-green">是</span>
              <span v-else class="g-txt-red">否</span>
            </template>
          </el-table-column>
          <el-table-column align="left" label="创建时间" prop="createTime" />
          <el-table-column align="center" width="160" label="操作" fixed="right">
            <template #op="scope">
              <div class="flex gap align-center justify-center">
                <span v-perm="{ key: '[权限key]+4', roles: ['admin'], visible: true }" @click="openModal(scope.row)">
                  <el-button type="primary" size="small" plain>编辑</el-button>
                </span>
                <span
                  v-perm="{ key: '[权限key]+2', roles: ['admin'], visible: true }"
                  @click="handleDelete(scope.$index, scope.row)"
                >
                  <el-button
                    type="danger"
                    size="small"
                    plain
                    :loading="tableData.requesting && tableData.rowIndex === scope.$index"
                    >删除</el-button
                  >
                </span>
              </div>
            </template>
          </el-table-column>
        </template>
      -->
      </slot>
      <el-table-column
        v-else
        v-for="(column, i) in currentColumns()"
        :key="i"
        :prop="column.prop"
        :label="column.label"
        :type="getColumnType(column.type)"
        :fixed="column.fixed"
        :align="column.align || 'center'"
        :width="column.width"
        :formatter="getColumnFormatter(column)"
      >
        <template v-if="column.headerSlot" #header="row">
          <!-- 列头插槽 -->
          <slot :name="column.headerSlot" v-bind="row" />
        </template>
        <template v-if="column.slot" #default="scope">
          <!-- 列插槽 -->
          <slot :name="column.slot" v-bind="scope" />
        </template>
        <template v-else-if="column.defaultValue" #default="{ row }">
          <!-- 默认内容 -->
          {{ row[column.prop || ""] || column.defaultValue || "" }}
        </template>
        <template v-if="column.type == 'image'" #default="{ row }">
          <img
            class="tableview__table--img"
            :src="row[column.prop || '']"
            alt="图片"
          />
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { inject } from "vue";
import { DataTableProps } from "../data";
import { ColumnProps } from "./types";

const props = withDefaults(
  defineProps<{
    data?: DataTableProps<any>;
    columns?: Array<ColumnProps>;
  }>(),
  {}
);

const tableProps = (props.data || inject("tableData")) as DataTableProps<any>;
const { tableData, tableRef, onSelectionChange } = tableProps;

/**
 * 返回当前显示的列
 */
const currentColumns = () => {
  return props.columns?.filter((a) => !a.show || a.show());
};

/**
 * 获取列格式化函数
 *
 * @param column 列配置对象
 * @returns 格式化函数，如果列配置中没有 formatter 或 prop 属性，则返回 null
 */
const getColumnFormatter = (column: ColumnProps) => {
  if (column.formatter && column.prop) {
    return (row: any) => column.formatter!(row[column.prop!], row);
  }
  return undefined;
};

/**
 * 获取ElementUI的列类型
 * @param type 列类型
 */
const getColumnType = (type?: string): string | undefined => {
  if (!type) {
    return undefined;
  }
  return ["selection", "index"].includes(type!) ? type : "default";
};
</script>
