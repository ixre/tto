#!kind:2#!target:vue-pro/ext/compose/modal.ts



import { ElDialog, ElButton } from "element-plus";
import { Component, h, ref, render } from "vue";

// 模态框配置
export interface ModalOption {
  // 标题
  title?: string | Component;
  // 模态框宽度,默认:50%
  width?: string | number;
  // 模态框宽度,高度:50%
  height?: string | number;
  // 是否可拖拽
  draggable?: boolean;
  // 是否显示关闭按钮
  showClose?: boolean;
  // 点击模态框关闭
  closeOnClickModal?: boolean;
  // 取消文本
  cancelLabel?: string;
  // 确定文本
  confirmLabel?: string;
  // 是否显示控制按钮(默认:显示)
  showControl?: boolean;
  // 自定义控制按钮内容, @apply: 用于生成绑定参数的函数
  controlSlot?: (modalRef: ModalRef<any>, apply: (fn?: Function) => void) => Component | string;
}

/**  模态框子组件引用
 * @type {C} 组件或字符
 */
export type ModalRef<R> = any & {
  submit?: (r: R) => void
  cancel?: () => void
};

/**
 * 弹出模态框, 样式伪类: `dyn-modal`
 * @param component 组件
 * @param opt 配置
 *
 * @example
 * await showModal("haha");
 *
 * @example
 * const data = await showModal(TProjectModal,{modelValue: 2},{
 *  title: "创建连接",
 *  width: 600,
 *  controlSlot(modalRef: ModalRef<any>,apply) {
 *     return h("div", { style: "display:flex; align-content: space-between" }, [
 *         h("div", {}, [
 *               h(ElButton, { type: "primary", onClick: () => apply(modalRef.dbConnection) }, "测试连接")
 *         ]),
 *         h("div", { style: "flex:1" }, [
 *           h(ElButton, { onClick: () => apply(modalRef.reset) }, "取消"),
 *           h(ElButton, { type: "primary", onClick: () => apply(modalRef.submit) }, "确定")
 *         ])
 *     ]);
 *   },
 * });
 *
 * @type {P} 子组件props
 * @type {M} 子组件modelValue
 * @type {C} 组件或字符
 * @type {O} 模态框选项,可使用ElDialog其他配置
 * @type {R} 返回类型
 */
export function showModal<
  C extends Component | string,
  P,
  M,
  O extends ModalOption,
  R
>(
  component: C,
  props?: P & { modelValue?: M },
  opt?: O
): Promise<R | undefined> {
  return new Promise((resolve, reject) => {
    const modelValue = true; // ElDialog是否显示
    const dialogProps = {
      ...opt,
      modelValue: modelValue,
      title: opt?.title,
      class: "dyn-modal",
      width: opt?.width,
      height: opt?.height,
      showClose: opt?.showClose,
      draggable: opt?.draggable ?? true,
      closeOnClickModal: opt?.closeOnClickModal ?? false,
      destroyOnClose: true,
    };

    // 接收close事件
    const closeHandler = (r?: R) => {
      // 返回错误
      if (r instanceof Error) return reject(r);
      removeChild();
      resolve(r);
    };

    // 模态框组件提交方法,需在子组件中使用defineExpose({submit(),reset()}),供外部的按钮调用
    let elRef = ref<ModalRef<R>>();
    if (props?.hasOwnProperty("ref")) {
      // 引用外部传入的Ref,方便外部操作Modal内部组件
      elRef = (props as any)["ref"];
    }
    // 添加元素到DOM
    const container = document.createElement("div");
    document.body.appendChild(container);

    // 关闭并移出DOM
    const removeChild = () => {
      // 隐藏模态框
      if (vNode.component?.props.modelValue) {
        vNode.component!.props.modelValue = false;
      }
      // 移除dom
      render(null, container);
      document.body.removeChild(container);
    };

    // 创建节点
    var vNode = h(
      ElDialog as Component,
      {
        ...dialogProps,
        onClose: closeHandler,
      },
      {
        default: () => {
          let type = typeof component;
          if (type == "string") {
            return h("div", component as string);
          } else {
            // 组件使用close-emit, emit必须以"on"开头,在子组件中用emits(`close`)
            // 添加resolve和reject两个emit
            return h(component as Component, {
              ...props,
              ref: elRef,
              onClose: closeHandler,
              onResolve: resolve,
              onReject: reject,
            });
          }
        },
        header: () => {
          // #header插槽
          if (opt && opt.title) {
            if (typeof opt.title == "string") {
              // 纯字符标题
              return h("span", opt.title as string);
            }
            // 使用组件作为标题
            return h(opt.title as Component, null);
          }
          return null;
        },
        footer: () => {
          // #footer插槽
          if (opt?.showControl == false) {
            // 不显示控制栏
            return null;
          }

          // 解析调用函数
          const applyResolveFunc = (fn?: Function) => {
            fn ? fn((v: any) => {
              // 当解析到预期结果时,关闭弹窗
              removeChild()
              resolve(v)
            }, reject) : closeHandler()
          }
          if (opt?.controlSlot) {
            const slotComponent = opt?.controlSlot(elRef.value, applyResolveFunc);
            // 自定义控制栏插槽
            if (typeof slotComponent === "string") {
              return h("div", slotComponent as string);
            }
            return slotComponent;
          }
          // 默认显示取消和确定按钮,确定按钮调用导出的submit方法,取消按钮调用导出的reset方法
          return [
            h(
              ElButton as Component,
              { onClick: () => applyResolveFunc(elRef.value.cancel) },
              opt?.cancelLabel || "取消"
            ),
            h(
              ElButton as Component,
              {
                type: "primary",
                onClick: () => applyResolveFunc(elRef.value.submit),
              },
              opt?.confirmLabel || "确定"
            ),
          ];
        },
      }
    );
    render(vNode, container);
  });
}
