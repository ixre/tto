#!kind:1

/**
该模板生成权限系统的SQL脚本

# 变量
VUE_PREFIX: 如:bz/
*/

/** #! INSERT INTO public.perm_res (name, res_type, pid, key, path, icon, permission, sort_num, is_external, is_hidden, create_time, component_path, cache_) */
/** #! VALUES ('注册表', 1, 1, 'registry', '/registry/index', 'hamburger', '', 0, 0, 0, 1607214871, '_features/registry/index', ''); */

{{$prefix := env "VUE_PREFIX"}}

{{range $k,$tables :=  .groups}}
{{if ne $k "perm"}}

/* ====== insert data : {{$k}} ===== */

INSERT INTO public.perm_res (name, res_type, pid, key, path, icon, permission, sort_num, is_external, is_hidden, create_time, component_path, cache_)
VALUES ('{{$k}}', 1, 0, '{{$k}}', '/{{$k}}/index', 'hamburger', '', 0, 0, 0, 1607214871, '', '');

{{range $i,$table := $tables}}
   /* {{$table.Comment}}前端页面 */
   INSERT INTO public.perm_res (name, res_type, pid, key, path, icon, permission, sort_num, is_external, is_hidden, create_time, component_path, cache_)
   VALUES ('{{$table.Comment}}', 1, (SELECT distinct id FROM perm_res where name='{{$k}}'),
    '{{replace (name_path $table.Name) "/" ":"}}', '{{$prefix}}{{name_path $table.Name}}/index', 'hamburger', '',
    0, 0, 0, 1607214871, 'features/{{$prefix}}{{name_path $table.Name}}/index', '');

   /* 新增{{$table.Comment}} */
   INSERT INTO public.perm_res (name, res_type, pid, key, path, icon, permission, sort_num, is_external, is_hidden, create_time, component_path, cache_)
   VALUES ('新增{{$table.Comment}}(接口)', 0, (SELECT distinct id FROM perm_res where name='{{$k}}'),
    '{{replace (name_path $table.Name) "/" ":"}}:create', '', 'create', '',0, 0, 0, 0, '', '');

    /* 更新{{$table.Comment}} */
    INSERT INTO public.perm_res (name, res_type, pid, key, path, icon, permission, sort_num, is_external, is_hidden, create_time, component_path, cache_)
    VALUES ('更新{{$table.Comment}}(接口)', 0, (SELECT distinct id FROM perm_res where name='{{$k}}'),
        '{{replace (name_path $table.Name) "/" ":"}}:update', '', 'update', '',0, 0, 0, 0, '', '');

    /* 查询{{$table.Comment}} */
    INSERT INTO public.perm_res (name, res_type, pid, key, path, icon, permission, sort_num, is_external, is_hidden, create_time, component_path, cache_)
    VALUES ('查询{{$table.Comment}}(接口)', 0, (SELECT distinct id FROM perm_res where name='{{$k}}'),
        '{{replace (name_path $table.Name) "/" ":"}}:get', '', 'update', '',0, 0, 0, 0, '', '');

    /* 查询{{$table.Comment}} */
    INSERT INTO public.perm_res (name, res_type, pid, key, path, icon, permission, sort_num, is_external, is_hidden, create_time, component_path, cache_)
    VALUES ('{{$table.Comment}}列表(接口)', 0, (SELECT distinct id FROM perm_res where name='{{$k}}'),
        '{{replace (name_path $table.Name) "/" ":"}}:list', '', 'update', '',0, 0, 0, 0, '', '');

    /* 删除{{$table.Comment}} */
    INSERT INTO public.perm_res (name, res_type, pid, key, path, icon, permission, sort_num, is_external, is_hidden, create_time, component_path, cache_)
    VALUES ('删除{{$table.Comment}}(接口)', 0, (SELECT distinct id FROM perm_res where name='{{$k}}'),
        '{{replace (name_path $table.Name) "/" ":"}}:delete', '', 'update', '',0, 0, 0, 0, '', '');

{{end}}
{{end}}
{{end}}