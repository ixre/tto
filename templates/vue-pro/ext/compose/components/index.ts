#!kind:2#!target:vue-pro/ext/compose/components/index.ts

import { Ref } from "vue";
import { ElementProps } from "./types";
import { GenericOption } from "@/ext/types";

/**
 * 生成选项代理对象,包括远程读取数据的实现
 *
 * @param props
 * @param options
 * @returns
 */
export const applyOptionRef = async (
  props: ElementProps,
  options: Ref<Array<GenericOption>>
): Promise<Ref<Array<GenericOption>>> => {
  const propValue = props.options as Ref<GenericOption[]>;
  if (propValue.value) {
    // 使用Ref类型数据
    options = propValue;
  }
  if (typeof props.options === "function") {
    // 加载动态数据
    const ret = await props.options();
    options.value = ret;
  } else if (Array.isArray(props.options)) {
    // 绑定固定数据
    options.value = props.options;
  }
  //console.log("---options", propValue.value, options, props.options)
  return options;
};
