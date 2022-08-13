/// t_BOMM_GoodsMst 货品基本信息表
type BOMMGoodsMst struct{
    // 商品编号
    FGoodsID int `db:"fGoodsID" pk:"yes"`
    // 商品类型
    FGoodsType string `db:"fGoodsType"`
    // 商品编码
    FGoodsCode string `db:"fGoodsCode"`
    // 商品名称
    FGoodsName string `db:"fGoodsName"`
    // 版本号
    FEdition float64 `db:"fEdition"`
    // 货品别名
    FAlias string `db:"fAlias"`
    // 货品英文名称
    FEName string `db:"fEName"`
    // 规格描述
    FSizeDesc string `db:"fSizeDesc"`
    // 品牌
    FBrandCode string `db:"fBrandCode"`
    // 系列代号
    FKitCode string `db:"fKitCode"`
    // 款式
    FStyleCode string `db:"fStyleCode"`
    // 颜色
    FClrCode string `db:"fClrCode"`
    // 标准单位
    FStdUnit string `db:"fStdUnit"`
    // 交易单位
    FBusinessUnit string `db:"fBusinessUnit"`
    // 产品类型
    FFgType string `db:"fFgType"`
    // 客户代号
    FCCode string `db:"fCCode"`
    // 是否按订单管理库存
    FIfStkByOrd string `db:"fIfStkByOrd"`
    // 是否按供应商管理库存
    FIfStkByPrv string `db:"fIfStkByPrv"`
    // 是否需要进行预留
    FIfReserve string `db:"fIfReserve"`
    // 重量单位
    FWeightUnitCode string `db:"fWeightUnitCode"`
    // 是否分包
    FIfSubpackage string `db:"fIfSubpackage"`
    // 包装件数
    FPackPcs int `db:"fPackPcs"`
    // 包装箱数
    FPackCtn int `db:"fPackCtn"`
    // 箱体规格单位
    FGoodsOUnit string `db:"fGoodsOUnit"`
    // 外箱长度
    FCtnL float64 `db:"fCtnL"`
    // 外箱宽度
    FCtnW float64 `db:"fCtnW"`
    // 外箱高度
    FCtnH float64 `db:"fCtnH"`
    // 外箱体积系数
    FCubeRate float64 `db:"fCubeRate"`
    // 外箱体积(CUFT)
    FOutCuft float64 `db:"fOutCuft"`
    // 外箱体积(立方米)
    FOutM3 float64 `db:"fOutM3"`
    // 主要生产厂别
    FMkCode string `db:"fMkCode"`
    // 主要存储仓库
    FStkCode string `db:"fStkCode"`
    // 海关料号
    FCustomGoodsCode string `db:"fCustomGoodsCode"`
    // 海关料名
    FCustomGoodsName string `db:"fCustomGoodsName"`
    // 法定海关单位
    FCustomUnit string `db:"fCustomUnit"`
    // 安全库存量
    FSafeQty float64 `db:"fSafeQty"`
    // 最低安全库存量
    FMinSafeQty float64 `db:"fMinSafeQty"`
    // 储位代号
    FPlaceCode string `db:"fPlaceCode"`
    // 销售属性
    FSaleType string `db:"fSaleType"`
    // 总体积
    FAllOutM3 float64 `db:"fAllOutM3"`
    // 总材积
    FAllOutCuft float64 `db:"fAllOutCuft"`
    // 总毛重
    FAllGW float64 `db:"fAllGW"`
    // 总净重
    FAllNW float64 `db:"fAllNW"`
    // 是否启用
    FIfUse bool `db:"fIfUse"`
}