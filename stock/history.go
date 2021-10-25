package stock

import (
	"encoding/json"
	"fmt"
	"github.com/bqxtt/bond_filter/util"
	"math"
	"strconv"
)

type Results []*Result

type Result struct {
	Status       int         `json:"status"`
	HistoryQuery []*DataList `json:"hq"`
	Code         string      `json:"code"`
}

type DataList []string

func (data *DataList) GetDate() string {
	return (*data)[0]
}

func (data *DataList) GetStartPrice() float64 {
	price, err := strconv.ParseFloat((*data)[1], 64)
	if err != nil {
		panic("get start price error")
	}
	return price
}

func (data *DataList) GetEndPrice() float64 {
	price, err := strconv.ParseFloat((*data)[2], 64)
	if err != nil {
		panic("get end price error")
	}
	return price
}

type HistoryPrices struct {
	DataListSet    []*DataList
	DayVolatility  float64
	YearVolatility float64
}

func (his *HistoryPrices) CalDayVolatility() float64 {
	var yields []float64
	for i := 0; i < len(his.DataListSet)-1; i++ {
		yields = append(yields, (his.DataListSet[i].GetEndPrice()/his.DataListSet[i+1].GetEndPrice())-1)
	}
	var squaredSum, sum float64
	n := float64(len(his.DataListSet))
	for _, yield := range yields {
		squaredSum += yield * yield
		sum += yield
	}
	//STDEV.S
	his.DayVolatility = math.Sqrt((n*squaredSum - sum*sum) / n / (n - 1))
	return his.DayVolatility
}

func (his *HistoryPrices) CalYearVolatility() float64 {
	his.CalDayVolatility()
	his.YearVolatility = his.DayVolatility * math.Sqrt(244)
	return his.YearVolatility
}

type HistoryCondition struct {
	StockCode string
	StartTime string
	EndTime   string
}

func GetHistoryPrices(condition *HistoryCondition) (*HistoryPrices, error) {
	url := fmt.Sprintf("http://q.stock.sohu.com/hisHq?code=cn_%v&start=%v&end=%v", condition.StockCode, condition.StartTime, condition.EndTime)
	method := "GET"
	res, err := util.Client.Request(method, url, nil, "")
	if err != nil {
		return nil, err
	}
	var results Results
	err = json.Unmarshal([]byte(res), &results)
	if err != nil {
		return nil, err
	}
	codeMap := map[string]*Result{}
	for _, result := range results {
		codeMap[result.Code] = result
	}
	if result, exist := codeMap[fmt.Sprintf("cn_%v", condition.StockCode)]; exist {
		return &HistoryPrices{DataListSet: result.HistoryQuery}, nil
	} else {
		return nil, fmt.Errorf("stock code: %v, not found", condition.StockCode)
	}
}
