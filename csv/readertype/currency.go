package readertype

import (
	"strings"
	"unicode"
)

type Currency Type

const (
	DefaultCurrency             Currency = "Default"
	Gil                                  = "Gil"
	GrandCompanySeal                     = "Grand Company Seal"
	PoeticTomestone                      = "Allagan Tomestones of Poetics"
	UncappedTomestone                    = "Allagan Tomestones of Causality"
	CappedTomestone                      = "Allagan Tomestones of Comedy"
	WolfMark                             = "Wolf Mark"
	AlliedSeal                           = "Allied Seal"
	MandervilleGoldSaucerPoints          = "MGP"
	CenturioSeal                         = "Centurio Seal"
	SackOfNuts                           = "Sack of Nut"
	WhiteCraftersScrip                   = "White Crafters' Scrip"
	PurpleCraftersScrip                  = "Purple Crafters' Scrip"
	WhiteGatherersScrip                  = "White Gatherers' Scrip"
	PurpleGatherersScrip                 = "Purple Gatherers' Scrip"
	SkybuildersScrip                     = "Skybuilders' Scrip"
	BicolorGemstone                      = "Bicolor Gemstone"
	BozjanCluster                        = "Bozjan Cluster"
	FauxLeaf                             = "Faux Leaf"
	SteelAmaljok                         = "Steel Amalj'ok"
	SylphicGoldleaf                      = "Sylphic Goldleaf"
	TitanCobaltpiece                     = "Titan Cobaltpiece"
	RainbowtidePsashp                    = "Rainbowtide Psashp"
	IxaliOaknot                          = "Ixali Oaknot"
	VanuWhitebone                        = "Vanu Whitebone"
	BlackCopperGil                       = "Black Copper Gil"
	CarvedKupoNut                        = "Carved Kupo Nut"
	KojinSango                           = "Kojin Sangos"
	AnantaDreamstaff                     = "Ananta Dreamstaff"
	NamazuKoban                          = "Namazu Koban"
	FaeFancy                             = "Fae Fancy"
	QitariCompliment                     = "Qitari Compliment"
	HammeredFrogment                     = "Hammered Frogment"
	ArkasodaraPana                       = "Arkasodara Pana"
	OmicronOmnitoken                     = "Omicron Omnitoken"
	LoporritCarat                        = "Loporrit Carat"
)

func (c Currency) String() string {
	return string(c)
}

func (c Currency) GetPlural() string {
	switch c {
	case Gil:
		return "Gil"
	case GrandCompanySeal:
		return "Grand Company Seals"
	case PoeticTomestone:
		return "Allagan Tomestones of Poetics"
	case UncappedTomestone:
		return "Allagan Tomestones of Causality"
	case CappedTomestone:
		return "Allagan Tomestones of Comedy"
	case WolfMark:
		return "Wolf Marks"
	case AlliedSeal:
		return "Allied Seals"
	case MandervilleGoldSaucerPoints:
		return "MGP"
	case CenturioSeal:
		return "Centurio Seals"
	case SackOfNuts:
		return "Sack of Nuts"
	case WhiteCraftersScrip:
		return "White Crafters' Scrips"
	case PurpleCraftersScrip:
		return "Purple Crafters' Scrips"
	case WhiteGatherersScrip:
		return "White Gatherers' Scrips"
	case PurpleGatherersScrip:
		return "Purple Gatherers' Scrips"
	case SkybuildersScrip:
		return "Skybuilders' Scrips"
	case BicolorGemstone:
		return "Bicolor Gemstones"
	case BozjanCluster:
		return "Bozjan Clusters"
	case FauxLeaf:
		return "Faux Leaves"
	case SteelAmaljok:
		return "Steel Amalj'oks"
	case SylphicGoldleaf:
		return "Sylphic Goldleaves"
	case TitanCobaltpiece:
		return "Titan Cobaltpieces"
	case RainbowtidePsashp:
		return "Rainbowtide Psashp"
	case IxaliOaknot:
		return "Ixali Oaknots"
	case VanuWhitebone:
		return "Vanu Whitebones"
	case BlackCopperGil:
		return "Black Copper Gil"
	case CarvedKupoNut:
		return "Carved Kupo Nuts"
	case KojinSango:
		return "Kojin Sangos"
	case AnantaDreamstaff:
		return "Ananta Dreamstaffs"
	case NamazuKoban:
		return "Namazu Kobans"
	case FaeFancy:
		return "Fae Fancies"
	case QitariCompliment:
		return "Qitari Compliments"
	case HammeredFrogment:
		return "Hammered Frogments"
	case ArkasodaraPana:
		return "Arkasodara Panas"
	case OmicronOmnitoken:
		return "Omicron Omnitokens"
	case LoporritCarat:
		return "Loporrit Carats"
	}

	return ""
}

func (c Currency) GetEffort() float64 {
	switch c {
	case Gil:
		return 0.85
	case GrandCompanySeal:
		return 0.9
	case PoeticTomestone:
		return 0.875
	case UncappedTomestone:
		return 1.1
	case CappedTomestone:
		return 1.15
	case WolfMark:
		return 1.05
	case AlliedSeal:
		return 1.05
	case MandervilleGoldSaucerPoints:
		return 0.7
	case CenturioSeal:
		return 1.03
	case SackOfNuts:
		return 1.01
	case WhiteCraftersScrip:
		return 1.08
	case PurpleCraftersScrip:
		return 1.1
	case WhiteGatherersScrip:
		return 1.1
	case PurpleGatherersScrip:
		return 1.12
	case SkybuildersScrip:
		return 1.25
	case BicolorGemstone:
		return 1.75
	case BozjanCluster:
		return 2.25
	case FauxLeaf:
		return 2.75
	case SteelAmaljok:
		return 0.95
	case SylphicGoldleaf:
		return 0.95
	case TitanCobaltpiece:
		return 0.95
	case RainbowtidePsashp:
		return 0.95
	case IxaliOaknot:
		return 0.95
	case VanuWhitebone:
		return 0.95
	case BlackCopperGil:
		return 0.95
	case CarvedKupoNut:
		return 0.95
	case KojinSango:
		return 0.95
	case AnantaDreamstaff:
		return 0.95
	case NamazuKoban:
		return 0.95
	case FaeFancy:
		return 0.95
	case QitariCompliment:
		return 0.95
	case HammeredFrogment:
		return 0.95
	case ArkasodaraPana:
		return 0.95
	case OmicronOmnitoken:
		return 0.95
	case LoporritCarat:
		return 0.95
	case DefaultCurrency:
		return 1.0
	}

	return 1.0
}

func FromApiParam(s string) Currency {
	// Strip out non alpha characters, then convert to lower case
	stripped := strings.Map(
		func(r rune) rune {
			if unicode.IsLetter(r) {
				return r
			}
			return -1
		}, s,
	)

	switch strings.ToLower(stripped) {
	case "gil":
		return Gil
	case "grandcompanyseals":
		return GrandCompanySeal
	case "gcseals":
		return GrandCompanySeal
	case "gcseal":
		return GrandCompanySeal
	case "allagantomestonesofpoetics":
		return PoeticTomestone
	case "poetics":
		return PoeticTomestone
	case "uncappedtome":
		return UncappedTomestone
	case "allagantomestonesofcausality":
		return UncappedTomestone
	case "causality":
		return UncappedTomestone
	case "cappedtome":
		return CappedTomestone
	case "allagantomestonesofcomedy":
		return CappedTomestone
	case "comedy":
		return CappedTomestone
	case "wolfmark":
		return WolfMark
	case "alliedseals":
		return AlliedSeal
	case "mgp":
		return MandervilleGoldSaucerPoints
	case "centurioseals":
		return CenturioSeal
	case "sackofnuts":
		return SackOfNuts
	case "whitecraftersscrip":
		return WhiteCraftersScrip
	case "whitecrafter":
		return WhiteCraftersScrip
	case "purplecrafter":
		return PurpleCraftersScrip
	case "purplecraftersscrip":
		return PurpleCraftersScrip
	case "whitegatherersscrip":
		return WhiteGatherersScrip
	case "purplegatherersscrip":
		return PurpleGatherersScrip
	case "whitegatherer":
		return WhiteGatherersScrip
	case "purplegatherer":
		return PurpleGatherersScrip
	case "skybuildersscrip":
		return SkybuildersScrip
	case "bicolorgemstone":
		return BicolorGemstone
	case "bozjancluster":
		return BozjanCluster
	case "fauxleaf":
		return FauxLeaf
	case "steelamaljok":
		return SteelAmaljok
	case "sylphicgoldleaf":
		return SylphicGoldleaf
	case "titancobaltpiece":
		return TitanCobaltpiece
	case "rainbowtidepsashp":
		return RainbowtidePsashp
	case "ixalioaknot":
		return IxaliOaknot
	case "vanuwhitebone":
		return VanuWhitebone
	case "blackcoppergil":
		return BlackCopperGil
	case "carvedkuponut":
		return CarvedKupoNut
	case "kojinsango":
		return KojinSango
	case "anantadreamstaff":
		return AnantaDreamstaff
	case "namazukoban":
		return NamazuKoban
	case "faefancy":
		return FaeFancy
	case "qitaricompliment":
		return QitariCompliment
	case "hammeredfrogment":
		return HammeredFrogment
	case "arkasodarapana":
		return ArkasodaraPana
	case "omicronomnitoken":
		return OmicronOmnitoken
	case "loporritcarat":
		return LoporritCarat
	default:
		return DefaultCurrency
	}
}

var itemIdToCurrency = map[int]Currency{
	25:    WolfMark,
	27:    AlliedSeal,
	28:    PoeticTomestone,
	29:    MandervilleGoldSaucerPoints,
	44:    UncappedTomestone,
	45:    CappedTomestone,
	10307: CenturioSeal,
	26533: SackOfNuts,
	28063: SkybuildersScrip,
	25199: WhiteCraftersScrip,
	25200: WhiteGatherersScrip,
	33913: PurpleCraftersScrip,
	33914: PurpleGatherersScrip,
	26807: BicolorGemstone,
	31135: BozjanCluster,
	30341: FauxLeaf,
	21075: SylphicGoldleaf,
	21076: SteelAmaljok,
	21078: TitanCobaltpiece,
	21077: RainbowtidePsashp,
	21073: IxaliOaknot,
	21074: VanuWhitebone,
	21079: BlackCopperGil,
	21080: CarvedKupoNut,
	21081: KojinSango,
	21935: AnantaDreamstaff,
	22525: NamazuKoban,
	28186: FaeFancy,
	28187: QitariCompliment,
	28188: HammeredFrogment,
	36657: ArkasodaraPana,
	37854: OmicronOmnitoken,
	38952: LoporritCarat,
}

var currencyToItemId = map[Currency]int{}

func init() {
	for itemId, currency := range itemIdToCurrency {
		currencyToItemId[currency] = itemId
	}
}

func FromItemId(itemId int) Currency {
	if currency, ok := itemIdToCurrency[itemId]; ok {
		return currency
	}

	return DefaultCurrency
}

func ToItemId(c Currency) int {
	if itemId, ok := currencyToItemId[c]; ok {
		return itemId
	}

	return 0
}
