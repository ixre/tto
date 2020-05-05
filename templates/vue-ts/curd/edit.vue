#!lang:ts＃!name:编辑界面
#!target:ts/feature/{{.table.Prefix}}/{{.table.Name}}/edit.vue
<template>
  <{{replace .table.Name "_" "-"}}-form :id="id"/>
</template>

<script lang="ts">
  import { Component, Vue,Prop } from 'vue-property-decorator';
  import {{.table.Title}}Form from './form.vue';

  @Component({
    name: 'Create{{.table.Title}}',
    components: {
      {{.table.Title}}Form
    }
  })
  export default class extends Vue {
    @Prop({ default: 0 }) private id!: number|string;
  };
</script>
