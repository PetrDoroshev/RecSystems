package main

import (
	"fmt"
)

type Assumption struct {
	Text string
	MD   float64
	NMD  float64
}

func (a *Assumption) And(rhs *Assumption) *Assumption {

	return &Assumption{Text: fmt.Sprintf("%s и %s", a.Text, rhs.Text),
		MD:  min(a.MD, rhs.MD),
		NMD: max(a.NMD, rhs.NMD)}
}

func (a *Assumption) Or(rhs *Assumption) *Assumption {

	return &Assumption{Text: fmt.Sprintf("%s или %s", a.Text, rhs.Text),
		MD:  max(a.MD, rhs.MD),
		NMD: min(a.NMD, rhs.NMD)}
}

func (a Assumption) String() string {
	return fmt.Sprintf("%-45s MD: %.2f, NMD: %.2f", a.Text, a.MD, a.NMD)
}

func MD(a1 Assumption, a2 Assumption) float64 {
	return a1.MD + a2.MD*(1-a1.MD)
}

func NMD(a1 Assumption, a2 Assumption) float64 {
	return a1.NMD + a2.NMD*(1-a1.NMD)
}

func main() {

	a1 := Assumption{Text: "X проживает в Y", MD: 0.8, NMD: 0.3}
	a2 := Assumption{Text: "X является членом партии Z", MD: 0.75, NMD: 0.25}

	a3 := Assumption{Text: "X имеет возраст T", MD: 0.4, NMD: 0.5}
	a4 := Assumption{Text: "X является ИП", MD: 0.85, NMD: 0.2}

	fmt.Println(a1.String())
	fmt.Println(a2.String())
	fmt.Println(a3.String())
	fmt.Println(a4.String())

	a1 = *a1.And(&a2)
	a2 = *a3.Or(&a4)

	fmt.Printf("\nE1: %s\nE2: %s\n", a1, a2)

	MD := MD(a1, a2)
	NMD := NMD(a1, a2)

	fmt.Printf("\nMD[H:E1, E2] = %f\n", MD)
	fmt.Printf("NMD[H:E1, E2] = %f\n", NMD)

	fmt.Printf("\nКоэффициент уверенности: \nKU[H:E] = %f", MD-NMD)

}
