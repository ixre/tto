#!kind:2#!target:vue-pro/ext/compose/components/elements/Element.vue

<template>
  <!-- 普通 -->
  <el-input
    v-if="config.type === 'input'"
    v-model="modelValue"
    :placeholder="config.placeholder || '请输入'"
    :disabled="config.disabled ?? false"
    :size="size"
    :clearable="config.clearable"
    @input="config.onInput"
  />
  <el-input
    v-if="config.type === 'number'"
    v-model.number="modelValue"
    :placeholder="config.placeholder || '请输入'"
    :disabled="config.disabled ?? false"
    :size="size"
    :clearable="config.clearable"
    @input="config.onInput"
  />
  <!-- 文本域 -->
  <el-input
    v-if="config.type === 'textarea'"
    v-model="modelValue"
    type="textarea"
    :rows="4"
    :disabled="config.disabled ?? false"
    :size="size"
    :placeholder="config.placeholder || '请输入'"
    @input="config.onInput"
  />
  <!-- 下拉 -->
  <SelectElement
    v-if="config.type === 'select'"
    v-model="modelValue"
    :placeholder="config.placeholder || '请选择'"
    :size="size"
    :disabled="config.disabled ?? false"
    :clearable="config.clearable ?? true"
    @input="modelHandle($event, config.key)"
    :options="config.options"
  />
  <!-- 日期选择 -->
  <el-date-picker
    v-if="config.type === 'date'"
    v-model="modelValue"
    align="right"
    :type="config.datetype ?? 'date'"
    :disabled="config.disabled ?? false"
    :value-format="config.dateFormat ?? 'yyyy-MM-dd'"
    :placeholder="config.placeholder || '选择日期'"
    :size="size"
    clearable
    @input="config.onInput"
  />
  <!-- 日期时间 -->
  <el-date-picker
    v-if="config.type === 'datetime'"
    v-model="modelValue"
    type="datetimerange"
    range-separator="至"
    start-placeholder="开始日期"
    end-placeholder="结束日期"
    :disabled="config.disabled ?? false"
    :size="size"
    clearable
    :value-format="config.dateFormat ?? 'yyyy-MM-dd'"
    @input="config.onInput"
  />
  <!-- 开关 -->
  <el-switch
    v-if="config.type === 'switch'"
    v-model="modelValue"
    :active-color="config.activeColor"
    :disabled="config.disabled ?? false"
    :inactive-color="config.inactiveColor"
    :active-text="config.activeText"
    :inactive-text="config.inactiveText"
    :active-value="config.activeValue"
    :inactive-value="config.inactiveValue"
    :size="size"
    @input="config.onInput"
  />
  <!-- 级联选择 -->
  <el-cascader
    v-if="config.type === 'cascader'"
    v-model="modelValue"
    :options="options as any"
    :size="size"
    :disabled="config.disabled ?? false"
    clearable
    @input="config.onInput"
  />
  <!-- 多选 -->
  <el-checkbox-group
    v-if="config.type === 'checkbox'"
    v-model="modelValue"
    :size="size"
    @input="modelHandle($event, config.key)"
  >
    <el-checkbox
      v-for="(item, i) in options"
      :key="i"
      :label="item.label"
      :value="item.value"
    >
      {{ item.label }}
    </el-checkbox>
  </el-checkbox-group>
  <!-- 单选 -->
  <el-radio-group
    v-if="config.type === 'radio'"
    v-model="modelValue"
    :size="props.size"
    @input="modelHandle($event, config.key)"
  >
    <el-radio
      v-for="item in options"
      :key="item.value"
      :label="item.label"
      :value="item.value"
    >
      {{ item.label }}
    </el-radio>
  </el-radio-group>
  <!-- 富文本 -->
  <text-editor
    v-if="config.type === 'editor'"
    :content="modelValue"
    @change="modelValue = $event"
  />
  <!-- 上传 -->
  <img-video-uploader
    v-if="config.type === 'upload'"
    :value="getImgName(modelValue)"
    :max="config.extraCof?.limit || 1"
    @upload-success="uploadSuccess($event, config.key)"
  ></img-video-uploader>
  <!-- 插槽 -->
  <template v-if="config.type === 'slot'">
    <slot :name="config.prop"></slot>
  </template>
</template>
<script setup lang="ts">
import { onBeforeMount, onMounted, ref, watch } from "vue";
// import textEditor from "@/components/text-editor/index.vue"
// import imgVideoUploader from "@/components/img-video-uploader/index.vue"
// import { ErpSelect } from "../erp-select"
// import { ErpSelectProps } from "../erp-select/erp-select.vue"
import { ElementProps } from "../types";
import { SelectElement } from ".";
import { GenericOption } from "@/ext/types";
import { applyOptionRef } from "..";

const modelValue = defineModel();
const props = withDefaults(
  defineProps<{
    config: ElementProps;
    size?: "" | "default" | "small" | "large";
  }>(),
  {
    size: "",
  }
);

// 向选择组件提供数据
let options = ref<GenericOption[]>([]);
const isRemote = ref(false);

onBeforeMount(async () => {
  if (["select", "radio", "checkbox", "cascader"].includes(props.config.type)) {
    isRemote.value = typeof props.config.options === "function";
    options = await applyOptionRef(props.config, options);
  }
});

const formatStyle = (style) => {
  if (!style) return [""];
  return [style];
};
const getImgName = (str) => {
  if (typeof str !== "string") return [];
  let arr = str.split("/");
  return [arr[arr.length - 1]];
};

const uploadSuccess = (e, key) => {
  modelData.value[key] = e && e[0] ? e[0].fullUrl() : "";
  emit("input", modelData.value);
};
</script>
<style scoped lang="less"></style>
