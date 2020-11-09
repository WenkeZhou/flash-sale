package errcode

var (
	ErrorBuyStock           = NewError(20010001, "购买商品失败")
	NotFoundStock           = NewError(20010002, "该商品不存在")
	SellOutStock            = NewError(20010003, "商品卖完了")
	ErrorPessimisticLock    = NewError(20010004, "悲观锁并发更新数据失败")
	ErrorOptimisticLock     = NewError(20010005, "乐观锁并发更新数据失败")
	ErrorTooManyRequest     = NewError(20010006, "太多并发请求")
	RedisGetVerifyHashError = NewError(20010007, "redis 获取值失败")
	VerifyHashNotEqual      = NewError(20010008, "verifyHash 不相等")
	UserRequestFrequently   = NewError(20010009, "用户访问太频繁")
)
