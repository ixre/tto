#!lang:ts＃!name:新增界面
#!target:ts/feature/{{.table.Prefix}}/{{.table.Name}}/create.vue
<template>
  <{{replace .table.Name "_" "-"}}-form :is-edit="false" />
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import {{.table.Title}}Form from './form.vue';

@Component({
  name: 'Create{{.table.Title}}',
  components: {
    {{.table.Title}}Form
  }
})
export default class extends Vue {};
</script>
