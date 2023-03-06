package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

type App struct {
	BaseURL string
	TsFrom  int64
	TsTo    int64
	Auth    string
}

type Data struct {
	ChartData []Values `json:"chartData"`
	RealData  []Values `json:"realData"`
}

type Values struct {
	Date  string  `json:"datef"`
	C8102 float64 `json:"ic2"`
	C8104 float64 `json:"ic4"`
	C8106 float64 `json:"ic6"`
	C8107 float64 `json:"ic7"`
	C8108 float64 `json:"ic8"`
	C8110 float64 `json:"ic10"`
	C8111 float64 `json:"ic11"`
	C8113 float64 `json:"ic13"`
	C8114 float64 `json:"ic14"`
	C8116 float64 `json:"ic16"`
	C8118 float64 `json:"ic18"`
	C8120 float64 `json:"ic20"`
	C8122 float64 `json:"ic22"`
	C8125 float64 `json:"ic25"`
	C8126 float64 `json:"ic26"`
}

var lastDay map[int]int = map[int]int{
	1: 31,
	// priestupny rok
	2:  28,
	3:  31,
	4:  30,
	5:  31,
	6:  30,
	7:  31,
	8:  31,
	9:  30,
	10: 31,
	11: 30,
	12: 31,
}

func (a *App) getUrl() string {
	return fmt.Sprintf("%stsfrom=%d&tsto=%d", a.BaseURL, a.TsFrom, a.TsTo)
}

func NewApp(tsfrom, tsto int64) *App {
	return &App{
		BaseURL: "https://iot-api.aiwater.io/iot/data/samples/graph-series-daily/icg:1373?",
		Auth:    "Basic ZnJvbnRlbmQ6RzF2ZU0zRDQrNA==",
		TsFrom:  tsfrom,
		TsTo:    tsto,
	}
}

func getMonth() int {
	var month int
	fmt.Println("Select month [1 2 3 4 5 6 7 8 9 10 11 12]")
	_, err := fmt.Scanf("%d \n", &month)
	if err != nil {
		log.Panicf("error while getting value: %v", err)
	}

	return month
}

func getTimeStamp(month int) (int64, int64) {
	tsfrom := time.Date(2023, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location()).Unix()
	tsto := time.Date(2023, time.Month(month), lastDay[month], 0, 0, 0, 0, time.Now().Location()).Unix()

	return tsfrom, tsto
}

func (a *App) getData() Data {
	req, err := http.NewRequest(http.MethodGet, a.getUrl(), nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", a.Auth)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("could not read response body: %s\n", err)
		os.Exit(1)
	}
	// fmt.Printf("client: response body: %s\n", resBody)

	var data Data
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		fmt.Printf("could not unmarshal data: %v", err)
	}

	return data
}

func (a *App) writeExcel(data Data, month int) {
	f := excelize.NewFile()

	f.SetCellValue("Sheet1", "B1", "Infocodes")
	f.SetCellValue("Sheet1", "A2", "Date")
	f.SetCellValue("Sheet1", "B2", 8102)
	f.SetCellValue("Sheet1", "C2", 8104)
	f.SetCellValue("Sheet1", "D2", 8106)
	f.SetCellValue("Sheet1", "E2", 8108)
	f.SetCellValue("Sheet1", "F2", 8110)
	f.SetCellValue("Sheet1", "G2", 8111)
	f.SetCellValue("Sheet1", "H2", 8113)
	f.SetCellValue("Sheet1", "I2", 8114)
	f.SetCellValue("Sheet1", "J2", 8116)
	f.SetCellValue("Sheet1", "K2", 8118)
	f.SetCellValue("Sheet1", "L2", 8120)
	f.SetCellValue("Sheet1", "M2", 8122)
	f.SetCellValue("Sheet1", "N2", 8125)

	// for i := 0; i < lastDay[month]; i++ {
	for i, value := range data.RealData {
		// value := data.RealData[i]
		date := fmt.Sprintf("%d.%d.2023", i+1, month)
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+3), date)

		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+3), value.C8102)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+3), value.C8104)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+3), value.C8106)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+3), value.C8108)
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(i+3), value.C8110)
		f.SetCellValue("Sheet1", "G"+strconv.Itoa(i+3), value.C8111)
		f.SetCellValue("Sheet1", "H"+strconv.Itoa(i+3), value.C8113)
		f.SetCellValue("Sheet1", "I"+strconv.Itoa(i+3), value.C8114)
		f.SetCellValue("Sheet1", "J"+strconv.Itoa(i+3), value.C8116)
		f.SetCellValue("Sheet1", "K"+strconv.Itoa(i+3), value.C8118)
		f.SetCellValue("Sheet1", "L"+strconv.Itoa(i+3), value.C8120)
		f.SetCellValue("Sheet1", "M"+strconv.Itoa(i+3), value.C8122)
		f.SetCellValue("Sheet1", "N"+strconv.Itoa(i+3), value.C8125)
	}
	// f.SetCellValue("Sheet1", "A4", now.Format(time.ANSIC))

	if err := f.SaveAs("februar.xlsx"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	month := getMonth()
	app := NewApp(getTimeStamp(month))
	data := app.getData()
	app.writeExcel(data, month)
}
