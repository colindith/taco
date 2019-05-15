package main

import "fmt"

type Saiyan struct {
	Name  string
	Power int
}

func super(saiyan *Saiyan) {

	fmt.Println("------", fmt.Sprintf("%T", saiyan))

	saiyan.Power++
	fmt.Println(saiyan)
}

func main() {
	gokuPointer := &Saiyan{
		Name:  "Goku",
		Power: 9000,
	}
	fmt.Println(fmt.Sprintf("%T", gokuPointer))
	super(gokuPointer)
	fmt.Println(gokuPointer.Power)

}
