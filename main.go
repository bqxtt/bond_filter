package main

import (
	"code.byted.org/tcg/bond/bond"
	"fmt"
)

func main() {
	bonds, err := bond.GetBonds(&bond.FilterCondition{
		ToPrice: 105,
		YtmRt:   0,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	bonds.
		FilterMoney().
		FilterRating().
		FilterRedeem().
		Sort().
		Print()
}
