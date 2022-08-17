<template>
  <div class="container">
    <div slot="header" class="mod-card-header">
      <span>平台设置</span>
    </div>
    <el-form :model="data.formData"
             :rules="dataRule"
             ref="dataForm"
             class="mod-form"
             label-width="140px">
      <el-form-item label="平台名称" prop="platform_name">
        <el-input v-model="data.formData.platform_name"></el-input>
      </el-form-item>
      <el-form-item label="LOGO" prop="platform_logo">
        <el-input v-model="data.formData.platform_logo"></el-input>
      </el-form-item>
      <el-form-item label="对比色LOGO" prop="platform_inverse_color_logo">
        <el-input v-model="data.formData.platform_inverse_color_logo"></el-input>
      </el-form-item>
      <el-form-item label="零售门户LOGO" prop="platform_retail_site_logo">
        <el-input v-model="data.formData.platform_retail_site_logo">
        </el-input>
      </el-form-item>

      <el-form-item label="批发门户LOGO" prop="platform_wholesale_site_logo">
        <el-input v-model="data.formData.platform_wholesale_site_logo">
        </el-input>
      </el-form-item>
      <el-form-item label="客服电话" prop="platform_service_tel">
        <el-input v-model="data.formData.platform_service_tel">
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="submitForm()" :loading="loading">保存</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script lang="ts">
import {onMounted, reactive, ref} from "@vue/composition-api";
import {queryPlatformConf, updatePlatformConf} from "@/views/feature/settings/system/api";
import {Form, Message} from "element-ui";
import {parseResult} from "@/fx";

const dataRule = {
  platform_name: {required: true, message: '平台名称不能为空'},
  platform_logo: {required: true, message: '未上传平台Logo'},
  // platform_inverse_color_logo:{required: true, message: '平台名称不能为空'},
  // platform_retail_site_logo:{required: true, message: '平台名称不能为空'},
  // platform_wholesale_site_logo:{required: true, message: '平台名称不能为空'},
  platform_service_tel: {required: true, message: '客服电话不能为空'},
}

export default {
  setup(props, context) {
    // 使用context.refs.dataForm访问,无须创建引用
    // 使用ref创建引用, 使用dataForm.value访问
    let dataForm = ref();
    let data = reactive({
      requesting: 0,
      formData: {
        platform_name: "",
        platform_logo: "",
        platform_inverse_color_logo: "",
        platform_retail_site_logo: "",
        platform_wholesale_site_logo: "",
        platform_service_tel: ""
      }
    });
    data.formData.platform_name = "hello"
    onMounted(async () => {
      const ret = await queryPlatformConf();
      console.log(ret.data);
      data.formData = ret.data;
    })

    const submitForm = () => {
      (dataForm.value as Form).validate(async valid => {
        console.log(data.formData.platform_name);
        if (valid) {
          if (data.requesting === 1) return
          data.requesting = 1;
          let ret = await (updatePlatformConf(data.formData)).finally(() => data.requesting = 0);
          const {errCode, errMsg} = parseResult(ret.data);
          if (errCode === 0) {
            Message({
              type: "success",
              message: '操作成功'
            })
            //this.callback({state:1})
          } else {
            Message({
              type:"error",
              message:errMsg
            })
          }
        } else {
          return false
        }
      })
    }

    return {
      data,
      dataRule: ref(dataRule),
      loading: false,
      submitForm,
      dataForm
    }
  }
}
</script>

<style scoped>
.mod-form {
  max-width: 50%;
}
</style>
