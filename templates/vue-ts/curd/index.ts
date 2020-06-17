#!target:ts/feature/{{name_path .table.Name}}/index.ts
import Layout from "@/layout/index.vue";
import {RouteConfig} from "vue-router";

{{$path := ""}}

/** 如果新增编辑在新的窗口打开, 去掉create/edit/list 的hidden属性 */
export const {{.table.Title}}Routes : RouteConfig = {
    path: '/{{.table.Prefix}}/',
    component: Layout,
    redirect: '{{$path}}/index',
    meta: {
        //roles: ["admin"],
        title: '{{.table.Comment}}',
        icon: 'example'
    },
    children: [
        {
            path: '{{$path}}/index',
            component: () => import(/* webpackChunkName: "{{.table.Name}}-all" */ './all.vue'),
            name: '{{.table.Title}}List',
            meta: {
                //roles: ["admin"],
                title: '{{.table.Comment}}',
                icon: 'list',
            }
        },
        {
            path: '{{$path}}/create',
            component: () => import(/* webpackChunkName: "{{.table.Name}}-create" */ './create.vue'),
            name: 'Create{{.table.Title}}',
            meta: {
                //roles: ["admin"],
                title: '创建{{.table.Comment}}',
                icon: 'edit'
            }
        },
        {
            path: '{{$path}}/edit/:id',
            component: () => import(/* webpackChunkName: "{{.table.Name}}-edit" */ './edit.vue'),
            name: 'Edit{{.table.Title}}',
            meta: {
                //roles: ["admin"],
                title: '修改{{.table.Comment}}',
                noCache: true,
                activeMenu: '{{.table.Name}}/list',
                hidden: true
            }
        },
        {
            path: '{{$path}}/list',
            component: () => import(/* webpackChunkName: "{{.table.Name}}-list" */ './list.vue'),
            name: '{{.table.Title}}List',
            meta: {
                //roles: ["admin"],
                title: '{{.table.Comment}}列表',
                icon: 'list',
                hidden: true
            }
        }
    ]
};