package readertype

import (
	"errors"
	"github.com/level-5-pidgey/MarketMoogle/util"
)

type ClassJobCategory struct {
	Id             int
	JobCategory    string
	JobsInCategory []Job
}

func jobFromCsvRow(rowNum int) Job {
	switch rowNum {
	case 3:
		return JobGladiator
	case 4:
		return JobPugilist
	case 5:
		return JobMarauder
	case 6:
		return JobLancer
	case 7:
		return JobArcanist
	case 8:
		return JobConjurer
	case 9:
		return JobThaumaturge
	case 10:
		return JobCarpenter
	case 11:
		return JobBlacksmith
	case 12:
		return JobArmourer
	case 13:
		return JobGoldsmith
	case 14:
		return JobLeatherworker
	case 15:
		return JobWeaver
	case 16:
		return JobAlchemist
	case 17:
		return JobCulinarian
	case 18:
		return JobMiner
	case 19:
		return JobBotanist
	case 20:
		return JobFisher
	case 21:
		return JobPaladin
	case 22:
		return JobMonk
	case 23:
		return JobWarrior
	case 24:
		return JobDragoon
	case 25:
		return JobBard
	case 26:
		return JobWhiteMage
	case 27:
		return JobBlackMage
	case 28:
		return JobArcanist
	case 29:
		return JobSummoner
	case 30:
		return JobScholar
	case 31:
		return JobRogue
	case 32:
		return JobNinja
	case 33:
		return JobMachinist
	case 34:
		return JobDarkKnight
	case 35:
		return JobAstrologian
	case 36:
		return JobSamurai
	case 37:
		return JobRedMage
	case 38:
		return JobBlueMage
	case 39:
		return JobGunbreaker
	case 40:
		return JobDancer
	case 41:
		return JobReaper
	case 42:
		return JobSage
	case 43:
		return JobViper
	case 44:
		return JobPictomancer
	default:
		return JobNone
	}
}

func (r ClassJobCategory) CreateFromCsvRow(record []string) (*ClassJobCategory, error) {
	categoryName := record[1]

	if categoryName == "" {
		return nil, errors.New("ClassJobCategory has no name")
	}

	jobs := make([]Job, 0, 41)

	for i, val := range record[2:] {
		if util.SafeStringToBool(val) {
			fromCsv := jobFromCsvRow(i + 2)

			if fromCsv == JobNone {
				continue
			}

			jobs = append(jobs, jobFromCsvRow(i+2))
		}
	}

	if len(jobs) == 0 {
		return nil, errors.New("ClassJobCategory has no jobs")
	}

	return &ClassJobCategory{
		Id:             util.SafeStringToInt(record[0]),
		JobCategory:    record[1],
		JobsInCategory: jobs,
	}, nil
}

func (r ClassJobCategory) GetKey() int {
	return r.Id
}
