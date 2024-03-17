package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type CollectableShopRewardScrip struct {
	Key        int
	Currency   int
	LowReward  int
	MidReward  int
	HighReward int
}

func (c CollectableShopRewardScrip) GetKey() int {
	return c.Key
}

func (c CollectableShopRewardScrip) CreateFromCsvRow(record []string) (*CollectableShopRewardScrip, error) {
	return &CollectableShopRewardScrip{
		Key:        util.SafeStringToInt(record[0]),
		Currency:   util.SafeStringToInt(record[1]),
		LowReward:  util.SafeStringToInt(record[2]),
		MidReward:  util.SafeStringToInt(record[3]),
		HighReward: util.SafeStringToInt(record[4]),
	}, nil
}
