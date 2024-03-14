package exchange

type GatheringPoint struct {
	Point
	Level      int
	GatherType string
	PointType  string
	Region     string
	Area       string
	Place      string
}

func (g *GatheringPoint) GetRegion() string {
	return g.Region
}

func (g *GatheringPoint) GetPlace() string {
	return g.Place
}

func (g *GatheringPoint) GetArea() string {
	return g.Area
}
