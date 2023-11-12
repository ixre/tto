# Vue2模板-curd版

*此模板兼容`Vue2.7`和`vue3`*

模板设置接口前缀,在配置文件中设置`base_path`

```ini
[global]
base_path ="/app/admin"
```

生成适用于Java程序的代码,指定参数:`-lang java`

```shell
tto -lang java -conf app.conf
```

## 权限指令

使用权限指令代码如下:

```vue
<span v-perm="{ key: 'C02+1', roles: ['admin'], visible:true}" @click="handleCreate">
    <el-button class="filter-item" icon="el-icon-plus">创建页面</el-button>
</span>
```

指令包含3个属性:

- key  权限key, 编码规则如下文;
- roles  角色数组, 优先匹配角色,如果当前用户不存在于数组中,再使用key进行验证权限.
- visible 是否可见,可见状态弹出提示,　不可见状态直接隐藏组件

权限key的编码规则:

```txt
(模块[A-Z])+(子模块[01-99])+(页面[01-99])+(组件[01-99])
```

额外权限:

```txt
key添加添加权限值,实现对增(1)删(2)改(4)的控制,　如:"C010102+1"
```
