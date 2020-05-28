#!lang:ts＃!name:表单界面
#!target:ts/feature/{{.table.Prefix}}/{{.table.Name}}/form.vue
<template>
  <div class="createPost-container">
    <el-form ref="formData" :model="formData" :rules="rules" class="form-container mod-form">
      <div class="createPost-main-container mod-form-container">
        {{range $i,$c := .columns}}
        {{if ne $c.IsPk true}}{{$name:= lower_title $c.Prop}}
        <el-row>
          <el-col :span="24">
            <el-form-item class="mod-form-item" label-width="80px" label="{{$c.Comment}}:"　prop="{{$name}}">
              <el-input v-model="formData.{{$name}}" class="mod-form-input" autosize　placeholder="请输入内容"/>
            </el-form-item>
          </el-col>
        </el-row>
        {{end}}{{end}}
        <sticky :z-index="10" class="mod-form-bar">
          <el-button v-loading="loading" style="margin-left: 80px;" type="success" @click="submitForm">
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
  @Prop({ default: false }) private isEdit!: boolean;
  @Prop({ default: 0 }) private id!: number|string;

  private validateRequire = (rule: any, value: string, callback: Function) => {
    if (value === '') {
      callback(new Error(rule.field + '为必填字段'))
    } else {
      callback()
    }
  }


  private formData :I{{$Class}} = default{{$Class}}();
  private loading = false;
  // 设置验证表单字段的规则
  private rules = {
    name: [{ validator: this.validateRequire }]
  };

  get lang() {
    return AppModule.language;
  }

  created() {
    if(!this.id && this.$route.params)this.id = this.$route.params.id;
    if (this.id) {
      this.fetchData(this.id);
    }
  }

  private async fetchData(id: any) {
    try {
      const { data } = await get{{$Class}}(id, { /* Your params here */ });
      this.formData = data;
      document.title = (this.lang === 'zh' ? '编辑{{$Comment}}' : 'Edit {{$Class}}')+'-'+id;
    } catch (err) {
      console.error(err);
    }
  }

  private submitForm() {
    (this.$refs.formData as Form).validate(async valid => {
      if (valid) {
        this.loading = true;
        let ret = await (this.id?update{{$Class}}(this.id,this.formData): create{{$Class}}(this.formData));
        const {errCode,errMsg} = ret.data;
        if(errCode === 0){
          this.$notify({
            title: '操作成功',
            message: '操作成功',
            type: 'success',
            duration: 2000
          });
          this.$emit("refresh",{state:1});
        }else{
          this.$notify({
            title: '操作失败',
            message: errMsg,
            type: 'error',
            duration: 2000
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
    padding: 40px 45px 20px 50px;

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
