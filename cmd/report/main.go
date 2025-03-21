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
	ICG     int
	Values  []float64
	Dates   []string
	Month   int
	Daily   bool
	Year    int
}

type Data struct {
	ChartData []Values `json:"chartData"`
	RealData  []Values `json:"realData"`
}

type Values struct {
  Date string  `json:"datef"`
  IC01 float64 `json:"ic1,omitempty"`
  IC02 float64 `json:"ic2,omitempty"`
  IC03 float64 `json:"ic3,omitempty"`
  IC04 float64 `json:"ic4,omitempty"`
  IC05 float64 `json:"ic5,omitempty"`
  IC06 float64 `json:"ic6,omitempty"`
  IC07 float64 `json:"ic7,omitempty"`
  IC08 float64 `json:"ic8,omitempty"`
  IC09 float64 `json:"ic9,omitmepty"`
  IC10 float64 `json:"ic10,omitempty"`
  IC11 float64 `json:"ic11,omitempty"`
  IC12 float64 `json:"ic12,omitempty"`
  IC13 float64 `json:"ic13,omitempty"`
  IC14 float64 `json:"ic14,omitempty"`
  IC15 float64 `json:"ic15,omitempty"`
  IC16 float64 `json:"ic16,omitempty"`
  IC17 float64 `json:"ic17,omitempty"`
  IC18 float64 `json:"ic18,omitempty"`
  IC19 float64 `json:"ic19,omitempty"`
  IC20 float64 `json:"ic20,omitempty"`
  IC21 float64 `json:"ic21,omitempty"`
  IC22 float64 `json:"ic22,omitempty"`
  IC23 float64 `json:"ic23,omitempty"`
  IC24 float64 `json:"ic24,omitempty"`
  IC25 float64 `json:"ic25,omitempty"`
  IC26 float64 `json:"ic26,omitempty"`
  IC27 float64 `json:"ic27,omitempty"`
  IC28 float64 `json:"ic28,omitempty"`
  IC29 float64 `json:"ic29,omitempty"`
  IC30 float64 `json:"ic30,omitempty"`
  IC32 float64 `json:"ic32,omitempty"`
  IC33 float64 `json:"ic33,omitempty"`
  IC34 float64 `json:"ic34,omitempty"`
  IC36 float64 `json:"ic36,omitempty"`
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

var nameMonth map[int]string = map[int]string{
	1:  "Januar",
	2:  "Februar",
	3:  "Marec",
	4:  "April",
	5:  "Maj",
	6:  "Jun",
	7:  "Jul",
	8:  "August",
	9:  "September",
	10: "Oktober",
	11: "November",
	12: "December",
}

var devices map[int]string = map[int]string{
	1373: "KTLZM-demistanica",
	1364: "RETTLH-kotolna-turbiny",
	1355: "GOTEC-reverzna-osmoza",
	1421: "SM62+RO-B1-3",
	1346: "VINCENTE",
	1435: "GPV-osmoza",
	1281: "TESGAL-RO-podesta",
	1393: "Dialyza-Levice",
	1145: "SEMI-RO1",
	1489: "SEMI-RO2",
	1361: "RONA",
	1475: "MORO",
	1545: "ICSNR",
	1091: "MORO3",
	1396: "MORO4",
  1510: "TANAWA",
  1408: "AEROSOL",
  1556: "TERMINAL-R1",
  1516: "BUCANY",
}

var icgs map[int]int = map[int]int{
	1:  1373,
	2:  1364,
	3:  1355,
	4:  1421,
	5:  1373,
	6:  1346,
	7:  1435,
	8:  1281,
	9:  1393,
	10: 1145,
	11: 1489,
	12: 1361,
	13: 1475,
	14: 1545,
	15: 1091,
	16: 1396,
  17: 1510,
  18: 1408,
  19: 1556,
  20: 1516,
}

var excelInfocodes map[int][]int = map[int][]int{
	1373: {8102, 8104, 8106, 8108, 8110, 8111, 8113, 8114, 8116, 8118, 8120, 8122, 8125},
	1364: {2102, 2104, 2106, 2108, 2110, 2112, 2114, 2116},
	1355: {2102, 2104, 2106, 2108, 2110, 2112, 2114, 2116},
	1421: {2102, 2104, 2106, 2108, 2110, 2112, 2114, 2116},
	1346: {2001, 2002, 2003, 2004, 2005, 2006, 2007, 2011, 2012, 2013},
	1435: {2102, 2104, 2106, 2108, 2110, 2112, 2114, 2116},
	1281: {2102, 2104, 2106, 2108, 2110, 2112, 2114, 2116},
	1393: {2102, 2104, 2106, 2108, 2110, 2112, 2114, 2116},
	1145: {2102, 2104, 2106, 2108, 2110, 2112, 2114, 2116},
	1489: {2102, 2104, 2106, 2108, 2110, 2112, 2114, 2116},
	1475: {2102, 2104, 2106, 2108, 2110, 2112, 2114, 2116},
	1361: {5102, 5105, 5106, 5107, 5108, 5109, 5110, 5112, 5114, 5116, 5118, 5120, 5122, 5124, 5126, 5128, 5130, 5132, 5133, 5134, 5136},
	1545: {2102, 2104, 2106, 2108, 2110, 2112, 2114, 2116},
	1091: {2102, 2104, 2106, 2108, 2110, 2112, 2114, 2116},
	1396: {2102, 2104, 2106, 2108, 2110, 2112, 2114, 2116},
  1510: {4102, 4104, 4106, 4108, 4110, 4111, 4112, 4113, 4114, 4115, 4116, 4117, 4118, 4119, 4120, 4121, 4122, 4123, 4124, 4125, 4126, 4127, 4128, 4129, 4130},
  1408: {8102, 8104, 8106, 8108, 8110, 8111, 8112, 8113, 8114, 8115, 8116, 8117, 8118, 8119, 8120, 8121, 8122, 8123, 8124, 8125, 8126, 8127},
  1556: {3102, 3104, 3106, 3108, 3110, 3111, 3112, 3114, 3116, 3118, 3119, 3120, 3121, 3122, 3123, 3124, 3126, 3127, 3128, 3130},
  1516: {4102, 4104, 4106, 4108, 4110, 4111, 4112, 4114, 4116, 4118, 4119, 4120, 4121, 4122, 4123, 4124, 4126, 4127, 4130},
}

var alphabet = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "X", "Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH"}

func (a *App) getUrl(tsfrom, tsto int64) string {
	fmt.Println(fmt.Sprintf("%s%d?tsfrom=%d&tsto=%d", a.BaseURL, a.ICG, tsfrom, tsto))
	return fmt.Sprintf("%s%d?tsfrom=%d&tsto=%d", a.BaseURL, a.ICG, tsfrom, tsto)
}

func (a *App) fillArray(data Data) {

	if a.ICG == int(icgs[1]) {
		for _, value := range data.RealData {
			a.Values = append(a.Values, value.IC02)
			a.Values = append(a.Values, value.IC04)
			a.Values = append(a.Values, value.IC06)
			a.Values = append(a.Values, value.IC08)
			a.Values = append(a.Values, value.IC10)
			a.Values = append(a.Values, value.IC11)
			a.Values = append(a.Values, value.IC13)
			a.Values = append(a.Values, value.IC14)
			a.Values = append(a.Values, value.IC16)
			a.Values = append(a.Values, value.IC18)
			a.Values = append(a.Values, value.IC20)
			a.Values = append(a.Values, value.IC22)
			a.Values = append(a.Values, value.IC25)

			// formatedDate, err := time.Parse("2006-01-02", strings.Split(value.Date, " ")[0])
			// if err != nil {
			// 	panic(err)
			// }
			// a.Dates = append(a.Dates, formatedDate)
			// formatedDate, _ := time.Parse("2006-01-02T15:04:05", strings.Split(value.Date, " ")[0])
			a.Dates = append(a.Dates, value.Date)
		}
	} else if a.ICG == int(icgs[2]) || a.ICG == int(icgs[11]) || a.ICG == int(icgs[13]) || a.ICG == int(icgs[14]) || a.ICG == int(icgs[15]) || a.ICG == int(icgs[16]) {
		for _, value := range data.RealData {
			a.Values = append(a.Values, value.IC02)
			a.Values = append(a.Values, value.IC04)
			a.Values = append(a.Values, value.IC06)
			a.Values = append(a.Values, value.IC08)
			a.Values = append(a.Values, value.IC10)
			a.Values = append(a.Values, value.IC12)
			a.Values = append(a.Values, value.IC14)
			a.Values = append(a.Values, value.IC16)

			// formatedDate, err := time.Parse("2006-01-02", strings.Split(value.Date, " ")[0])
			// if err != nil {
			// 	panic(err)
			// }
			a.Dates = append(a.Dates, value.Date)
		}
	} else if a.ICG == int(icgs[3]) || a.ICG == int(icgs[7]) || a.ICG == int(icgs[8]) {
		for _, value := range data.RealData {
			a.Values = append(a.Values, value.IC02)
			a.Values = append(a.Values, value.IC04)
			a.Values = append(a.Values, value.IC06)
			a.Values = append(a.Values, value.IC08)
			a.Values = append(a.Values, value.IC10)
			a.Values = append(a.Values, value.IC12)
			a.Values = append(a.Values, value.IC14)
			a.Values = append(a.Values, value.IC16)

			// formatedDate, err := time.Parse("2006-01-02", strings.Split(value.Date, " ")[0])
			// if err != nil {
			// 	panic(err)
			// }
			a.Dates = append(a.Dates, value.Date)
		}
	} else if a.ICG == int(icgs[4]) || a.ICG == int(icgs[9]) || a.ICG == int(icgs[10]) {
		for _, value := range data.RealData {
			a.Values = append(a.Values, value.IC02)
			a.Values = append(a.Values, value.IC04)
			a.Values = append(a.Values, value.IC06)
			a.Values = append(a.Values, value.IC08)
			a.Values = append(a.Values, value.IC10)
			a.Values = append(a.Values, value.IC12)
			a.Values = append(a.Values, value.IC14)
			a.Values = append(a.Values, value.IC16)

			// formatedDate, err := time.Parse("2006-01-02", strings.Split(value.Date, " ")[0])
			// if err != nil {
			// 	panic(err)
			// }
			a.Dates = append(a.Dates, value.Date)
		}
  } else if a.ICG == int(icgs[18]) {
    for _, value := range data.RealData {
      a.Values = append(a.Values, value.IC02)
      a.Values = append(a.Values, value.IC04)
      a.Values = append(a.Values, value.IC06)
      a.Values = append(a.Values, value.IC08)
      a.Values = append(a.Values, value.IC10)
      a.Values = append(a.Values, value.IC11)
      a.Values = append(a.Values, value.IC12)
      a.Values = append(a.Values, value.IC13)
      a.Values = append(a.Values, value.IC14)
      a.Values = append(a.Values, value.IC15)
      a.Values = append(a.Values, value.IC16)
      a.Values = append(a.Values, value.IC17)
      a.Values = append(a.Values, value.IC18)
      a.Values = append(a.Values, value.IC19)
      a.Values = append(a.Values, value.IC20)
      a.Values = append(a.Values, value.IC21)
      a.Values = append(a.Values, value.IC22)
      a.Values = append(a.Values, value.IC23)
      a.Values = append(a.Values, value.IC24)
      a.Values = append(a.Values, value.IC25)
      a.Values = append(a.Values, value.IC26)
      a.Values = append(a.Values, value.IC27)

      a.Dates = append(a.Dates, value.Date)
    }
  } else if a.ICG == int(icgs[17]) {
    for _, value := range data.RealData {
      a.Values = append(a.Values, value.IC02)
      a.Values = append(a.Values, value.IC04)
      a.Values = append(a.Values, value.IC06)
      a.Values = append(a.Values, value.IC08)
      a.Values = append(a.Values, value.IC10)
      a.Values = append(a.Values, value.IC11)
      a.Values = append(a.Values, value.IC12)
      a.Values = append(a.Values, value.IC13)
      a.Values = append(a.Values, value.IC14)
      a.Values = append(a.Values, value.IC15)
      a.Values = append(a.Values, value.IC16)
      a.Values = append(a.Values, value.IC17)
      a.Values = append(a.Values, value.IC18)
      a.Values = append(a.Values, value.IC19)
      a.Values = append(a.Values, value.IC20)
      a.Values = append(a.Values, value.IC21)
      a.Values = append(a.Values, value.IC22)
      a.Values = append(a.Values, value.IC23)
      a.Values = append(a.Values, value.IC24)
      a.Values = append(a.Values, value.IC25)
      a.Values = append(a.Values, value.IC26)
      a.Values = append(a.Values, value.IC27)
      a.Values = append(a.Values, value.IC28)
      a.Values = append(a.Values, value.IC29)
      a.Values = append(a.Values, value.IC30)

      // formatedDate, err := time.Parse("2006-01-02", strings.Split(value.Date, " ")[0])
      // if err != nil {
      //  panic(err)
      // }
      a.Dates = append(a.Dates, value.Date)
    }
	} else if a.ICG == int(icgs[6]) {
		for _, value := range data.RealData {
			a.Values = append(a.Values, value.IC01)
			a.Values = append(a.Values, value.IC02)
			a.Values = append(a.Values, value.IC03)
			a.Values = append(a.Values, value.IC04)
			a.Values = append(a.Values, value.IC05)
			a.Values = append(a.Values, value.IC06)
			a.Values = append(a.Values, value.IC07)
			a.Values = append(a.Values, value.IC11)
			a.Values = append(a.Values, value.IC12)
			a.Values = append(a.Values, value.IC13)

			a.Dates = append(a.Dates, value.Date)
		}
	} else if a.ICG == int(icgs[19]) {
		for _, value := range data.RealData {
			a.Values = append(a.Values, value.IC02)
			a.Values = append(a.Values, value.IC04)
			a.Values = append(a.Values, value.IC06)
			a.Values = append(a.Values, value.IC08)
			a.Values = append(a.Values, value.IC10)
			a.Values = append(a.Values, value.IC11)
			a.Values = append(a.Values, value.IC12)
			a.Values = append(a.Values, value.IC14)
			a.Values = append(a.Values, value.IC16)
			a.Values = append(a.Values, value.IC18)
			a.Values = append(a.Values, value.IC19)
			a.Values = append(a.Values, value.IC20)
			a.Values = append(a.Values, value.IC21)
			a.Values = append(a.Values, value.IC22)
			a.Values = append(a.Values, value.IC23)
			a.Values = append(a.Values, value.IC24)
			a.Values = append(a.Values, value.IC26)
			a.Values = append(a.Values, value.IC27)
			a.Values = append(a.Values, value.IC28)
			a.Values = append(a.Values, value.IC30)

			a.Dates = append(a.Dates, value.Date)
		}
	} else if a.ICG == int(icgs[20]) {
		for _, value := range data.RealData {
			a.Values = append(a.Values, value.IC02)
			a.Values = append(a.Values, value.IC04)
			a.Values = append(a.Values, value.IC06)
			a.Values = append(a.Values, value.IC08)
			a.Values = append(a.Values, value.IC10)
			a.Values = append(a.Values, value.IC11)
			a.Values = append(a.Values, value.IC12)
			a.Values = append(a.Values, value.IC14)
			a.Values = append(a.Values, value.IC16)
			a.Values = append(a.Values, value.IC18)
			a.Values = append(a.Values, value.IC19)
			a.Values = append(a.Values, value.IC20)
			a.Values = append(a.Values, value.IC21)
			a.Values = append(a.Values, value.IC22)
			a.Values = append(a.Values, value.IC23)
			a.Values = append(a.Values, value.IC24)
			a.Values = append(a.Values, value.IC26)
			a.Values = append(a.Values, value.IC27)
			a.Values = append(a.Values, value.IC30)

			a.Dates = append(a.Dates, value.Date)
		}
	} else if a.ICG == int(icgs[12]) {
		for _, value := range data.RealData {
			a.Values = append(a.Values, value.IC02)
			a.Values = append(a.Values, value.IC05)
			a.Values = append(a.Values, value.IC06)
			a.Values = append(a.Values, value.IC07)
			a.Values = append(a.Values, value.IC08)
			a.Values = append(a.Values, value.IC09)
			a.Values = append(a.Values, value.IC10)
			a.Values = append(a.Values, value.IC12)
			a.Values = append(a.Values, value.IC14)
			a.Values = append(a.Values, value.IC16)
			a.Values = append(a.Values, value.IC18)
			a.Values = append(a.Values, value.IC20)
			a.Values = append(a.Values, value.IC22)
			a.Values = append(a.Values, value.IC24)
			a.Values = append(a.Values, value.IC26)
			a.Values = append(a.Values, value.IC28)
			a.Values = append(a.Values, value.IC30)
			a.Values = append(a.Values, value.IC32)
			a.Values = append(a.Values, value.IC33)
			a.Values = append(a.Values, value.IC34)
			a.Values = append(a.Values, value.IC36)

			a.Dates = append(a.Dates, value.Date)
		}
	}
}

func NewApp(tsfrom int64, tsto int64, icg int, month int, daily bool, year int) *App {
	var url = ""
	if icg == 1346 {
		url = "https://iot-api.aiwater.io/iot/data/samples/graph-series/icg:"
	} else {
		url = "https://iot-api.aiwater.io/iot/data/samples/graph-series-daily/icg:"
	}
	return &App{
		BaseURL: url,
		Daily:   daily,
		Auth:    "Basic ZnJvbnRlbmQ6RzF2ZU0zRDQrNA==",
		TsFrom:  tsfrom,
		TsTo:    tsto,
		ICG:     icg,
		Month:   month,
		Year:    year,
	}
}

func getYear() int {
	var year int
	fmt.Println("Insert year")
	_, err := fmt.Scanf("%d \n", &year)
	if err != nil {
		log.Panicf("error while getting value: %v", err)
	}

	return year
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

func getICG() (int, bool) {
	var device int
	fmt.Print("Select device [1-KTZLM: deminstanica, 2-RETTLH: kotolna turbiny, 3-GOTEC: reverzna osmoza,")
	fmt.Print("4-SM62+RO-B1-3, 5-KTZLM: deminstanica tych vela, 6-VINCENTE, 7-GPV-rev. osmoza, ")
	fmt.Print("8-TESGAL:RO podesta, 9-Dialyza Levice, 10-SEMI: RO1, 11-SEMI: RO2, 12-RONA: chladenie, ")
	fmt.Print("13-MORO, 14 - ICSNR, 15-MORO3, 16-MORO4, 17-TANAWA, 18-AEROSOL, 19-TERMINAL-R1, 20-BUCANY]\n")
	_, err := fmt.Scanf("%d \n", &device)
	if err != nil {
		log.Panicf("error while getting value: %v", err)
	}

	daily := true
	if device == 5 || device == 6 {
		daily = false
	}

	return icgs[device], daily
}

func getTimeStamp(month int, year int) (int64, int64) {
	tsfrom := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location()).Unix()
	tsto := time.Date(year, time.Month(month), lastDay[month], 0, 0, 0, 0, time.Now().Location()).Unix()

	return tsfrom, tsto
}

func getTimeStampOfDay(month int, day int, half int) (int64, int64) {
	tsfrom := time.Date(2023, time.Month(month), day, half*12-12, 0, 0, 0, time.Now().Location()).Unix()
	// var tsto int64
	// if day+1 > lastDay[month] {
	// 	tsto = time.Date(2023, time.Month(month+1), 1, 0, 0, 0, 0, time.Now().Location()).Unix()
	// } else {
	// 	tsto = time.Date(2023, time.Month(month), lastDay[month], 0, 0, 0, 0, time.Now().Location()).Unix()
	// }
	tsto := time.Date(2023, time.Month(month), day, half*12, 0, 0, 0, time.Now().Location()).Unix()

	return tsfrom, tsto
}

func (a *App) noIdeaHowToNameThis() Data {
	if a.Daily {
		return a.getData(a.getUrl(a.TsFrom, a.TsTo))
	}
	var data Data
	for i := 1; i <= lastDay[a.Month]; i++ {
		tmpData := a.getData(a.getUrl(getTimeStampOfDay(a.Month, i, 1)))
		data.ChartData = append(data.ChartData, tmpData.ChartData...)
		data.RealData = append(data.RealData, tmpData.RealData...)
		tmpData = a.getData(a.getUrl(getTimeStampOfDay(a.Month, i, 2)))
		data.ChartData = append(data.ChartData, tmpData.ChartData...)
		data.RealData = append(data.RealData, tmpData.RealData...)
	}
	return data

}

func (a *App) getData(url string) Data {
	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	// fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("could not read response body: %s\n", err)
		os.Exit(1)
	}

	var data Data
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		fmt.Printf("could not unmarshal data: %v", err)
	}
	return data
}

func (a *App) writeExcel(data []float64) {
	f := excelize.NewFile()

	err := f.SetCellValue("Sheet1", "B1", "Infocodes")
	if err != nil {
		panic(err)
	}
	err = f.SetCellValue("Sheet1", "A2", "Date")
	if err != nil {
		panic(err)
	}

	codes := excelInfocodes[int(a.ICG)]

	for i := 0; i < len(codes); i++ {
		err = f.SetCellValue("Sheet1", alphabet[i+1]+"2", codes[i])
		if err != nil {
			panic(err)
		}
	}

	row := 2
	for i, value := range data {
		if i%(len(codes)) == 0 {
			// date := fmt.Sprint(a.Dates[row-2].AddDate(0, 0, 1).Format("02.01.2006"))
			// date := fmt.Sprint(a.Dates[row-2].AddDate(0, 0, 1).Format("02.01.2006"))
			err = f.SetCellValue("Sheet1", "A"+strconv.Itoa(row+1), a.Dates[row-2])
			row += 1
			if err != nil {
				panic(err)
			}
		}

		err = f.SetCellValue("Sheet1", alphabet[(i+len(codes))%len(codes)+1]+strconv.Itoa(row), value)
		if err != nil {
			panic(err)
		}
	}

	fileName := fmt.Sprintf("%s-%s.xlsx", devices[a.ICG], nameMonth[a.Month])
	fmt.Print(fileName)
	if err := f.SaveAs(fileName); err != nil {
		log.Fatal(err)
	}
}

func main() {
  fmt.Println("Lucia je najkrajšia na svete")
	// daily := true
	year := getYear()
	month := getMonth()
	icg, daily := getICG()
	tsfrom, tsto := getTimeStamp(month, year)
	app := NewApp(tsfrom, tsto, icg, month, daily, year)
	fmt.Println(app)
	data := app.noIdeaHowToNameThis()
	// fmt.Println(data)
	app.fillArray(data)
	app.writeExcel(app.Values)
}

