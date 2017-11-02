package Common

type Node struct {
	Name string
	Address string
}

type PID struct {
	Name string
	Address string
}

type Config struct {
	Myself PID
	Coordinators []PID
}