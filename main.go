package main

import (
	"fmt"
	"github.com/bqxtt/bond_filter/bond"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("html/*")
	r.GET("/bond", GetBonds)
	if err := r.Run(":1113"); err != nil {
		log.Fatalf("router run error: %v", err)
	}
}

type BondRequest struct {
	Price float64 `json:"price" form:"price"`
}

func GetBonds(c *gin.Context) {
	request := &BondRequest{}
	var price float64 = 105
	if err := c.ShouldBindQuery(request); err == nil {
		price = request.Price
	}
	bonds, err := bond.GetBonds(&bond.FilterCondition{
		ToPrice: price,
		YtmRt:   0,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	filterBonds := bonds.
		FilterMoney().
		FilterRating().
		FilterRedeem().
		Sort()
	c.HTML(http.StatusOK, "bond.html", gin.H{"title": "bqx的可转债筛选", "price": fmt.Sprint(price), "bonds": *filterBonds})
}
