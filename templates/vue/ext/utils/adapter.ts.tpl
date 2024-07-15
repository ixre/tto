#!kind:2#!target:vue/ext/utils/adapter.ts

/** 您可以编辑此文件,并进行适配,需要在utils包中导出 */

import axios from "axios";
export {Message,MessageBox} from "element-ui"
// vue3 + elementui-plus
// export {ElMessage as Message,ElMessageBox as MessageBox} from "element-plus";

// vue2只支持vue-router v3, 需要导入全局的router变量 
// import route from "@/router"
// export const router = route
import {router,route} from "vue-router"
export const router = useRouter();
export const route = useRoute();

// vue3 vue-router
//export {useRouter, useRoute} from "vue-router"


export const request = axios.create({
    baseURL: '',
    timeout: 5000
})




// parse multiple kind response to standard result like 
// {errCode:1,errMsg:"success"}
export function parseResult(data:any) {
    const {error} = data;
    if (error) return {errCode: 1, errMsg: error}
    const {message} = data;
    if (message) return {errCode: 1, errMsg: error}
    const {ErrCode, ErrMsg} = data;
    if (ErrCode || ErrMsg) return {errCode: ErrCode, errMsg: ErrMsg}
    const {errCode, errMsg} = data;
    if (errCode || errMsg) return {errCode: errCode, errMsg: errMsg}
    return {
        errCode: 0,
        errMsg: "success"
    }
}

export const formatColTime = (_,_1,v)=>parseTime(v);
export const formatColDate = (_,_1,v)=>parseDate(v);

export const parseDate = (time)=>parseTime(time,"{y}-{m}-{d}")

// Parse the time to string
export const parseTime = (time, format) => {
    if (!time)return "-";
    const fmt = format || '{y}-{m}-{d} {h}:{i}'
    let date = time;
    if (typeof time === 'string') {
        if(/^[0-9]+$/.test(time)){
            time = parseInt(time)
        }
        date = new Date(time)
    }
    if (typeof time === 'number') {
        date = new Date(time.toString().length === 10?time * 1000:time)
    }
    const formatObj = {
        y: date.getFullYear(),
        m: date.getMonth() + 1,
        d: date.getDate(),
        h: date.getHours(),
        i: date.getMinutes(),
        s: date.getSeconds(),
        a: date.getDay()
    }
    const timeStr = fmt.replace(/{([ymdhisa])+}/g, (result, key) => {
        const value = formatObj[key]
        // Note: getDay() returns 0 on Sunday
        if (key === 'a') return ['日', '一', '二', '三', '四', '五', '六'][value]
        return value.toString().padStart(2, '0')
    })
    return timeStr
}
