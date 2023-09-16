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
	IC10 float64 `json:"ic10,omitempty"`
	IC11 float64 `json:"ic11,omitempty"`
	IC12 float64 `json:"ic12,omitempty"`
	IC13 float64 `json:"ic13,omitempty"`
	IC14 float64 `json:"ic14,omitempty"`
	IC16 float64 `json:"ic16,omitempty"`
	IC18 float64 `json:"ic18,omitempty"`
	IC20 float64 `json:"ic20,omitempty"`
	IC22 float64 `json:"ic22,omitempty"`
	IC25 float64 `json:"ic25,omitempty"`
	IC26 float64 `json:"ic26,omitempty"`
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
}

var icgs map[int]int = map[int]int{
	1: 1373,
	2: 1364,
	3: 1355,
	4: 1421,
	5: 1373,
	6: 1346,
	7: 1435,
}

var excelInfocodes map[int][]int = map[int][]int{
	1373: {8102, 8104, 8106, 8108, 8110, 8111, 8113, 8114, 8116, 8118, 8120, 8122, 8125},
	1364: {2102, 2104, 2106, 2108, 2110, 2112, 2114, 2116},
	1355: {2102, 2104, 2106, 2108, 2110, 2112, 2114, 2116},
	1421: {2102, 2104, 2106, 2108, 2110, 2112, 2114, 2116},
	1346: {2001, 2002, 2003, 2004, 2005, 2006, 2007, 2011, 2012, 2013},
	1435: {2102, 2104, 2106, 2108, 2110, 2112, 2114, 2116},
}

var alphabet = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "X", "Y", "Z"}

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
	} else if a.ICG == int(icgs[2]) {
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
	} else if a.ICG == int(icgs[3]) || a.ICG == int(icgs[7]) {
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
	} else if a.ICG == int(icgs[4]) {
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
	}
}

func NewApp(tsfrom int64, tsto int64, icg int, month int, daily bool, year int) *App {
	var url = ""
	if icg == 1373 || icg == 1364 || icg == 1355 || icg == 1421 || icg == 1435 {
		url = "https://iot-api.aiwater.io/iot/data/samples/graph-series-daily/icg:"
	} else {
		url = "https://iot-api.aiwater.io/iot/data/samples/graph-series/icg:"
	}
	return &App{
		BaseURL: url,
		Daily:   daily,
		Auth:    "Basic ZnJvbnRlbmQ6RzF2ZU0zRDQrNA==",
		TsFrom:  tsfrom,
		TsTo:    tsto,
		ICG:     icg,
		Month:   month,
		Year: year,
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
	fmt.Println("Select device [1-KTZLM: deminstanica, 2-RETTLH: kotolna turbiny, 3-GOTEC: reverzna osmoza, 4-SM62+RO-B1-3, 5-KTZLM: deminstanica tych vela, 6-VINCENTE, 7-GPV-rev. osmoza]")
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
	if err := f.SaveAs(fileName); err != nil {
		log.Fatal(err)
	}
}

func main() {
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
