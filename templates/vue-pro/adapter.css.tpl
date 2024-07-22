#!kind:2#!target:vue-pro/adapter.css

/* 需要手动引用到全局样式 */

.mod-grid-container{align-content:flex-start}
.mod-grid-header{
  display:flex;flex-flow:row unwrap;max-width: 100%;
}
.mod-grid-footer{height:45px;overflow:hidden;text-align:right;padding:2px;}
.mod-grid-body{max-width: 100%;}

/** 解决el-select不显示的问题 */
.filter-select{width:240px}