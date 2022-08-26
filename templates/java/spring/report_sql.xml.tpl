#!target:query/{{.table.Prefix}}/{{substr_n .table.Name "_" 1}}_list.xml
<?xml version="1.0" encoding="utf-8" ?>
<ExportItemConfig>
    <ColumnMapping>
    {{range $i,$v := .columns}}{{if ne $i 0}};{{end}}{{$v.Name}}:{{$v.Comment}}{{end}}
    </ColumnMapping>
    <Query>
        <![CDATA[
        SELECT * FROM {{.table.Name}}
        WHERE {where}
        ORDER BY {order_by}
        LIMIT {{if eq .global.db "pgsql"}}{page_size} OFFSET {page_offset}{{else}}{page_offset},{page_size}{{end}}
     ]]>
    </Query>
    <Import><![CDATA[]]></Import>
    <Total>
        <![CDATA[
            SELECT COUNT(1) FROM {{.table.Name}} WHERE {where}
        ]]>
    </Total>
</ExportItemConfig>
