package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strconv"
)

type Input struct {
	Keys struct {
		N int `json:"n"`
		K int `json:"k"`
	} `json:"keys"`
	Data map[string]struct {
		Base  string `json:"base"`
		Value string `json:"value"`
	} `json:"data"`
}

func parseInput(jsonData string) (int, int, [][2]*big.Int) {
	var input Input
	err := json.Unmarshal([]byte(jsonData), &input)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return 0, 0, nil
	}

	n, k := input.Keys.N, input.Keys.K
	if n < k {
		fmt.Println("Error: Not enough roots provided.")
		return 0, 0, nil
	}

	var points [][2]*big.Int
	for key, val := range input.Data {
		base, err := strconv.Atoi(val.Base)
		if err != nil {
			fmt.Println("Error converting base for key:", key)
			continue
		}

		y := new(big.Int)
		y.SetString(val.Value, base)

		x := new(big.Int)
		x.SetString(key, 10)

		points = append(points, [2]*big.Int{x, y})
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i][0].Cmp(points[j][0]) < 0
	})

	return n, k, points[:k]
}

func gaussSolve(matrix [][]*big.Int, k int) []*big.Int {
	coeffs := make([]*big.Int, k)
	for i := range coeffs {
		coeffs[i] = new(big.Int)
	}

	for i := 0; i < k; i++ {
		for j := i + 1; j < k; j++ {
			factor := new(big.Int).Div(matrix[j][i], matrix[i][i])
			for l := i; l <= k; l++ {
				matrix[j][l].Sub(matrix[j][l], new(big.Int).Mul(factor, matrix[i][l]))
			}
		}
	}

	for i := k - 1; i >= 0; i-- {
		coeffs[i].Set(matrix[i][k])
		for j := i + 1; j < k; j++ {
			coeffs[i].Sub(coeffs[i], new(big.Int).Mul(matrix[i][j], coeffs[j]))
		}
		coeffs[i].Div(coeffs[i], matrix[i][i])
	}

	return coeffs
}

func main() {
	jsonData := `{
		"keys": {
			"n": 4,
			"k": 3
		},
		"data": {
			"1": { "base": "10", "value": "4" },
			"2": { "base": "2", "value": "111" },
			"3": { "base": "10", "value": "12" },
			"6": { "base": "4", "value": "213" }
		}
	}`

	n, k, points := parseInput(jsonData)
	if n == 0 {
		return
	}

	fmt.Println("Total roots provided (n):", n)
	fmt.Println("Minimum required roots (k):", k)

	matrix := make([][]*big.Int, k)
	for i := range matrix {
		matrix[i] = make([]*big.Int, k+1)
		for j := range matrix[i] {
			matrix[i][j] = new(big.Int)
		}
	}

	for i, p := range points {
		x := p[0]
		y := p[1]
		for j := 0; j < k; j++ {
			matrix[i][j].Exp(x, big.NewInt(int64(j)), nil)
		}
		matrix[i][k].Set(y)
	}

	coeffs := gaussSolve(matrix, k)
	fmt.Println("Constant term (c):", coeffs[0])
}
