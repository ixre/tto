#!lang:ts＃!name:表单界面
#!target:ts/feature/{{name_path .table.Name}}/form.vue
<template>
  <div class="createPost-container">
    <el-form ref="formData" :model="formData" :rules="rules" size="small" class="form-container mod-form">
      <div class="createPost-main-container mod-form-container">
        {{range $i,$c := .columns}}\
        {{if not $c.IsPk}}{{$name:= lower_title $c.Prop}}{{$ele:= $c.Render.Element}}\
        <el-row>
          <el-col :span="24">
            <el-form-item class="mod-form-item" label-width="78px" label="{{$c.Comment}}"　prop="{{$name}}">
            {{if eq $ele "radio"}}\
                <el-radio-group v-model="formData.{{$name}}">
                  <el-radio :label="1">是</el-radio>
                  <el-radio :label="0">否</el-radio>
                </el-radio-group>
            {{else if eq $ele "checkbox"}}\
                <el-checkbox v-model="formData.{{$name}}"></el-checkbox>
            {{else if eq $ele "textarea"}}\
                <el-input type="textarea" v-model="formData.{{$name}}" class="mod-form-input" :autosize="{ minRows: 2, maxRows: 4}"　placeholder="请输入内容"/>
            {{else if eq $ele "select"}}\
                <el-select v-model="formData.{{$name}}">
                   <el-option v-for="(value,attr) in {"选项1":1,"选项2":2}" :label="attr" :value="value"/>
                </el-select>
            {{else}}\
                <el-input v-model="formData.{{$name}}" class="mod-form-input" autosize　placeholder="请输入内容"/>
            {{end}}
            </el-form-item>
          </el-col>
        </el-row>
        {{end}}{{end}}
        <sticky :z-index="10" class="mod-form-bar">
          <el-button v-loading="requesting" style="margin-left: 80px;" type="success" @click="submitForm">
            提交
          </el-button>
        </sticky>
      </div>
    </el-form>
  </div>
</template>
{{$Class := .table.Title}}{{$Comment := .table.Comment}}
<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import { AppModule } from '@/store/modules/app'
import Sticky from '@/components/Sticky/index.vue'
import { Form } from 'element-ui'
import {I{{$Class}}, default{{$Class}},get{{$Class}},create{{$Class}},update{{$Class}} } from './api'

@Component({
  name: '{{$Class}}Form',
  components: {
    Sticky,
  }
})
export default class extends Vue {
  @Prop({ default: 0 }) private id!: number|string; /* 如果id有值,则为更新.　反之为新增.　如：:disabled="id" */

  private formData :I{{$Class}} = default{{$Class}}();
  private requesting = 0;

  {{$validateColumns := exclude .columns .table.Pk "create_time" "update_time" "state"}}
  private validate = (rule: any, value: string, callback: Function) => {
    const label = rule.label || rule.field;\
  {{range $i,$c := $validateColumns}}\
    {{if eq (type "ts" $c.Type) "number" }}
    if(rule.field === "{{lower_title $c.Prop}}" && !/^\d[\d\.]*$/.test(value))return callback(new Error(label+ '不正确'));\
    {{end}}{{end}}\
    if (value === '') {
      callback(new Error(label + '为必填字段'))
    } else {
      callback()
    }
  }

  // 设置验证表单字段的规则,取消验证请注释对应的规则
  private rules = {
    {{range $i,$c := $validateColumns}}{{if ne $c.IsPk true}}
    {{lower_title $c.Prop}}: [{label:"{{$c.Comment}}", validator: this.validate }] \
    {{if not (is_last $i $validateColumns)}},{{end}}{{end}}{{end}}
  };


  get lang() {
    return AppModule.language;
  }

  created() {
    if(!this.id && this.$route.params)this.id = this.$route.params.id;
    if (this.id)this.fetchData(this.id);
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

  private submitForm() {
    (this.$refs.formData as Form).validate(async valid => {
      if (valid) {
        if(this.requesting === 1)return;this.requesting = 1;
        let ret = await (this.id?update{{$Class}}(this.id,this.formData):create{{$Class}}(this.formData)).finally(()=>this.requesting = 0);
        const {errCode,errMsg} = ret.data;
        if(errCode === 0){
          this.$notify.success({
            title: '操作成功',
            message: '操作成功',
            duration:2000
          });
          this.$emit("refresh",{state:1});
        }else{
          this.$notify.error({
            title: '操作失败',
            message: errMsg,
            duration:2000
          });
          this.$emit("refresh",{state:0});
        }
      } else {
        return false
      }
    })
  }
}
</script>

<style lang="scss">
.article-textarea {
  textarea {
    padding-right: 40px;
    resize: none;
    border: none;
    border-radius: 0px;
    border-bottom: 1px solid #bfcbd9;
  }
}
</style>

<style lang="scss" scoped>
.createPost-container {
  position: relative;

  .createPost-main-container {
    padding:0 45px 20px 50px;

    .postInfo-container {
      position: relative;
      @include clearfix;
      margin-bottom: 10px;

      .postInfo-container-item {
        float: left;
      }
    }
  }

  .word-counter {
    width: 40px;
    position: absolute;
    right: 10px;
    top: 0px;
  }
}
</style>
