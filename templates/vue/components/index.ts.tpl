#!kind:2#!target:vue/components/index.ts
import { parseResult } from "../utils";
export * from "./ModalDialog";

// 列表引用
export interface ListDataRef<R> {
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
    selectedRows?: Array<R>,
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