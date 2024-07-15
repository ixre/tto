#!kind:2#!target:vue/ext/compose/data.ts

import { Ref, nextTick } from "vue"
import { parseResult, Message } from "../utils"


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

/** 分页接口响应 */
type PR<E> = { data: PagingResult<E> }

/** 数据源 */
type DS<P, R> = (page: number, size: number, p: P) => Promise<R>

/** 函数 */
type EF<E, R> = (e: E) => Promise<R>

/** 更新函数 */
type UF<E, R> = (pk: any, e: E) => Promise<R>

/** 删除删除 */
type DF<P, R> = EF<Array<P>, R>

/** 操作结果 */
type Result<R> = R & { errCode: number; errMsg: string }

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
  dataRef: ListRef<E>
}): Promise<R | undefined> {
  const dr = opt.dataRef
  if (!dr.loading) {
    dr.loading = true
    const ret = await opt.fn(dr.page, dr.size, opt.data).finally(() => (dr.loading = false))
    const { data } = ret
    dr.rows = data.rows
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
  dataRef: ListRef<E>
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
async function revertSelectionState<E>(queryData: ListRef<E>) {
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
              return primary(a) == primary(b)
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
function onSelectionChange<E>(queryData: ListRef<E>, rows: Array<E>, row?: E) {
  const { primary, selectedRows, rows: dataList } = queryData
  const fi = (r: E) => (queryData.selectedRows || [])?.findIndex((a) => primary(a) == primary(r))
  if (row == null) {
    // 全选/取消全选
    if (rows.length > 0) {
      // 全选
      rows.forEach((a) => {
        if (fi(a) == -1) {
          selectedRows.push(a)
        }
      })
    } else {
      // 取消全选
      dataList?.forEach((a) => {
        const i = fi(a)
        if (i != -1) {
          selectedRows.splice(i, 1)
        }
      })
    }
  } else {
    if (rows.findIndex((a) => primary(a) == primary(row)) >= 0) {
      // 选中单行
      selectedRows.push(row)
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
function resetSelection<E>(queryData: ListRef<E>) {
  queryData.selectedRows = []
  queryData.tableRef().value.clearSelection()
}

// 弹出提示
function alertTips<R>(ret: Result<R>, onClose?: () => void) {
  const { errCode, errMsg } = ret
  if (errCode) {
    Message.error({ message: errMsg, duration: 2000 })
  } else {
    Message.success({ message: "操作成功", duration: 2000, onClose })
  }
}

/**
 * 生成数据表格的函数
 *
 * @param r 数据引用对象
 * @returns
 */
export function useDataTable<E>(r: ListRef<E>) {
  return {
    async onQueryData<P, R extends PR<E>>(ds: DS<P, R>, p: P, page?: number): Promise<R | undefined> {
      if (page) {
        if (r.loading) return
        r.page = page
      }
      return queryData({ fn: ds, data: p, dataRef: r })
    },
    async onDelete<P, R extends { data: any }>(opt: {
      fn: DF<P, R>
      index?: number
      row?: E
      tips?: boolean
      onClose?: () => void
    }): Promise<Result<R>> {
      r.rowIndex = opt.index
      if (r.requesting) {
        throw "cancel"
      }
      var pkArr = []
      if (opt.row) {
        pkArr = [r.primary(opt.row)]
      } else {
        pkArr = (r.selectedRows || []).map(r.primary)
      }
      const ret = await deleteData({ fn: opt.fn, ids: pkArr, dataRef: r })
      // 触发提示
      opt.tips != false && alertTips(ret, opt.onClose)
      return ret
    },
    onSelectionChange(rows: Array<E>, row?: E) {
      onSelectionChange(r, rows, row)
    },
    onRevertSelection() {
      revertSelectionState(r)
    },
    onResetSelection() {
      resetSelection(r)
    },
    showModal: null
  }
}

export interface FormData<E> {
  // 数据表引用
  formRef: () => Ref
  // 请求状态
  requesting?: boolean
  // 主键
  pk?: number
  // 数据
  data: E
}

/** 数据表单 */
export function useDataForm<E>(f: FormData<E>) {
  return {
    async onSubmitForm<R extends { data: any }>(opt: {
      tips?: boolean
      updateFn: UF<E, R>
      createFn: EF<E, R>
      onClose?: () => void
    }): Promise<Result<E>> {
      return new Promise(async (resolve, reject) => {
        const valid = await (f.formRef().value as any)?.validate()
        if (!valid || f.requesting) {
          reject("")
        }
        f.requesting = true
        const { data } = await (f.pk ? opt.updateFn(f.pk, f.data) : opt.createFn(f.data)).finally(
          () => (f.requesting = false)
        )
        const ret = {
          ...data,
          ...parseResult(data),
          ...f.data
        }
        // 触发提示
        opt.tips != false && alertTips(ret, opt.onClose)
        resolve(ret)
      })
    }
  }
}
