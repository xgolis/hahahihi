package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	Dates   []time.Time
	Month   int
}

type Data struct {
	ChartData []Values `json:"chartData"`
	RealData  []Values `json:"realData"`
}

type Values struct {
	Date string  `json:"datef"`
	IC02 float64 `json:"ic2,omitempty"`
	IC04 float64 `json:"ic4,omitempty"`
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
}

var icgs map[int]int = map[int]int{
	1: 1373,
	2: 1364,
	3: 1355,
	4: 1421,
}

var excelInfocodes map[int][]int = map[int][]int{
	1373: {8102, 8104, 8106, 8108, 8110, 8111, 8113, 8114, 8116, 8118, 8120, 8122, 8125},
	1364: {2102, 2104, 2106, 2108, 2110, 2112, 2114, 2116},
	1355: {2102, 2104, 2106, 2108, 2110, 2112, 2114, 2116},
	1421: {2102, 2104, 2106, 2108, 2110, 2112, 2114, 2116},
}

var alphabet = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "X", "Y", "Z"}

func (a *App) getUrl() string {
	return fmt.Sprintf("%s%d?tsfrom=%d&tsto=%d", a.BaseURL, a.ICG, a.TsFrom, a.TsTo)
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

			formatedDate, err := time.Parse("2006-01-02", strings.Split(value.Date, " ")[0])
			if err != nil {
				panic(err)
			}
			a.Dates = append(a.Dates, formatedDate)
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

			formatedDate, err := time.Parse("2006-01-02", strings.Split(value.Date, " ")[0])
			if err != nil {
				panic(err)
			}
			a.Dates = append(a.Dates, formatedDate)
		}
	} else if a.ICG == int(icgs[3]) {
		for _, value := range data.RealData {
			a.Values = append(a.Values, value.IC02)
			a.Values = append(a.Values, value.IC04)
			a.Values = append(a.Values, value.IC06)
			a.Values = append(a.Values, value.IC08)
			a.Values = append(a.Values, value.IC10)
			a.Values = append(a.Values, value.IC12)
			a.Values = append(a.Values, value.IC14)
			a.Values = append(a.Values, value.IC16)

			formatedDate, err := time.Parse("2006-01-02", strings.Split(value.Date, " ")[0])
			if err != nil {
				panic(err)
			}
			a.Dates = append(a.Dates, formatedDate)
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

			formatedDate, err := time.Parse("2006-01-02", strings.Split(value.Date, " ")[0])
			if err != nil {
				panic(err)
			}
			a.Dates = append(a.Dates, formatedDate)
		}
	}
}

func NewApp(tsfrom int64, tsto int64, icg int, month int) *App {
	return &App{
		BaseURL: "https://iot-api.aiwater.io/iot/data/samples/graph-series-daily/icg:",
		Auth:    "Basic ZnJvbnRlbmQ6RzF2ZU0zRDQrNA==",
		TsFrom:  tsfrom,
		TsTo:    tsto,
		ICG:     icg,
		Month:   month,
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

func getICG() int {
	var device int
	fmt.Println("Select device [1-KTZLM: deminstanica, 2-RETTLH: kotolna turbiny, 3-GOTEC: reverzna osmoza, 4-SM62+RO-B1-3]")
	_, err := fmt.Scanf("%d \n", &device)
	if err != nil {
		log.Panicf("error while getting value: %v", err)
	}

	return icgs[device]
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
			date := fmt.Sprint(a.Dates[row-2].AddDate(0, 0, 1).Format("02.01.2006"))
			row += 1
			err = f.SetCellValue("Sheet1", "A"+strconv.Itoa(row), date)
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
	month := getMonth()
	icg := getICG()
	tsfrom, tsto := getTimeStamp(month)
	app := NewApp(tsfrom, tsto, icg, month)
	data := app.getData()
	app.fillArray(data)
	app.writeExcel(app.Values)
}
