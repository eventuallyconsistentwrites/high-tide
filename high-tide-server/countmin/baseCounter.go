package countmin

type BaseCounter interface {
	String() string
	Update(string)
	PointQuery(string) int
	Reset()
}
