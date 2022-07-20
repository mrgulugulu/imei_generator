package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// generateImei 生成随机imei到指定路径的文件
// TAC(6位,0-5位):类型分配码，有一定的规则可以从excel中获取
// FAC(2位,6-7位):最终装配地代码，可以随机生成
// SNR(6位, 8-13位):序列号，可以随机生成
// CD(1位):验证码，有固定的算法
func generateImei(inputPath string, amount int) []string {

	tac := readFromExcel(inputPath)

	res := make([]string, amount)
	for i := range res {
		middleNumStr := strconv.Itoa(rand.Int() % 100000000)
		res[i] = tac[rand.Int()%(len(tac)-1)] + middleNumStr
		res[i] = res[i] + strconv.Itoa(luhn(res[i]))
	}
	return res

}

// luhn 生成最后一位的验证码
func luhn(prefix string) int {
	var total, sum1, sum2 int
	n := len(prefix)
	for i := 0; i < n; i++ {
		num, _ := strconv.Atoi(string(prefix[i]))
		// 奇数
		if i%2 == 0 {
			sum1 += num
		} else { // 偶数
			tmp := num * 2
			if tmp < 10 {
				sum2 += tmp
			} else {
				sum2 = sum2 + tmp + 1 - 10
			}
		}
	}
	total = sum1 + sum2
	if total%10 == 0 {
		return 0
	} else {
		return 10 - (total % 10)
	}
}

// readFromExcel 从excel中读取数据
func readFromExcel(path string) []string {
	res := []string{}
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Errorf("open errors %v", err)
	}
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Errorf("get rows errors %v", err)
	}

	for _, row := range rows {
		for i := range row {
			if i == 0 {
				res = append(res, row[i])
			}
		}
	}
	return res
}

func writeToExcel(path string, data []string) {
	file := excelize.NewFile()
	streamWriter, err := file.NewStreamWriter("Sheet1")
	if err != nil {
		fmt.Println(err)
	}
	if err := streamWriter.SetRow("A1", []interface{}{
		excelize.Cell{Value: "Data"}}); err != nil {
		fmt.Println(err)
	}
	for rowId := 1; rowId < len(data)+1; rowId++ {
		row := make([]interface{}, 3)
		row[0], row[1] = data[rowId-1], encode(data[rowId-1])
		cell, _ := excelize.CoordinatesToCellName(1, rowId)
		if err := streamWriter.SetRow(cell, row); err != nil {
			fmt.Println(err)
		}
	}
	if err := streamWriter.Flush(); err != nil {
		fmt.Println(err)
	}
	if err := file.SaveAs(path); err != nil {
		fmt.Println(err)
	}

}

func encode(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
