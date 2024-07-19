#!kind:2#!target:vue-pro/ext/compose/components/elements/SelectElement.vue

<template>
  <div class="form__select--element">
    <el-select
      style="width: 100%"
      :filterable="isRemote"
      :remote="isRemote"
      :remote-method="getRemoteOptions"
      v-bind="$attrs"
      v-on="$listeners"
      v-model="modelValue"
    >
      <el-option
        v-for="item in options"
        :key="item.value"
        :label="item.label"
        :value="item.value"
        :disabled="item.disabled"
      >
        <span style="float: left1; margin-right: 5px">{{ item.label }}</span>
        <el-badge
          :value="item.count"
          :max="999"
          v-if="item.count && item.count > 0"
        ></el-badge>
      </el-option>
    </el-select>
  </div>
</template>
<script setup lang="ts">
import { onBeforeMount, Ref, ref } from "vue";
import { ElementProps } from "../types";
import { GenericOption } from "@/ext/types";
import { applyOptionRef } from "..";

const props = defineProps<ElementProps>();
let options = ref<GenericOption[]>([]);
const isRemote = ref(false);

const modelValue = defineModel();

onBeforeMount(async () => {
  isRemote.value = typeof props.options === "function";
  options = await applyOptionRef(props, options);
});
</script>
