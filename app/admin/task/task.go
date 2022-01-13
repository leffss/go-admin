package task

import "fmt"

type Ab struct {
	Id int
	Name string
}

func Test(id int) *Ab {
	return &Ab{
		Id: id,
		Name: "leffss",
	}
}

func Test2(id int) string {
	return fmt.Sprintf("Admin Test Task: %d", id)
}