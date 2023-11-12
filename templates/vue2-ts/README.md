# Vue2+TypeScript代码模板

*此模板仅支持vue2, vue3版本请使用模板: vue*

模板设置接口前缀,在配置文件中设置`base_path`
```
[global]
base_path ="/qkto/admin"
```
生成适用于Java程序的代码,指定参数:`-lang java`
```
tto -lang java -conf app.conf
```

# 依赖文件
文件: fx/index.ts
```
// standard error
export interface ErrResult{
    errCode:number
    errMsg:string
}

// parse multiple kind response to standard result
export function parseResult(data:any):ErrResult {
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
```
