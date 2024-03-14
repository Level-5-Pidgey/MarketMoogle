package exchange

type Point interface {
	GetRegion() string

	GetZone() string

	GetPlace() string

	GetArea() string
}
