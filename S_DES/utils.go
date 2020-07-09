package S_DES

import (
	"log"
	"strconv"
	"strings"
)

func StringToIntArray(s string) []int {
	var res []int

	cache := strings.Split(s, " ")
	for i := 0; i < len(cache); i++ {
		num, err := strconv.Atoi(cache[i])
		if err != nil {
			log.Fatalln("Problem converting string to Int")
		}
		res = append(res, num)
	}

	return res
}

func StringTo2DIntArray(s string, rowSize int, colSize int) [][]int {
	var res [][]int
	res = make([][]int, rowSize)
	for i := 0; i < rowSize; i++ {
		res[i] = make([]int, colSize)
	}

	cache := strings.Split(s, " ")
	for i := 0; i < rowSize; i++ {
		for j := 0; j < colSize; j++ {
			num, err := strconv.Atoi(cache[i*colSize+j])
			if err != nil {
				log.Fatalln("Problem converting String to Int. String to matrix routine")
			}
			res[i][j] = num
		}
	}

	return res
}

func MinInt(a int, b int) int {
	if a <= b {
		return a
	}
	return b
}

func MaxInt(a int, b int) int {
	if a >= b {
		return a
	}
	return b
}

func Substr(s string, startPos int, length int) string {
	var res string
	for i := startPos; i < MinInt(startPos+length, len(s)); i++ {
		res = res + (string)(s[i])
	}

	return res
}
