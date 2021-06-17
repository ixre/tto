#!lang:ts＃!name:表单界面
#!target:ts/feature/{{name_path .table.Name}}/form.vue
<template>
  <div class="createPost-container">
    <el-form ref="formData" class="form-container mod-form" size="small"
             label-position="right" :model="formData" :rules="rules">
      <div class="createPost-main-container mod-form-container">
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
        <sticky :z-index="10" class="mod-form-bar">
            <el-button @click="()=>callback({close: true})">取消</el-button>
            <el-button v-loading="requesting" type="primary" @click="submitForm">提交</el-button>
        </sticky>
      </div>
    </el-form>
  </div>
</template>
{{$Class := .table.Title}}{{$Comment := .table.Comment}}
{{$validateColumns := exclude .columns .table.Pk "create_time" "update_time" "state"}}

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { AppModule } from '@/store/modules/app';
import Sticky from '@/components/Sticky/index.vue';
import { Form } from 'element-ui';
import {I{{$Class}}, default{{$Class}},get{{$Class}},create{{$Class}},update{{$Class}} } from './api'
import {parseResult} from "@/fx";

@Component({
  name: '{{$Class}}Form',
  components: {
    Sticky,
  }
})
export default class extends Vue {
  @Prop({ default: 0 }) private value!: number|string; /* 如果id有值,则为更新. 反之为新增. 如：:disabled="id" */

  private formData :I{{$Class}} = default{{$Class}}();
  private requesting = 0;

  // 设置验证表单字段的规则,取消验证请注释对应的规则
  /** #! 验证规则会反应到组件,比如required,所以不用在组件上再加required */
  private rules = {
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


  get lang() {
    return AppModule.language;
  }

  created() {
    if(!this.value && this.$route.params){
      this.value = this.$route.query.id as string;
    }
    if (this.value)this.fetchData(this.value);
  }

  private async fetchData(id: any) {
    try {
      /* document.title = (this.lang === 'zh' ? '编辑{{$Comment}}' : 'Edit {{$Class}}')+'-'+id; */
      const { data } = await get{{$Class}}(id, { /* Your params here */ });
      this.formData = data;
    } catch (err) {
      console.error(err);
    }
  }

  private callback(arg:any={state:0,close:true,args:{}}){
    this.$emit("callback",arg);
  }
  private submitForm() {
    (this.$refs.formData as Form).validate(async valid => {
      if (valid) {
        if(this.requesting === 1)return;this.requesting = 1;
        let ret = await (this.value?update{{$Class}}(this.value,this.formData):create{{$Class}}(this.formData)).finally(()=>this.requesting = 0);
        const {errCode,errMsg} = parseResult(ret.data);
        if(errCode === 0){
          this.$notify.success({
            title: '提示',
            message: '操作成功',
            duration:2000,
          });
          this.callback({state:1,close:true,args:{}})
        }else{
          await this.$alert(errMsg,"提示");
        }
      } else {
        return false
      }
    })
  }
}
</script>

<style lang="scss" scoped>
.createPost-container {
  position: relative;

  .createPost-main-container {
    padding:0 15px;

    .postInfo-container {
      position: relative;
      @include clearfix;
      margin-bottom: 10px;

      .postInfo-container-item {
        float: left;
      }
    }
  }
}
</style>
