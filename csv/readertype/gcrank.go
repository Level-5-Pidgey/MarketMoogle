package readertype

type GrandCompanyRank int8

const (
	GrandCompanyRankNone GrandCompanyRank = iota
	PrivateThirdClass
	PrivateSecondClass
	PrivateFirstClass
	Corporal
	SergeantThirdClass
	SergeantSecondClass
	SergeantFirstClass
	ChiefSergeant
	SecondLieutenant
	FirstLieutenant
	Captain
)

func (r GrandCompanyRank) String() string {
	return [...]string{
		"None",
		"Private Third Class",
		"Private Second Class",
		"Private First Class",
		"Corporal",
		"Sergeant Third Class",
		"Sergeant Second Class",
		"Sergeant First Class",
		"Chief Sergeant",
		"Second Lieutenant",
		"First Lieutenant",
		"Captain",
	}[r]
}
