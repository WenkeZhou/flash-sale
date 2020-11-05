package errcode

var (
	ErrorBuyStock = NewError(20010001, "购买商品失败")
	NotFoundStock = NewError(20010002, "该商品不存在")
	SellOutStock  = NewError(20010003, "商品卖完了")
)
