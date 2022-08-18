#!kind:2#!target:vue/{{.prefix}}/route.js
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
  {{$path := substr_n $table.Name "_" 1}}
    {
      path: '{{$path}}/index',
      name: '{{$table.Title}}Index',
      component: () => import(/* webpackChunkName: "{{$table.Name}}-index" */ './{{$path}}/index.vue'),
      meta: {
        // roles: ["admin"],
        title: '{{$table.Comment}}',
        icon: 'list'
      }
    }{{if not (is_last $i $tables)}},{{end}}
  {{end}}\
  ]
}
