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
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("router run error: %v", err)
	}
}

func GetBonds(c *gin.Context) {
	bonds, err := bond.GetBonds(&bond.FilterCondition{
		ToPrice: 105,
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
	c.HTML(http.StatusOK, "bond.html", gin.H{"title": "可转债筛选", "bonds": *filterBonds})
	//c.JSON(http.StatusOK, filterBonds.ToString())
}
