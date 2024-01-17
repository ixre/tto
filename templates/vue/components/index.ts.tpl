#!kind:2#!target:vue/components/index.ts

import { Ref, nextTick } from "vue";
import { parseResult } from "../utils";
export * from "./ModalDialog";

// 列表引用
export interface ListDataRef<R> {
    // 返回唯一主键
    primary: (r: R) => any;
    // 数据表引用
    tableRef: () => Ref,
    // 总条数
    total?: number;
    // 页码
    page: number;
    // 数量
    size: number;
    // 分页数据
    dataList?: Array<R>;
    // 当前页码
    rowIndex?: number,
    // 加载数据状态
    loading?: boolean;
    // 操作请求状态
    requesting?: boolean,
    // 已选择的行
    selectedRows: Array<R>,
}

/**
 * 查询数据列表
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
export async function queryDataList<E, P, R extends { data: { rows: Array<E>, total: number } }>(
    dataRef: ListDataRef<E>, fetch: (page: number, size: number, p: P) => Promise<R>, queryParams: P): Promise<R | undefined> {
    if (!dataRef.loading) {
        dataRef.loading = true;
        const ret = await fetch(dataRef.page, dataRef.size, queryParams)
            .finally(() => dataRef.loading = false);
        const { data } = ret;
        dataRef.dataList = data.rows;
        if (data.total > 0) dataRef.total = data.total;
        revertSelectionState(dataRef);
        return ret;
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
export async function deleteData<E, P, R extends { data: any }>(dataRef: ListDataRef<E>, deleteFn: (...p: Array<P>) => Promise<R>, ...p: Array<P>)
    : Promise<R & { errCode: number, errMsg: string }> {
    if (dataRef.requesting) throw new Error("requesting");
    dataRef.requesting = true;
    let ret = await deleteFn(...p)
        .finally(() => { dataRef.requesting = false; dataRef.rowIndex = 0; });
    return {
        ...ret.data,
        ...parseResult(ret.data),
    };
}


// 回显选中状态
async function revertSelectionState<E>(queryData: ListDataRef<E>) {
    await nextTick();
    const { tableRef, dataList, primary, selectedRows } = queryData;
    if (tableRef) {
        const { toggleRowSelection, clearSelection } = tableRef().value;
        if (selectedRows) {
            dataList?.filter(a => selectedRows.findIndex(b => {
                if (typeof (a) != typeof (b)) {
                    console.error(`selectedRows基础类型为${typeof (b)}, 无法自动选中行`,)
                }
                return primary(a) == primary(b);
            }) != -1).map(a => toggleRowSelection(a, true));
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
export function onSelectionChange<E>(queryData: ListDataRef<E>, rows: Array<E>, row?: E) {
    const { primary, selectedRows, dataList } = queryData;
    const fi = (r: E) => (queryData.selectedRows || [])?.findIndex(a => primary(a) == primary(r));
    if (row == null) {
        // 全选/取消全选
        if (rows.length > 0) {
            // 全选
            rows.forEach(a => {
                if (fi(a) == -1) {
                    selectedRows.push(a)
                }
            });
        } else {
            // 取消全选
            dataList?.forEach(a => {
                const i = fi(a);
                if (i != -1) {
                    selectedRows.splice(i, 1)
                }
            })
        }
    } else {
        if (rows.findIndex(a => primary(a) == primary(row)) >= 0) {
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
export function onResetSelection<E>(queryData: ListDataRef<E>) {
    queryData.selectedRows = [];
    queryData.tableRef().value.clearSelection();
}