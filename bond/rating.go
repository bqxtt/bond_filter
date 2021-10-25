package bond

type Rating string

var RatingValue = map[Rating]int{
	"AAA":  10,
	"AA+":  9,
	"AA":   8,
	"AA-":  7,
	"A+":   6,
	"A":    5,
	"A-":   4,
	"BBB+": 3,
	"BBB":  2,
	"BBB-": 1,
	"BB":   0,
}

func (rating Rating) IsLessThan(compareRating Rating) bool {
	return RatingValue[rating] < RatingValue[compareRating]
}
