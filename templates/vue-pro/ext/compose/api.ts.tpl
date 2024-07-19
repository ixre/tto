#!kind:2#!target:vue-pro/ext/compose/api.ts

import { ref, Ref } from "vue"
import { ElMessage as Message } from "element-plus"
type ApiResult = any & { code: number; msg: string }
export function defineApiCompose() {
  const loadingRef: Ref<boolean> = ref(false)
  return {
    // 请求状态
    loading: loadingRef,
    // 请求方法
    async request(fn: () => Promise<any & ApiResult>): Promise<any & ApiResult> {
      if (loadingRef.value) {
        return Promise.reject("requesting..")
      }
      loadingRef.value = true
      const { data } = await fn().finally(() => ((loadingRef as any).value = false))
      const { code, msg } = data
      if (code > 0) {
        Message.error(msg)
        return null
      }
      return Promise.resolve(data)
    }
  }
}
