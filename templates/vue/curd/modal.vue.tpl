#!lang:ts＃!name:表单界面
#!target:vue/{{name_path .table.Name}}/modal.vue
<template>
  <div class="mod-form-container">
    <el-form ref="formRef" class="mod-form" size="small"
             label-position="right" :model="formData" :rules="rules">
        {{range $i,$c := exclude .columns "create_time" "update_time"}}\
        {{if not $c.IsPk}}{{$name:= $c.Prop}}{{$ele:= $c.Render.Element}}\
            <el-form-item class="mod-form-item" label-position="left" label-width="85px" label="{{$c.Comment}}:" prop="{{$name}}">
            {{/*<el-col :span="12">LEFT...</el-col><el-col :span="12">RIGHT...</el-col>*/}}
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
            {{else if num_type $c.Type }}\
                <el-input v-model.number="formData.{{$name}}" class="mod-form-input" autosize placeholder="请输入{{$c.Comment}}"/>
            {{else}}\
                <el-input v-model="formData.{{$name}}" class="mod-form-input" autosize placeholder="请输入{{$c.Comment}}"/>
            {{end}}
                <span class="mod-form-remark"></span>
         </el-form-item>
         {{end}}{{end}}
         <el-container class="mod-form-bar">
            <el-button @click="()=>cancelForm({close: true})">取消</el-button>
            <el-button type="primary" :loading="requesting" @click="submitForm">提交</el-button>
        </el-container>
     </el-form>
  </div>
</template>
{{$Class := .table.Title}}{{$Comment := .table.Comment}}
{{$validateColumns := exclude .columns .table.Pk "create_time" "update_time" "state"}}

<script setup>
import {onMounted,ref} from "vue";
import {{`{`}}{{$Class}},default{{$Class}},get{{$Class}},create{{$Class}},update{{$Class}} } from "./api"
import {Message,MessageBox,router,parseResult} from "@/utils/adapter";

/** #! 定义属性,接收父组件的参数 */
const props = defineProps({
  value:{ type: {{title (type "ts" .table.PkType)}}, default: {{default "ts" .table.PkType}}}
})

/** #! 定义Emit向父组件传递数据 */
const emit = defineEmits(['callback']);

const formRef = ref(null); /* #! formRef关联表单ref属性 */
let formData  = ref(default{{$Class}}()); /** #! 表单数据 */
let requesting = ref(false);
let isModal = props.value;

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
  {{if num_type $c.Type }}\
  {{$c.Prop}}: [{required: true, message:"{{$c.Comment}}不能为空"}, \
      {type:"number", message:"{{$c.Comment}}必须为数字"}], \
  {{else if $c.NotNull}}\
  {{$c.Prop}}: [{required: true, message:"{{$c.Comment}}不能为空"}], \
  {{end}}\
  {{end}}{{end}}
};

onMounted(()=>{
  const {{`{`}}{{.table.Pk}}{{`}`}} = router.currentRoute.query;
  if(id)props.value = {{if num_type .table.PkType }}parseInt({{.table.Pk}}){{else}}{{.table.Pk}}{{end}};
  if(props.value)fetchFormData(props.value);
})

const cancelForm = (arg)=> (isModal && emit("callback",arg)) || router.go(-1)

const fetchFormData = async(id) =>{
  try {
    const { data } = await get{{$Class}}(id, { /* Your params here */ });
    formData.value = data;
  } catch (err) {
    console.error(err);
    Message.warning("数据加载失败:"+err.message)
  }
}

const submitForm = ()=> {
  formRef.value.validate(async valid => {
    if (!valid)return;
    if(requesting.value)return;requesting.value = true;
    let ret = await (props.value?update{{$Class}}(props.value,formData.value):create{{$Class}}(formData.value))
      .finally(()=>requesting.value=false);
    const {errCode,errMsg} = parseResult(ret.data);
    if(errCode === 0){
      Message.success({message:'操作成功',duration:2000});
      callback({state:1,close:true});
    }else{
      MessageBox.alert(errMsg,"错误");
    }
  })
}
</script>

