#!lang:ts＃!name:表单界面
#!target:vue2/{{name_path .table.Name}}/modal.vue
<template>
  <div class="mod-form-container">
    <el-form ref="formRef" class="mod-form" size="small"
             label-position="right" :model="formData" :rules="rules">
        <el-row :gutter="15">
        {{range $i,$c := exclude .columns "create_time" "update_time"}}\
        {{if not $c.IsPk}}{{$name:= $c.Prop}}{{$ele:= $c.Render.Element}}\
          <el-col :md="12" :xs="24">
            <el-form-item class="mod-form-item" label-width="85px" label="{{$c.Comment}}:" prop="{{$name}}">
            {{if eq $ele "radio"}}\
              <el-switch v-model="formData.{{$name}}"
                         active-text=""
                         inactive-text=""
                         :active-value="1"
                         :inactive-value="0">
              </el-switch>
            {{else if eq $ele "checkbox"}}\
                <el-checkbox v-model="formData.{{$name}}"></el-checkbox>
            {{else if eq $ele "textarea"}}\
                <el-input type="textarea" v-model="formData.{{$name}}" class="mod-form-input" :autosize="{ minRows: 2, maxRows: 4}" placeholder=""/>
            {{else if eq $ele "select"}}\
                <el-select v-model="formData.{{$name}}">
                   <el-option v-for="(value,attr) in {"选项1":1,"选项2":2}" :label="attr" :value="value"/>
                </el-select>
            {{else if ends_with $c.Name "_time"}}\
                <el-date-picker v-model="formData.{{$name}}" type="date" value-format="timestamp" class="mod-form-input"
                    placeholder="选择日期" format="yyyy 年 MM 月 dd 日" >
                </el-date-picker>
            {{else if equal_any $c.Type 3 4 5}}\
                <el-input v-model.number="formData.{{$name}}" class="mod-form-input" autosize placeholder="请输入{{$c.Comment}}"/>
            {{else}}\
                <el-input v-model="formData.{{$name}}" class="mod-form-input" autosize placeholder="请输入{{$c.Comment}}"/>
            {{end}}
                <span class="mod-form-remark"></span>
            </el-form-item>
          </el-col>
        {{end}}{{end}}
        </el-row>
        <el-container class="mod-form-bar">
            <el-button @click="()=>callback({close: true})">取消</el-button>
            <el-button v-loading="requesting" type="primary" @click="submitForm">提交</el-button>
        </el-container>
     </el-form>
  </div>
</template>
{{$Class := .table.Title}}{{$Comment := .table.Comment}}
{{$validateColumns := exclude .columns .table.Pk "create_time" "update_time" "state"}}

<script setup>
import {onMounted, reactive, ref} from "vue";
import {Message,MessageBox} from "element-ui"
import {{`{`}}{{$Class}},get{{$Class}},create{{$Class}},update{{$Class}} } from "./api"
import {parseResult} from "@/hook";

/** #! 定义属性,接收父组件的参数 */
const props = defineProps({
  value:{ type: {{title (type "ts" .table.PkType)}}, default: {{default "ts" .table.PkType}}}
})

/** #! 定义Emit向父组件传递数据 */
const emit = defineEmits(['callback']);

const formRef = ref(null); /* #! formRef关联表单ref属性 */
let formData  = reactive(new {{$Class}}()); /** #! 表单数据 */
let requesting = ref(false);

// 设置验证表单字段的规则,取消验证请注释对应的规则
/** #! 验证规则会反应到组件,比如required,所以不用在组件上再加required */
let rules = {
  // 自定义验证规则：
  // phone: [{label:"phone", validator: this.validate }]
  // private async validate(rule: any, value: string, callback: Function){
  //   const label = rule.label || rule.field;
  //   if (value === '') { callback(new Error(label + '为必填字段'))} else {callback()}
  // } \
  {{range $i,$c := $validateColumns}}{{if ne $c.IsPk true}}
  {{if equal_any $c.Type 3 4 5}}\
  {{$c.Prop}}: [{required: true, message:"{{$c.Comment}}不能为空"}, \
      {type:"number", message:"{{$c.Comment}}必须为数字值"}], \
  {{else if $c.NotNull}}\
  {{$c.Prop}}: [{required: true, message:"{{$c.Comment}}不能为空"}], \
  {{end}}\
  {{end}}{{end}}
};

onMounted(()=>{
  {{/*if(props.value == null)props.value = this.$route.query.id;*/}} 
  if(props.value)fetchFormData(props.value);
})

const callback = (arg)=>emit("callback",arg);

const fetchFormData = async(id) =>{
  try {
    const { data } = await get{{$Class}}(id, { /* Your params here */ });
    formData = data;
  } catch (err) {
    console.error(err);
    Message.warning("数据加载失败:"+err.message)
  }
}

const submitForm = ()=> {
  formRef.value.validate(async valid => {
    if (valid) {
      if(requesting.value)return;requesting.value = true;
      let ret = await (props.value?update{{$Class}}(props.value,formData):create{{$Class}}(formData))
        .catch((ex)=>MessageBox.alert(ex.message,"错误"))
        .finally(()=>requesting.value=false);
      const {errCode,errMsg} = parseResult(ret.data);
      if(errCode === 0){
        Message.success({message:'操作成功',duration:2000});
        callback({state:1,close:true,args:{}});
      }else{
        MessageBox.alert(errMsg,"错误");
      }
    } else {
      return false
    }
  })
}
</script>

