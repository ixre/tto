#!kind:2#!target:vue-pro/ext/compose/components/types.ts
import { GenericOption } from "@/ext/types";
import { parseTime } from "@/ext/utils";
import { Ref } from "vue";

export type ColumnDisplayType = "text" | "image" | "selection" | "index";

/**
 * 列配置属性
 */
export interface ColumnProps {
  /**
   * 列的唯一标识
   */
  prop?: string;
  /**
   * 控制列是否显示
   * @returns 返回函数
   */
  show?: () => boolean;
  /**
   * 列头插槽
   */
  headerSlot?: string;
  /**
   * 插槽,支持列自定义显示
   */
  slot?: string;
  /**
   * 展示类型
   */
  type?: ColumnDisplayType;
  /**
   * 列标题
   */
  label?: string;
  /**
   * 固定列头
   */
  fixed?: string | boolean;
  /**
   * 对齐方式,默认left
   */
  align?: "left" | "center" | "right";
  /**
   * 宽度
   */
  width?: string | number;

  /**
   * 占位符,当数据为空时显示
   */
  defaultValue?: string;

  /**
   * 格式化器
   * @param value 值
   * @param row 行
   * @returns 内容
   */
  formatter?: (value: any, row?: any) => string;
}

/**
 * 列构建器
 */
export interface IColumnBuilder {
  /**
   * 创建序号列
   * @param opt 属性
   * @returns
   */
  index(opt?: ColumnProps): ColumnProps;
  /**
   * 创建选择列
   * @param opt 属性
   * @returns
   */
  selection(opt?: ColumnProps): ColumnProps;

  /**
   * 创建文本列
   * @param field 字段
   * @param opt 属性
   */
  text(field: string, opt?: ColumnProps): ColumnProps;

  /**
   * 创建图片列
   * @param field 字段
   * @param opt 属性
   */
  image(field: string, opt?: ColumnProps): ColumnProps;

  /**
   * 创建可选值列
   * @param field 字段
   * @param values 选项,如: {1:"男",2:"女"}
   * @param opt 属性
   */
  option(field: string, opt: ColumnProps, values: OptionValues): ColumnProps;
  /**
   * 创建时间戳列
   * @param field 字段
   * @param opt 属性
   */
  timestamp(field: string, opt?: ColumnProps): ColumnProps;

  /**
   * 自定义列
   * @param opt 属性
   */
  custom(opt: ColumnProps): ColumnProps;
}

export interface OptionValues {
  [key: string | number]: any;
}
/**
 * 列构建器
 */
export const Columns: IColumnBuilder = {
  index(opt?: ColumnProps): ColumnProps {
    return {
      type: "index",
      label: "序号",
      width: 80,
      fixed: true,
      align: "right",
      ...opt,
    };
  },
  selection(opt?: ColumnProps): ColumnProps {
    return {
      type: "selection",
      width: 60,
      align: "center",
      fixed: true,
      ...opt,
    };
  },
  text(field: string, opt?: ColumnProps): ColumnProps {
    return {
      width: 120,
      align: "left",
      ...opt,
      prop: field,
    };
  },
  option(field: string, opt: ColumnProps, values: OptionValues): ColumnProps {
    return {
      width: 80,
      align: "left",
      ...opt,
      prop: field,
      formatter: (value, row) => {
        return values[value] || value;
      },
    };
  },
  timestamp(field: string, opt?: ColumnProps): ColumnProps {
    return {
      width: 180,
      align: "left",
      ...opt,
      prop: field,
      formatter: (value, row) => parseTime(value) || "nnn",
    };
  },
  custom: (opt: ColumnProps): ColumnProps => {
    return {
      ...opt,
    };
  },
  image(field: string, opt?: ColumnProps): ColumnProps {
    return {
      align: "center",
      ...opt,
      prop: field,
      type: "image",
    };
  },
};

type SizeType = "medium" | "small" | "mini";

/** 表单项类型  */
type FormElementType =
  | "input"
  | "number"
  | "select"
  | "date"
  | "editor"
  | "datetime"
  | "checkbox"
  | "radio"
  | "switch"
  | "cascader"
  | "textarea"
  | "slot"
  | "upload";

/**
 * 元素基本属性
 */
export interface ElementNormalProps {
  // 字段
  prop: string;
  // 标签
  label?: string;
  // 组件类型
  type: FormElementType;
  // 是否可见
  visible?: boolean;
  // 是否禁用
  disabled?: boolean;
  // 自定义样式
  style?: Record<string, string> | string;
  /** 其他配置 */
  [key: string]: any;
}

/**
 * 表单元素属性
 */
export interface FormElementProps {
  // 占位提示文本
  placeholder?: string;
}

/**
 * 选择(含check/radio)元素属性
 */
export interface SelectElementProps {
  /** 选项 */
  options?:
    | GenericOption[]
    | (() => Promise<GenericOption[]>)
    | Ref<Array<GenericOption>>;
}

/**
 * 输入元素属性
 */
export interface InputElementProps {
  /** 是否可清除输入值 */
  clearable?: boolean;
  /** 输入框类型 */
  onInput?: ((value: string) => any) | undefined;
}

/**
 * 元素属性
 */
export type ElementProps = ElementNormalProps &
  SelectElementProps &
  FormElementProps &
  InputElementProps;
