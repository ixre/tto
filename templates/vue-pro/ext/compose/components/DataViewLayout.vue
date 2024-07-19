#!kind:2#!target:vue-pro/ext/compose/components/DataViewLayout.vue

<template>
  <div class="container tableview">
    <div class="tableview__header">
      <slot name="header"></slot>
    </div>
    <!-- 查询和其他操作 -->
    <div class="tableview__filter filter flex justify-between margin-vert">
      <el-form
        ref="filterRef"
        class="flex align-start"
        :model="queryParams"
        :inline="true"
      >
        <div class="filter__group flex flex-wrap">
          <slot name="filter" v-bind:params="queryParams">
            <el-form-item label="关键词" class="filter__item" prop="keyword">
              <el-input
                v-model="queryParams.keyword"
                clearable
                class="filter-input"
                placeholder="请输入关键词"
              />
            </el-form-item>
          </slot>
        </div>
        <el-form-item class="filter__item filter__action">
          <span
            v-perm="{
              key: `${props.data.permKey}+1`,
              roles: ['admin'],
              visible: true,
            }"
            class="filter__item"
            @click="() => props.query()"
          >
            <el-button type="primary">搜索</el-button>
          </span>
          <span class="filter__item" @click="resetFilter">
            <el-button v-permission="['admin']" class="filter__item" plain
              >重置</el-button
            >
          </span>
        </el-form-item>
      </el-form>
    </div>
    <!-- 表单数据 -->
    <div class="tableview__main">
      <div class="tableview__toolbar filter flex justify-between">
        <slot name="toolbar">
          <!-- 工具条 -->
          <div class="flex filter__group">
            <div class="filter__item">
              <span
                v-perm="{
                  key: `${props.data.permKey}+1`,
                  roles: ['admin'],
                  visible: true,
                }"
                @click="openModal"
              >
                <el-button plain type="primary">新增</el-button>
              </span>
            </div>
          </div>
          <div class="flex filter__group">
            <div class="filter__item">
              <span
                v-show="tableData.selectedRows?.length"
                @click="onResetSelection"
              >
                <el-button plain
                  >清除选择项({{ tableData.selectedRows?.length }})</el-button
                >
              </span>
              <span
                v-perm="{
                  key: `${props.data.permKey}+2`,
                  roles: ['admin'],
                  visible: true,
                }"
                @click="handleDelete"
              >
                <el-button
                  v-show="tableData.selectedRows?.length"
                  type="danger"
                  :loading="tableData.requesting"
                  >删除</el-button
                >
              </span>
            </div>
          </div>
        </slot>
      </div>
      <!-- 表单数据 -->
      <slot />

      <div class="tableview__footer flex justify-between">
        <slot name="footer">
          <span class="tableview__footer--total"
            >共查询到结果 {{ tableData.total || 0 }}条</span
          >
          <!-- 分页 -->
          <el-pagination
            v-show="tableData.total && tableData.total > 0"
            :total="tableData.total"
            :page-size.sync="tableData.size"
            :current-page="tableData.page"
            :page-sizes="[10, 20, 50, 100]"
            @current-change="props.query"
            @size-change="(v) => (tableData.size = v)"
            background
            small
            layout="prev, pager, next, sizes, jumper"
          ></el-pagination>
        </slot>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { provide, toRefs } from "vue";
import { DataTableProps } from "../data";

const props = withDefaults(
  defineProps<{
    data: DataTableProps<any>;
    query: (page?: number) => void;
  }>(),
  {}
);

const { tableData, queryParams, filterRef } = toRefs(props.data);
const { resetFilter, onResetSelection } = props.data;

const emit = defineEmits(["open-modal"]);

// 向下级数据表提供数据
provide("tableData", props.data);

const openModal = (row?: any) => {
  emit("open-modal", row);
};
</script>
