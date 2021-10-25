package bond

import (
	"code.byted.org/tcg/bond/stock"
	"code.byted.org/tcg/bond/util"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

type Result struct {
	Page int    `json:"page"`
	Rows []*Row `json:"rows"`
}

type Row struct {
	ID   string `json:"id"`
	Cell *Bond  `json:"cell"`
}

type Bond struct {
	BondID           string  `json:"bond_id"`
	BondNM           string  `json:"bond_nm"`
	BondPY           string  `json:"bond_py"`
	Price            float64 `json:"price"`
	IncreaseRT       float64 `json:"increase_rt"`
	StockID          string  `json:"stock_id"`
	StockNM          string  `json:"stock_nm"`
	StockPY          string  `json:"stock_py"`
	StockPrice       float64 `json:"sprice"`
	StockIncreaseRT  float64 `json:"sincrease_rt"`
	ConvertPrice     float64 `json:"convert_price"`      // 转股价
	ConvertValue     float64 `json:"convert_value"`      // 转股价值
	PremiumRT        float64 `json:"premium_rt"`         // 转股溢价率
	RatingCD         string  `json:"rating_cd"`          // 评级
	PutConvertPrice  float64 `json:"put_convert_price"`  // 回售触发价
	ForceRedeemPrice float64 `json:"force_redeem_price"` // 强赎回触发价
	YearLeft         float64 `json:"year_left"`          // 剩余年限
	RemainMoney      float64 `json:"curr_iss_amt"`       // 剩余规模
	YtmRt            float64 `json:"ytm_rt"`             // 到期税前收益
	RedeemDt         *string `json:"redeem_dt"`          // 强赎回触发日期，为nil表示没有触发强赎

	VolatilityRate float64 // 波动率
	RatingScore    float64 // 计算排名得分
}

type Bonds []*Bond

func (bonds *Bonds) Print() {
	for _, b := range *bonds {
		fmt.Println(fmt.Sprintf("代码：%v 名称：%v 现价：%v 评级：%v 正股波动率：%v 排名得分：%v", b.BondID, b.BondNM, b.Price, b.RatingCD, b.VolatilityRate, b.RatingScore))
	}
}

func (bonds *Bonds) FilterRating() *Bonds {
	var filterBonds Bonds
	for _, b := range *bonds {
		if !Rating(b.RatingCD).IsLessThan("AA") {
			filterBonds = append(filterBonds, b)
		}
	}
	return &filterBonds
}

func (bonds *Bonds) FilterMoney() *Bonds {
	var filterBonds Bonds
	for _, b := range *bonds {
		if b.RemainMoney >= 0.5 {
			filterBonds = append(filterBonds, b)
		}
	}
	return &filterBonds
}

func (bonds *Bonds) FilterRedeem() *Bonds {
	var filterBonds Bonds
	for _, b := range *bonds {
		if b.RedeemDt == nil {
			filterBonds = append(filterBonds, b)
		}
	}
	return &filterBonds
}

func (bonds *Bonds) Sort() *Bonds {
	for _, b := range *bonds {
		history, err := stock.GetHistoryPrices(&stock.HistoryCondition{
			StockCode: b.StockID,
			StartTime: "20181024",
			EndTime:   "20211024",
		})
		if err != nil {
			panic(fmt.Sprintf("sort bonds error: %v", err))
		}
		b.VolatilityRate = history.CalYearVolatility()
	}
	bonds.calRatingScore()
	sort.Slice(*bonds, func(i, j int) bool {
		return (*bonds)[i].RatingScore > (*bonds)[j].RatingScore
	})
	return bonds
}

func (bonds *Bonds) calRatingScore() {
	sort.Slice(*bonds, func(i, j int) bool {
		return (*bonds)[i].VolatilityRate < (*bonds)[j].VolatilityRate
	})
	for i, b := range *bonds {
		b.RatingScore += float64(i)
	}
	sort.Slice(*bonds, func(i, j int) bool {
		return (*bonds)[i].PremiumRT > (*bonds)[j].PremiumRT
	})
	for i, b := range *bonds {
		b.RatingScore += float64(i)
	}
	sort.Slice(*bonds, func(i, j int) bool {
		return (*bonds)[i].YearLeft > (*bonds)[j].YearLeft
	})
	for i, b := range *bonds {
		b.RatingScore += float64(i)
	}
}

type FilterCondition struct {
	ToPrice float64 `json:"to_price"`
	YtmRt   float64 `json:"ytm_rt"`
}

func GetBonds(condition *FilterCondition) (Bonds, error) {
	url := fmt.Sprintf("https://www.jisilu.cn/data/cbnew/cb_list_new/?___jsl=LST___t=%v", time.Now().Unix()*1000)
	method := "POST"
	payload := fmt.Sprintf(`fprice=&tprice=%v&curr_iss_amt=&volume=&svolume=&premium_rt=&ytm_rt=%v&rating_cd=&is_search=Y`, condition.ToPrice, condition.YtmRt) + "&market_cd%5B%5D=shmb&market_cd%5B%5D=shkc&market_cd%5B%5D=szmb&market_cd%5B%5D=szcy&btype=C&listed=Y&qflag=N&sw_cd=&bond_ids=&rp=50&page=1"
	header := http.Header{}
	header.Add("Connection", "keep-alive")
	header.Add("sec-ch-ua", "\"Chromium\";v=\"94\", \"Google Chrome\";v=\"94\", \";Not A Brand\";v=\"99\"")
	header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	header.Add("X-Requested-With", "XMLHttpRequest")
	header.Add("sec-ch-ua-mobile", "?0")
	header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")
	header.Add("sec-ch-ua-platform", "\"Windows\"")
	header.Add("Origin", "https://www.jisilu.cn")
	header.Add("Sec-Fetch-Site", "same-origin")
	header.Add("Sec-Fetch-Mode", "cors")
	header.Add("Sec-Fetch-Dest", "empty")
	header.Add("Referer", "https://www.jisilu.cn/data/cbnew/")
	header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	header.Add("Cookie", "kbzw__Session=lpog3dvlfbag03440rduhg5hj5; kbz_newcookie=1; Hm_lvt_164fe01b1433a19b507595a43bf58262=1635061476,1635061597; kbzw_r_uname=bqxtt; kbzw__user_login=7Obd08_P1ebax9aX2dPu1euYrqXR0dTn8OTb3crUjabErduqrZKjxaet2pyrztnFp5WpqtXcxaiYpq3Xz7HNmJ2j1uDb0dWMoZWmrqyhso2yj8ui1dSexdDqyuDl1piumqeCnrjg5dfn2OOBws2Vmqmap52WuODlqayckNmqrZ6Jutznztu43Nm-4dWflqewo5yvjJ-tvrXEw5-YzdnM2Zm8ztzX5ouWpN_p4uXGn5qop6WXraKnmKSZqJfG2cfR092oqpywmqqY; Hm_lpvt_164fe01b1433a19b507595a43bf58262=1635065447")

	res, err := util.Client.Request(method, url, header, payload)
	if err != nil {
		return nil, err
	}
	var result Result
	err = json.Unmarshal([]byte(res), &result)
	if err != nil {
		return nil, err
	}
	var bonds Bonds
	for _, row := range result.Rows {
		bonds = append(bonds, row.Cell)
	}
	return bonds, nil
}
