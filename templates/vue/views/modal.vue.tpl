#!lang:ts＃!name:表单界面
#!target:vue/views/{{.table.Title}}/modal.vue
<template>
  <div class="mod-form-container">
    <el-form ref="formRef" class="mod-form" size="small"
            label-position="left" :label-width="85"
            :model="form.data" :rules="rules">
        {{range $i,$c := exclude .columns "create_time" "update_time"}}\
        {{if not $c.IsPk}}{{$name:=lower_title $c.Prop}}{{$ele:= $c.Render.Element}}\
            <el-form-item class="mod-form-item" label="{{$c.Comment}}:" prop="{{$name}}">
            {{/*<el-col :span="12">LEFT...</el-col><el-col :span="12">RIGHT...</el-col>*/}}
            {{if eq $ele "radio"}}\
              <el-switch v-model="form.data.{{$name}}"
                         active-text=""
                         inactive-text=""
                         :active-value="1"
                         :inactive-value="0">
              </el-switch>
            {{else if eq $ele "checkbox"}}\
                <el-checkbox v-model="form.data.{{$name}}"></el-checkbox>
            {{else if eq $ele "textarea"}}\
                <el-input type="textarea" v-model="form.data.{{$name}}" class="mod-form-input" :autosize="{ minRows: 2, maxRows: 4}" placeholder="请输入{{$c.Comment}}"/>
            {{else if eq $ele "select"}}\
                <el-select v-model="form.data.{{$name}}">
                   <el-option v-for="(value,attr) in {"选项1":1,"选项2":2}" :label="attr" :value="value"/>
                </el-select>
            {{else if ends_with $c.Name "_time"}}\
                <el-date-picker v-model="form.data.{{$name}}" type="date" value-format="timestamp" class="mod-form-input"
                    placeholder="选择日期" format="yyyy 年 MM 月 dd 日" >
                </el-date-picker>
            {{else if num_type $c.Type }}\
                <el-input v-model.number="form.data.{{$name}}" class="mod-form-input" autosize placeholder="请输入{{$c.Comment}}"/>
            {{else}}\
                <el-input v-model="form.data.{{$name}}" class="mod-form-input" autosize placeholder="请输入{{$c.Comment}}"/>
            {{end}}
                <span class="mod-form-remark"></span>
         </el-form-item>
         {{end}}{{end}}
         {{/* 在模态框中显示按钮,默认使用外部的按钮
         <el-container class="mod-form-bar">
            <el-button @click="()=>applyCallback()">取消</el-button>
            <el-button type="primary" :loading="form.requesting" @click="submitForm">提交</el-button>
          </el-container>
          */}}
     </el-form>
  </div>
</template>
{{$Class := .table.Title}}{{$Comment := .table.Comment}}
{{$validateColumns := exclude .columns .table.Pk "create_time" "update_time" "state"}}

<script lang="ts" setup>
import {ref,reactive} from "vue";
import {{`{`}}{{$Class}},get{{$Class}},create{{$Class}},update{{$Class}} } from "@/api"
import {Message,MessageBox,parseResult} from "@/ext/utils";
import { useDataForm, FormData } from "@/ext/compose"
/** #! 定义属性,接收父组件的参数 */
const props = withDefaults(defineProps<{modelValue?:{{type "ts" .table.PkType}}}>(),{});

const formRef = ref();
const form = reactive<FormData<{{$Class}}>>({
  formRef: () => formRef,
  pk: props.modelValue,
  data: new {{$Class}}()
})
const { onSubmitForm } = useDataForm(form)

/** #! 验证规则会反应到组件,比如required,所以不用在组件上再加required */
// 设置验证表单字段的规则,取消验证请注释对应的规则
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

const fetchData = async ()=>{
  {{/** 作为单独的页面
  //if(!form.pk){
  //  const {{`{`}}{{.table.Pk}}{{`}`}} = useRoute().currentRoute.query;
  //  if(id) form.pk = {{if num_type .table.PkType }}parseInt({{.table.Pk}} as string){{else}}{{.table.Pk}}{{end}};
  //}
  */}} \
  if(form.pk){
    try {
       form.data = (await get{{$Class}}(form.pk, { /* Your params here */ })).data;
    } catch (err) {
      console.error(err);
      Message.warning("数据加载失败:"+err)
    }
  }
}


// 取消
const cancel = (resolve:any)=>resolve()

// 提交表单
const submitForm = async (resolve: any, reject: any) => {
  const ret = await onSubmitForm({
    updateFn: update{{$Class}},
    createFn: create{{$Class}}
  })
  ret.errCode === 0 ? resolve(ret) : reject(ret.errMsg)
}

defineExpose({
  submit:submitForm,
  cancel: cancel
});

fetchData()
</script>

