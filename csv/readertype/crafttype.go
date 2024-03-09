package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type CraftType struct {
	Key  int
	Name string
	Job  Job
}

func jobFromCraftType(key int) Job {
	switch key {
	case 0:
		return JobCarpenter
	case 1:
		return JobBlacksmith
	case 2:
		return JobArmourer
	case 3:
		return JobGoldsmith
	case 4:
		return JobLeatherworker
	case 5:
		return JobWeaver
	case 6:
		return JobAlchemist
	case 7:
		return JobCulinarian
	default:
		return JobNone
	}
}

func (c CraftType) CreateFromCsvRow(record []string) (*CraftType, error) {

	key := util.SafeStringToInt(record[0])
	return &CraftType{
		Key:  key,
		Name: record[3],
		Job:  jobFromCraftType(key),
	}, nil
}

func (c CraftType) GetKey() int {
	return c.Key
}
