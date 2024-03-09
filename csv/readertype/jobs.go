package readertype

type Job string

const (
	JobNone          Job = "None"
	JobGladiator         = "GLA"
	JobPugilist          = "PGL"
	JobMarauder          = "MRD"
	JobLancer            = "LNC"
	JobArcanist          = "ARC"
	JobConjurer          = "CNJ"
	JobThaumaturge       = "THM"
	JobCarpenter         = "CRP"
	JobBlacksmith        = "BSM"
	JobArmourer          = "ARM"
	JobGoldsmith         = "GSM"
	JobLeatherworker     = "LTW"
	JobWeaver            = "WVR"
	JobAlchemist         = "ALC"
	JobCulinarian        = "CUL"
	JobMiner             = "MIN"
	JobBotanist          = "BTN"
	JobFisher            = "FSH"
	JobPaladin           = "PLD"
	JobMonk              = "MNK"
	JobWarrior           = "WAR"
	JobDragoon           = "DRG"
	JobBard              = "BRD"
	JobWhiteMage         = "WHM"
	JobBlackMage         = "BLM"
	JobSummoner          = "SMN"
	JobScholar           = "SCH"
	JobRogue             = "ROG"
	JobNinja             = "NIN"
	JobBlueMage          = "BLU"
	JobMachinist         = "MCH"
	JobDarkKnight        = "DRK"
	JobAstrologian       = "AST"
	JobSamurai           = "SAM"
	JobRedMage           = "RDM"
	JobGunbreaker        = "GNB"
	JobDancer            = "DNC"
	JobReaper            = "RPR"
	JobSage              = "SGE"
	JobViper             = "VPR"
	JobPictomancer       = "PCT"
)

func (j Job) FromShortString(s string) Job {
	switch s {
	case "GLA":
		return JobGladiator
	case "PGL":
		return JobPugilist
	case "MRD":
		return JobMarauder
	case "LNC":
		return JobLancer
	case "ARC":
		return JobArcanist
	case "ACN":
		return JobArcanist
	case "CNJ":
		return JobConjurer
	case "THM":
		return JobThaumaturge
	case "CRP":
		return JobCarpenter
	case "BSM":
		return JobBlacksmith
	case "ARM":
		return JobArmourer
	case "GSM":
		return JobGoldsmith
	case "LTW":
		return JobLeatherworker
	case "WVR":
		return JobWeaver
	case "ALC":
		return JobAlchemist
	case "CUL":
		return JobCulinarian
	case "MIN":
		return JobMiner
	case "BTN":
		return JobBotanist
	case "FSH":
		return JobFisher
	case "PLD":
		return JobPaladin
	case "MNK":
		return JobMonk
	case "WAR":
		return JobWarrior
	case "DRG":
		return JobDragoon
	case "BRD":
		return JobBard
	case "WHM":
		return JobWhiteMage
	case "BLM":
		return JobBlackMage
	case "SMN":
		return JobSummoner
	case "SCH":
		return JobScholar
	case "ROG":
		return JobRogue
	case "NIN":
		return JobNinja
	case "BLU":
		return JobBlueMage
	case "MCH":
		return JobMachinist
	case "DRK":
		return JobDarkKnight
	case "AST":
		return JobAstrologian
	case "SAM":
		return JobSamurai
	case "RDM":
		return JobRedMage
	case "GNB":
		return JobGunbreaker
	case "DNC":
		return JobDancer
	case "RPR":
		return JobReaper
	case "SGE":
		return JobSage
	case "VPR":
		return JobViper
	case "PCT":
		return JobPictomancer
	default:
		return JobNone
	}
}

func (j Job) String() string {
	switch j {
	case JobGladiator:
		return "Gladiator"
	case JobPugilist:
		return "Pugilist"
	case JobMarauder:
		return "Marauder"
	case JobLancer:
		return "Lancer"
	case JobArcanist:
		return "Arcanist"
	case JobConjurer:
		return "Conjurer"
	case JobThaumaturge:
		return "Thaumaturge"
	case JobCarpenter:
		return "Carpenter"
	case JobBlacksmith:
		return "Blacksmith"
	case JobArmourer:
		return "Armourer"
	case JobGoldsmith:
		return "Goldsmith"
	case JobLeatherworker:
		return "Leatherworker"
	case JobWeaver:
		return "Weaver"
	case JobAlchemist:
		return "Alchemist"
	case JobCulinarian:
		return "Culinarian"
	case JobMiner:
		return "Miner"
	case JobBotanist:
		return "Botanist"
	case JobFisher:
		return "Fisher"
	case JobPaladin:
		return "Paladin"
	case JobMonk:
		return "Monk"
	case JobWarrior:
		return "Warrior"
	case JobDragoon:
		return "Dragoon"
	case JobBard:
		return "Bard"
	case JobRogue:
		return "Rogue"
	case JobNinja:
		return "Ninja"
	case JobWhiteMage:
		return "WhiteMage"
	case JobBlackMage:
		return "BlackMage"
	case JobSummoner:
		return "Summoner"
	case JobScholar:
		return "Scholar"
	case JobMachinist:
		return "Machinist"
	case JobDarkKnight:
		return "DarkKnight"
	case JobAstrologian:
		return "Astrologian"
	case JobSamurai:
		return "Samurai"
	case JobRedMage:
		return "RedMage"
	case JobGunbreaker:
		return "Gunbreaker"
	case JobDancer:
		return "Dancer"
	case JobReaper:
		return "Reaper"
	case JobSage:
		return "Sage"
	case JobViper:
		return "Viper"
	case JobPictomancer:
		return "Pictomancer"
	default:
		return "None"
	}
}
