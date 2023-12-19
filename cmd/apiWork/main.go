package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/net/html/charset"
)

// ValCurs структура для разбора XML-ответа
type ValCurs2 struct {
	XMLName 		xml.Name `xml:"ValCurs"`
	Name			string `xml:"name,attr"`
	Record []struct{
		Date		string `xml:"Date,attr"`
		Nominal		string `xml:"Nominal"`
		Value		string `xml:"Value"`
		VunitRate	string `xml:"VunitRate"`
	} `xml:"Record"`
}

type ValCurs struct {
	XMLName 		xml.Name `xml:"ValCurs"`
	Date			string `xml:"Date,attr"`
	Name			string `xml:"name,attr"`
	Valute []struct{
		ID			string `xml:"ID,attr"`
		NumCode		string `xml:"NumCode"`
		CharCode	string `xml:"CharCode"`
		Nominal		string `xml:"Nominal"`
		Name		string `xml:"Name"`
		Value		string `xml:"Value"`
		VunitRate	string `xml:"VunitRate"`
	} `xml:"Valute"`
}

func getInfo()(map[string]string){
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://www.cbr.ru/scripts/XML_daily.asp", nil)
	if err != nil{
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil{
		log.Fatalln(err)
	}
	//log.Println(string(body))
	// Разбор XML-ответа
	valCurs := new(ValCurs)
	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&valCurs)
	if err != nil {
		log.Fatalln("Ошибка при разборе XML:", err)
	}

	//mass := make([]string, len(valCurs.Valute))
	valuteMap := make(map[string]string)
	for _, valute := range valCurs.Valute{
		//log.Println(valCurs.Date, valute.ID, valute.NumCode, valute.CharCode, valute.Nominal, valute.Name, valute.Value, valute.VunitRate)
		valuteMap[valute.Name] = valute.ID
	}
	//log.Println(mass)
	return valuteMap
}

func getInfoFor90Day(str string, currentTime time.Time, previousTime time.Time, name string)([]string){
	url := fmt.Sprintf("http://www.cbr.ru/scripts/XML_dynamic.asp?date_req1=%02d/%02d/%d&date_req2=%02d/%02d/%d&VAL_NM_RQ=%s", previousTime.Day(), previousTime.Month(), previousTime.Year(), currentTime.Day(), currentTime.Month(), currentTime.Year(), str)
	client := &http.Client{}
	//log.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil{
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil{
		log.Fatalln(err)
	}
	//log.Println(string(body))
	// Разбор XML-ответа
	valCurs := new(ValCurs2)
	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&valCurs)
	if err != nil {
		log.Fatalln("Ошибка при разборе XML:", err)
	}
	min,err := strconv.ParseFloat(strings.ReplaceAll(valCurs.Record[0].Value, ",", "."),64)
	if err != nil {
		log.Fatalln("Ошибка при конвертации строки в число с плавающей точкой:", err)
	}
	var average_value float64
	max := min
	date_min:=valCurs.Record[0].Date
	date_max:=date_min
	for _, valute := range valCurs.Record{
		//log.Println(valute.Date, valute.Value, valute.VunitRate)
		current,err := strconv.ParseFloat(strings.ReplaceAll(valute.Value, ",", "."),64)
		if err != nil {
			log.Fatalln("Ошибка при конвертации строки в число с плавающей точкой:", err)
		}
		if min>current{
			min = current
			date_min = valute.Date
		}
		if max<current{
			max = current
			date_max = valute.Date
		}
		average_value+=current
	}
	average_value/=float64(len(valCurs.Record))
	average_value_cropped:=fmt.Sprintf("%.2f",average_value)
	//log.Print("min: ",min," date min: ", date_min," max: ", max," date max ", date_max," average value: ", average_value_cropped,"\n")
	output_string := []string{name, strconv.FormatFloat(min, 'f', -2, 64), date_min, strconv.FormatFloat(max, 'f', -2, 64), date_max, average_value_cropped}
	return output_string
}

func main() {
	log.SetFlags(0)
	currentTime := time.Now()
	previousTime := currentTime.AddDate(0, 0, -90)
	log.Printf("Текущая дата: %02d-%02d-%d\n", currentTime.Day(), currentTime.Month(), currentTime.Year())
	//log.Printf("Дата 90 дней назад: %02d-%02d-%d\n", previousTime.Day(), previousTime.Month(), previousTime.Year())
	log.Println("Курс валют за последнии 90 дней")
	var valuteMap map[string]string = getInfo()
	//log.Println("first_stage_end")
	mass := make([][]string, len(valuteMap))
	i:=0
	for name, ID := range valuteMap{
		//fmt.Printf("%s ", name)
		mass[i]=getInfoFor90Day(ID, currentTime, previousTime, name)
		i++
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "min", "date min", "max", "date max", "average cost"})
	for _, row := range mass {
		table.Append(row)
	}
	table.Render()
}