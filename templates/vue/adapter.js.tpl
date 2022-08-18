#!kind:2#!target:vue/adapter.js

/** 您可以编辑此文件,并进行适配 */

import axios from "axios";
export {Message,MessageBox} from "element-ui"

// parse multiple kind response to standard result like 
// {errCode:1,errMsg:"success"}
export function parseResult(data) {
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

export const request = axios.create({
    baseURL: '',
    timeout: 5000
})