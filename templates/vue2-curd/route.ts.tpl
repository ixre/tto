#!kind:2#!target:vue/{{.prefix}}/route.ts
{{$tables := .tables}}
{{$first := get_n .tables 0}}\
{{$first_path := substr_n $first.Name "_" 1}}\
/** 如果新增编辑在新的窗口打开, 去掉create/edit/list 的hidden属性 */
export const {{$first.Prefix}}Routes = {
  path: '/{{$first.Prefix}}/',
  redirect: '{{$first_path}}/index',
  meta: {
    // roles: ["admin"],
    title: '{{$first.Prefix}}',
    icon: 'example'
  },
  children: [\
  {{range $i,$table := .tables}}\
  {{$path := name_path $table.Name}}
    {
      path: '{{$path}}/index',
      name: '{{$table.Title}}Index',
      component: () => import(/* webpackChunkName: "{{$table.Name}}-index" */ './{{$path}}/index.vue'),
      meta: {
        // roles: ["admin"],
        title: '{{$table.Comment}}',
        icon: 'list'
      }
    },
    {
      path: '{{$path}}/detail',
      name: '{{$table.Title}}Detail',
      component: () => import(/* webpackChunkName: "{{$table.Name}}-detail" */ './{{$path}}/modal.vue'),
      meta: {
        // roles: ["admin"],
        title: '{{$table.Comment}}',
        hidden: true
      }
    } \
    {{if not (is_last $i $tables)}},{{end}}
  {{end}}\
  ]
}
