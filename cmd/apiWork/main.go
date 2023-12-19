package main

import (
	"bytes"
	"encoding/xml"
	//"fmt"
	"io"
	"log"
	"net/http"

	//"os"
	"golang.org/x/net/html/charset"
)

// ValCurs структура для разбора XML-ответа
type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valute []struct{
		NumCode		string `xml:"NumCode"`
		CharCode	string `xml:"CharCode"`
		Nominal		string `xml:"Nominal"`
		Name		string `xml:"Name"`
		Value		string `xml:"Value"`
		VunitRate	string `xml:"VunitRate"`
	} `xml:"Valute"`
}

// // Valute структура для представления информации о валюте
// type Valute struct {
// 	CharCode string  `xml:"CharCode"`
// 	Name     string  `xml:"Name"`
// 	Value    float64 `xml:"Value"`
// }

func main() {
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
		log.Println("Ошибка при разборе XML:", err)
		return
	}
	for _, valute := range valCurs.Valute{
		log.Println(valute.NumCode, valute.CharCode, valute.Nominal, valute.Name, valute.Value, valute.VunitRate)
	}

	// // Выведите информацию о курсах валют
	// log.Printf("Дата: %s\n", valCurs.Date)
	// for _, valute := range valCurs.Valutes {
	// 	log.Printf("%s (%s): %.2f\n", valute.Name, valute.CharCode, valute.Value)
	// }
}