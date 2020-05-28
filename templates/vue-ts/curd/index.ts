#!target:ts/feature/{{.table.Prefix}}/{{.table.Name}}/index.ts
import Layout from "@/layout/index.vue";
import {RouteConfig} from "vue-router";

/** 如果新增编辑在新的窗口打开, 去掉create/edit/list 的hidden属性 */
export const {{.table.Title}}Routes : RouteConfig = {
    path: '/{{.table.Name}}',
    component: Layout,
    redirect: '/{{.table.Name}}/list',
    meta: {
        title: '{{.table.Comment}}',
        icon: 'example'
    },
    children: [
        {
            path: 'create',
            component: () => import(/* webpackChunkName: "product-create" */ './create.vue'),
            name: 'Create{{.table.Title}}',
            meta: {
                title: '创建{{.table.Comment}}',
                icon: 'edit'
            }
        },
        {
            path: 'edit/:id',
            component: () => import(/* webpackChunkName: "product-edit" */ './edit.vue'),
            name: 'Edit{{.table.Title}}',
            meta: {
                title: '修改{{.table.Comment}}',
                noCache: true,
                activeMenu: '/{{.table.Name}}/list',
                hidden: true
            }
        },
        {
            path: 'index',
            component: () => import(/* webpackChunkName: "product-list" */ './all.vue'),
            name: '{{.table.Title}}List',
            meta: {
                title: '{{.table.Comment}}',
                icon: 'list'
            }
        },
        {
            path: 'list',
            component: () => import(/* webpackChunkName: "product-list" */ './list.vue'),
            name: '{{.table.Title}}List',
            meta: {
                title: '{{.table.Comment}}列表',
                icon: 'list',
                hidden: true
            }
        }
    ]
};