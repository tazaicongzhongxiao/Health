package app

import "fmt"

//
// GetRedisViewsStatisticsPlatformPvKey
// 获Platformr统计key
//
func GetRedisViewsStatisticsPlatformPvKey(coId, countTime int64, storeId int64, platform int) string {
	return fmt.Sprintf(RedisViewsStatisticsPlatformPvKey, coId, countTime, storeId, platform)
}

//
// GetRedisViewsStatisticsPlatformUvKey
// 获取Platformr统计key
//
func GetRedisViewsStatisticsPlatformUvKey(coId, countTime int64, storeId int64, platform int) string {
	return fmt.Sprintf(RedisViewsStatisticsPlatformUvKey, coId, countTime, storeId, platform)
}

//
// GetRedisViewsStatisticsAreaPvKey
// 获Area统计key
//
func GetRedisViewsStatisticsAreaPvKey(coId, countTime, storeId, pageId int64) string {
	return fmt.Sprintf(RedisViewsStatisticsAreaPvKey, coId, countTime, storeId, pageId)
}

//
// GetRedisViewsStatisticsAreaUvKey
// 获取Area统计key
//
func GetRedisViewsStatisticsAreaUvKey(coId, countTime, storeId, pageId int64) string {
	return fmt.Sprintf(RedisViewsStatisticsAreaUvKey, coId, countTime, storeId, pageId)
}

//
// GetRedisViewsStatisticsGoodsPvKey
// 获Goods统计key
//
func GetRedisViewsStatisticsGoodsPvKey(coId, countTime, storeId, goodsId int64) string {
	return fmt.Sprintf(RedisViewsStatisticsGoodsPvKey, coId, countTime, storeId, goodsId)
}

//
// GetRedisViewsStatisticsGoodsUvKey
// 获取Goods统计key
//
func GetRedisViewsStatisticsGoodsUvKey(coId, countTime, storeId, goodsId int64) string {
	return fmt.Sprintf(RedisViewsStatisticsGoodsUvKey, coId, countTime, storeId, goodsId)
}

//
// GetRedisViewsStatisticsViewPvKey
// 获View统计key
//
func GetRedisViewsStatisticsViewPvKey(coId, countTime int64) string {
	return fmt.Sprintf(RedisViewsStatisticsViewPvKey, coId, countTime)
}

//
// GetRedisViewsStatisticsViewUvKey
// 获取View统计key
//
func GetRedisViewsStatisticsViewUvKey(coId, countTime int64) string {
	return fmt.Sprintf(RedisViewsStatisticsViewUvKey, coId, countTime)
}
