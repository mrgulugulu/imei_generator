package main

func main() {

	res := generateImei("./data.xlsx", 500000)
	writeToExcel("res.xlsx", res)

}
