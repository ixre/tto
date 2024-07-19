#!lang:ts＃!name:表单界面
#!target:vue-pro/views/{{.table.Title}}/modal.vue
<template>
  <div class="mod-form-container">
    <el-form ref="formRef" class="mod-form"
            label-position="left" :label-width="85"
            :model="formData" :rules="rules">
      <ExtFormElement v-for="(ele, i) in columns" :key="i" v-model="formData[ele.prop]" :config="ele" />
      {{/* 在模态框中显示按钮,默认使用外部的按钮
         <el-container class="mod-form-bar">
            <el-button @click="()=>applyCallback()">取消</el-button>
            <el-button type="primary" :loading="loading" @click="submitForm">提交</el-button>
          </el-container>
          */}}
     </el-form>
  </div>
</template>
{{$Class := .table.Title}}{{$Comment := .table.Comment}}
{{$validateColumns := exclude .columns .table.Pk "create_time" "update_time" "state"}}

<script lang="ts" setup>
/**
 * {{.table.Comment}}表单组件
 * @author: {{.global.user}}
 * @description: {{.global.time}}
 */
import {{`{`}}{{$Class}},get{{$Class}},create{{$Class}},update{{$Class}} } from "@/api"
import {ElementProps, ExtFormElement, useDataForm } from "@/ext/compose"

const modelValue = defineModel()
const { formData, formRef, loading, buildQueryData, onSubmitForm } = useDataForm(modelValue,new {{$Class}}())

/**
 * 初始化数据
 */
const fetchData = buildQueryData(getSysGeneralOption)
fetchData()

// 定义数据列, 如果不指定,则默认使用ElementPlus列的定义方式
const columns: Array<ElementProps> = [
    {{range $i,$c := exclude .columns "create_time" "update_time"}}\
    {{if not $c.IsPk}}{{$name:=lower_title $c.Prop}}{{$ele:= $c.Render.Element}}\
        {{if eq $ele "radio"}}\
        {type: "switch",prop: "{{$name}}",label:"{{$c.Comment}}",remark:"", activeValue:1, inactiveValue:0},
        {{else if eq $ele "checkbox"}}\
        {type: "checkbox",prop: "{{$name}}",label:"{{$c.Comment}}",options:[{label:"选项一",value:0,label:"选项二",value:1}]},
        {{else if eq $ele "textarea"}}\
        {type: "textarea",prop: "{{$name}}",label:"{{$c.Comment}}",placeholder:"请输入{{$c.Comment}}",autosize:{ minRows: 2, maxRows: 4 }},
        {{else if eq $ele "select"}}\
        {type: "select",prop: "{{$name}}",label:"{{$c.Comment}}",options: [{label:"选项一",value:0,label:"选项二",value:1}]},
        {{else if ends_with $c.Name "_date"}}\
        {type: "date", prop: "{{$name}}",label:"{{$c.Comment}}",placeholder:"选择日期",valueFormat:"timestamp",format:"yyyy 年 MM 月 dd 日"},
        {{else if ends_with $c.Name "_time"}}\
        {type: "datetime", prop: "{{$name}}",label:"{{$c.Comment}}",placeholder:"选择时间",valueFormat:"timestamp",format:"yyyy 年 MM 月 dd 日 HH:mm"},
        {{else if num_type $c.Type }}\
        {type: "number", prop: "{{$name}}",label:"{{$c.Comment}}",placeholder:"请输入{{$c.Comment}}"},
        {{else}}\
        {type: "input", prop: "{{$name}}",label:"{{$c.Comment}}",placeholder:"请输入{{$c.Comment}}"},
        {{end}}
      {{end}}{{end}}
]

/**
 * 设置验证表单字段的规则,取消验证请注释对应的规则
 */
const rules = {
  // 自定义验证规则：
  // phone: [{label:"phone", validator: this.validate }]
  // private async validate(rule: any, value: string, callback: Function){
  //   const label = rule.label || rule.field;
  //   if (value === '') { callback(new Error(label + '为必填字段'))} else {callback()}
  // } \
  {{range $i,$c := $validateColumns}}{{if ne $c.IsPk true}}
  {{if num_type $c.Type }}\
  {{lower_title $c.Prop}}: [{required: true, message:"{{$c.Comment}}不能为空"}, \
      {type:"number", message:"{{$c.Comment}}必须为数字"}], \
  {{else if $c.NotNull}}\
  {{lower_title $c.Prop}}: [{required: true, message:"{{$c.Comment}}不能为空"}], \
  {{end}}\
  {{end}}{{end}}
};

/**
 * 提交表单
 * @param resolve 执行成功
 * @param reject 执行失败
 */
const submitForm = async (resolve: any, reject: any) => {
  await onSubmitForm({
    updateFn: update{{$Class}},
    createFn: create{{$Class}}
  },    
  { resolve, reject }
  )
}

/**
 * 导出函数
 */
defineExpose({
  submit:submitForm,
  cancel: (resolve:any)=>resolve(),
  loading,
})
</script>

