#!kind:2#!target:vue-pro/ext/compose/data.ts


import { ModelRef, Ref, UnwrapNestedRefs, nextTick, reactive, ref } from "vue"
import { parseResult, Message, Supplier } from "../utils"

// 列表引用
export interface ListRef<R> {
  // 返回唯一主键
  primary: (r: R) => any
  // 数据表引用
  tableRef: () => Ref
  // 总条数
  total?: number
  // 页码
  page: number
  // 数量
  size: number
  // 分页数据
  rows?: Array<R>
  // 当前行索引
  rowIndex?: number
  // 加载数据状态
  loading?: boolean
  // 操作请求状态
  requesting?: boolean
  // 已选择的行
  selectedRows: Array<R>
}

export interface PagingResult<E> {
  rows: Array<E>
  total: number
  extra?: any
}

/** 列表数据引用类型 */
type ERef<E> = UnwrapNestedRefs<ListRef<E>>

/** 分页接口响应 */
type PR<E> = { data: PagingResult<E> }

/** 数据源 */
type DS<P, R> = (page: number, size: number, p: P) => Promise<R>
/** 函数 */
type EF<E, R> = (e: E, ...args: any) => Promise<R>
/** 更新函数 */
type UF<E, R> = (pk: any, e: E) => Promise<R>
/** 删除删除 */
type DF<P, R> = EF<Array<P>, R>
/** 操作结果 */
type Result<R> = R & { code: number; msg: string }

/**
 * 查询数据
 *
 * @param dataRef 数据引用
 * @param fetch 获取数据方法
 * @param queryParams 查询参数
 * @returns 接口响应数据
 *
 * @type E 数据行类型
 * @type P 查询参数类型
 * @type R 响应对象类型
 */
async function queryData<E, P, R extends { data: PagingResult<E> }>(opt: {
  fn: DS<P, R>
  data: P
  dataRef: ERef<E>
}): Promise<R | undefined> {
  const dr = opt.dataRef
  if (!dr.loading) {
    dr.loading = true
    const ret = await opt.fn(dr.page, dr.size, opt.data).finally(() => (dr.loading = false))
    const { data } = ret
    dr.rows = data.rows as any
    if (data.total > 0) dr.total = data.total
    revertSelectionState(dr)
    return ret
  }
}

/**
 * 删除数据
 *
 * @param dataRef 数据引用
 * @param deleteFn 处理删除函数
 * @param p 查询参数
 * @returns 响应信息
 *
 * @type E 数据行类型
 * @type P 查询参数类型
 * @type R 响应对象类型
 */
async function deleteData<E, P, R extends { data: any }>(opt: {
  fn: DF<P, R>
  ids: Array<P>
  dataRef: ERef<E>
}): Promise<Result<R>> {
  const dataRef = opt.dataRef
  if (dataRef.requesting) throw new Error("requesting")
  dataRef.requesting = true
  let ret = await opt.fn(opt.ids).finally(() => {
    dataRef.requesting = false
    dataRef.rowIndex = 0
  })
  return {
    ...ret.data,
    ...parseResult(ret.data)
  }
}

// 回显选中状态
async function revertSelectionState<E>(queryData: ERef<E>) {
  await nextTick()
  const { tableRef, rows: dataList, primary, selectedRows } = queryData
  if (tableRef) {
    const { toggleRowSelection, clearSelection } = tableRef().value
    if (selectedRows) {
      dataList
        ?.filter(
          (a) =>
            selectedRows.findIndex((b) => {
              if (typeof a != typeof b) {
                console.error(`selectedRows基础类型为${typeof b}, 无法自动选中行`)
              }
              return primary(a as E) == primary(b as E)
            }) != -1
        )
        .map((a) => toggleRowSelection(a, true))
    } else {
      clearSelection()
    }
  }
}

/**
 * 更新选中行函数，当选中数据变化时发生
 *
 * @param queryData 查询对象
 * @param rows 已选中的行
 * @param row 选中或取消的单行，如果全选，则为空
 *
 * @type E 数据行类型
 */
function onSelectionChange<E>(queryData: ERef<E>, rows: Array<E>, row?: E) {
  const { primary, selectedRows, rows: dataList } = queryData
  const fi = (r: E) => (queryData.selectedRows || [])?.findIndex((a) => primary(a as E) == primary(r))
  if (row == null) {
    // 全选/取消全选
    if (rows.length > 0) {
      // 全选
      rows.forEach((a) => {
        if (fi(a) == -1) {
          selectedRows.push(a as any)
        }
      })
    } else {
      // 取消全选
      dataList?.forEach((a) => {
        const i = fi(a as E)
        if (i != -1) {
          selectedRows.splice(i, 1)
        }
      })
    }
  } else {
    if (rows.findIndex((a) => primary(a) == primary(row)) >= 0) {
      // 选中单行
      selectedRows.push(row as any)
    } else {
      // 取消选中单行
      selectedRows.splice(fi(row), 1)
    }
  }
}

/**
 * 清除选择的行
 * @param queryData 查询对象
 */
function resetSelection<E>(queryData: ERef<E>) {
  queryData.selectedRows = []
  queryData.tableRef().value.clearSelection()
}

// 弹出提示
function alertTips<R>(ret: Result<R>, onClose?: () => void) {
  const { code, msg } = ret
  if (code) {
    Message.error({ message: msg, duration: 2000 })
  } else {
    Message.success({ message: "操作成功", duration: 2000, onClose })
  }
}

export interface DataTableProps<E> {
  // 资源标识,用于控制权限
  permKey: string
  // 表格
  tableRef: Ref<any>
  // 表格数据
  tableData: UnwrapNestedRefs<ListRef<E>>
  // 查询参数
  queryParams: Ref<any>
  // 搜索栏
  filterRef: Ref<any>
  // 加载状态
  loading: Ref<boolean>
  /**
   * 重置筛选条件
   */
  resetFilter(): void
  /**
   *
   * @param ds 查询数据
   * @param p 查询参数
   * @param page
   */
  onQueryData<P, R extends PR<E | any>>(ds: DS<P, R>, p: P, page?: number): Promise<R | undefined>

  /**
   * 生成查询分页数据函数
   * @param ds 查询数据函数
   * @param extraParam 附加参数
   */
  buildQueryPagingData<P, R extends PR<E | any>>(
    ds: DS<P, R>,
    extraParam?: any
  ): (page?: number) => Promise<R | undefined>
  /**
   * 删除数据
   * @param opt 删除数据函数
   */
  onDelete<P, R extends { data: any }>(opt: {
    fn: DF<P, R>
    index?: number
    row?: E
    tips?: boolean
    onClose?: () => void
  }): Promise<Result<R>>

  /**
   *  当选择的行发生变化时
   * @param rows 选中的多行
   * @param row 当前选中行
   */
  onSelectionChange(rows: Array<E>, row?: E): void

  /**
   * 反选
   */
  onRevertSelection(): void

  /**
   * 清除选择
   */
  onResetSelection(): void
}

/**
 * 生成数据表格的函数
 *
 * @param r 数据引用对象
 * @returns
 */
export function useDataTable<E>(options: {
  primary: Supplier<E, any>
  permKey?: string
  defaultParams?: any
}): DataTableProps<E> {
  const tableRef = ref()
  const filterRef = ref()
  const loading = ref(false)
  const queryParams = ref(options.defaultParams || {})
  const tableData = reactive<ListRef<E>>({
    page: 1,
    size: 20,
    selectedRows: [],
    tableRef: () => tableRef,
    primary: options.primary
  })

  return {
    permKey: options.permKey || "",
    tableRef,
    tableData,
    filterRef,
    queryParams,
    /** 加载状态 */
    loading,
    /**
     * 重置筛选器
     */
    resetFilter() {
      filterRef.value?.resetFields()
    },
    async onQueryData<P, R extends PR<E | any>>(ds: DS<P, R>, p: P, page?: number): Promise<R | undefined> {
      if (page) {
        if (tableData.loading) return
        tableData.page = page
      }
      return queryData({ fn: ds, data: p, dataRef: tableData })
    },

    /**
     * 生成查询分页数据函数
     * @param ds 查询数据函数
     * @param extraParam 附加参数
     */
    buildQueryPagingData<P, R extends PR<E | any>>(ds: DS<P, R>, extraParam?: any) {
      return (page?: number): Promise<R | undefined> => {
        if (page) {
          if (tableData.loading) {
            return Promise.resolve(undefined)
          }
          tableData.page = page
        }
        return queryData({ fn: ds, data: { ...queryParams.value, ...extraParam }, dataRef: tableData })
      }
    },
    async onDelete<P, R extends { data: any }>(opt: {
      fn: DF<P, R>
      index?: number
      row?: E
      tips?: boolean
      onClose?: () => void
    }): Promise<Result<R>> {
      tableData.rowIndex = opt.index
      if (tableData.requesting) {
        throw "cancel"
      }
      var pkArr = []
      if (opt.row) {
        pkArr = [tableData.primary(opt.row)]
      } else {
        pkArr = (tableData.selectedRows || []).map(tableData.primary as any)
      }
      const ret = await deleteData({ fn: opt.fn, ids: pkArr, dataRef: tableData })
      // 触发提示
      opt.tips != false && alertTips(ret, opt.onClose)
      return ret
    },
    onSelectionChange(rows: Array<E>, row?: E) {
      onSelectionChange(tableData, rows, row)
    },
    onRevertSelection() {
      revertSelectionState(tableData)
    },
    onResetSelection() {
      resetSelection(tableData)
    }
  }
}

/** 数据表单 */
export function useDataForm<P, E>(modelValue: ModelRef<P>, data?: E) {
  const formRef = ref()
  const loading = ref(false)
  const f = ref<E & any>(data || ({} as E as any))
  return {
    formData: f,
    modelValue,
    formRef,
    loading,
    /**
     * 构建查询数据函数
     *
     * @template R 返回的数据类型，需要继承自 PR<E>
     * @param ds 数据源函数，类型为 EF<any, { data: E }>
     * @param extra 额外的参数，可选
     * @returns 返回一个返回 Promise<R | undefined> 的函数
     */
    buildQueryData<R extends PR<E>>(ds: EF<any, { data: E }>, extra?: any): () => Promise<R | undefined> {
      return async () => {
        return new Promise(async (resolve, reject) => {
          if (!modelValue || loading.value) {
            reject("已获取到数据，或正在请求中")
          }
          loading.value = true
          try {
            const { data } = await ds(modelValue.value, extra).finally(() => (loading.value = false))
            const { code, msg } = data as any
            if (code && code > 0) {
              Message.error(msg || "error")
            } else {
              f.value = data as any
            }
          } catch (ex) {
            Message.error(ex as any)
          }
          resolve(f.value)
        })
      }
    },
    //@deprecate
    async queryData<R extends PR<E>>(ds: EF<any, { data: E }>, ...args: any): Promise<R | undefined> {
      if (!modelValue || loading.value) return
      loading.value = true
      try {
        const { data } = await ds(modelValue.value, ...args).finally(() => (loading.value = false))
        const { code, msg } = data as any
        if (code && code > 0) {
          Message.error(msg || "error")
        } else {
          f.value = data as any
        }
      } catch (ex) {
        Message.error(ex as any)
      }
    },
    /**
     * 提交表单
     * @param opt 参数
     * @returns
     */
    async onSubmitForm<R extends { data: any }>(
      opt: {
        tips?: boolean
        updateFn?: UF<E, R>
        createFn?: EF<E, R>
        onClose?: () => void
      },
      promise?: { resolve?: (value: any) => void; reject?: (reason: string) => void }
    ): Promise<Result<E>> {
      return new Promise(async (resolve, reject) => {
        const valid = await (formRef.value as any)?.validate()
        if (!valid) {
          reject("数据格式不合法")
        }
        if (loading.value) {
          reject("重复请求,请检查loading是否为true")
        }
        if (modelValue.value && !opt.updateFn) {
          reject("未指定update函数,请确认是否支持更新操作")
        }
        if (!modelValue.value && !opt.createFn) {
          reject("未指定create函数,请确认是否支持新增操作")
        }
        loading.value = true
        const { data } = await (
          modelValue.value ? opt.updateFn!(modelValue.value, f.value as E) : opt.createFn!(f.value as E)
        ).finally(() => (loading.value = false))
        const ret = {
          ...data,
          ...parseResult(data),
          ...(f.value as E)
        }
        // 触发提示
        opt.tips != false && alertTips(ret, opt.onClose)
        resolve(ret)
      }).then((ret: any) => {
        if (promise) {
          // 触发外部传入的resolve和reject, 通过code判断操作是否符合预期
          if (ret.code > 0) {
            promise.reject && promise.reject(ret.msg)
          } else {
            promise.resolve && promise.resolve(ret)
          }
        }
        return ret
      })
    }
  }
}
