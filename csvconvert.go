package main

import (
	//"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var ADMS13OutputKeyOrder = []string{"Channel Number", "Receive Frequency", "Transmit Frequency", "Offset Frequency", "Offset Direction", "Operating Mode", "TAG", "Name", "Tone Mode", "CTCSS Frequency", "DCS", "User CTCSS", "Tx Power", "Scan", "Step", "Narrow", "Clock Shift", "Comment", "Zero"}
var ADMS13ToneValues = map[string]string{"Tone": "TONE", "": "OFF"}
var ADMS13ScanValues = map[string]string{"": "YES", "S": "NO"}
var ADMS13DummyRow = ",,,,,,,,,,,,,,,,,,0\n"

// var offsetValues = map[string]string{"5.00 MHz": "5.00", "7.85 MHz": "7.85", "8.00 MHz": "8.00", "9.15 MHz": "9.15", "600 kHz": "0.600000", " ": ""}

// var offsetDirection = map[string]string{"Minus": "-", "Plus": "+", "Simplex": "", "Split": ""}
// var offsetValues = map[string]string{"5.00 MHz": "5.00", "7.85 MHz": "7.85", "8.00 MHz": "8.00", "9.15 MHz": "9.15", "600 kHz": "0.600000", " ": ""}
// var toneValues = map[string]string{"Tone": "Tone", "None": ""}
// var skipValues = map[string]string{"Off": "", "Skip": "S", "XXX": "P"}

func main() {

	// read csv file
	InputCSVSliceMap, err := CSVFileToMap("XCZFreqListv1_01.csv")
	checkError("Can not read the csv file", err)

	createADMS13OutputFile(InputCSVSliceMap)
}

func createADMS13OutputFile(inputMap []map[string]string) {
	// open the output file
	outputfile, err := os.OpenFile("ADMS13.csv", os.O_CREATE|os.O_WRONLY, 0777)
	defer outputfile.Close()
	checkError("Couldn't create the output csv file", err)

	// Iterate through the records
	rowNumber := 1
	for _, record := range inputMap {
		thisRowNUmber, _ := strconv.Atoi(record["Location"])
		for rowNumber < thisRowNUmber {
			outputfile.WriteString(strconv.Itoa(rowNumber) + ADMS13DummyRow)
			rowNumber++
		}
		rowNumber++
		outputRow := createADMS13OutputRow(record)
		_, err = outputfile.WriteString(outputRow)
		checkError("Error writing the output csv file", err)
	}
	for rowNumber < 1000 {
		outputfile.WriteString(strconv.Itoa(rowNumber) + ADMS13DummyRow)
		rowNumber++
	}
}

func createADMS13OutputRow(inputMap map[string]string) string {
	outputMap := map[string]string{}

	fmt.Printf("Channel Number = %s \n", inputMap["Location"])

	// Location,Name,Frequency,Duplex,Offset,Tone,rToneFreq,cToneFreq,DtcsCode,DtcsPolarity,Mode,TStep,Skip,Comment,URCALL,RPT1CALL,RPT2CALL

	// calculate transmit frequency
	freq, _ := strconv.ParseFloat(inputMap["Frequency"], 64)
	offsetvalue, _ := strconv.ParseFloat(inputMap["Offset"], 64)
	offsetFrequency := ""
	offsetDirection := ""
	transmitFrequency := freq
	switch inputMap["Duplex"] {
	case "+":
		offsetDirection = "+RPT"
		offsetFrequency = inputMap["Offset"]
		transmitFrequency = freq + offsetvalue
	case "-":
		offsetDirection = "-RPT"
		offsetFrequency = inputMap["Offset"]
		transmitFrequency = freq - offsetvalue
	case "":
		offsetDirection = "OFF"
		offsetFrequency = "0.00000"
		transmitFrequency = freq
	case "off":
		offsetDirection = "OFF"
		offsetFrequency = "0.00000"
		transmitFrequency = freq
	case "split":
		offsetDirection = "-/+"
		offsetFrequency = "0.00000"
		transmitFrequency = offsetvalue
	}
	if offsetFrequency == "0.600000" {
		offsetFrequency = "0.60000"
	}

	outputMap["Channel Number"] = inputMap["Location"]
	outputMap["Receive Frequency"] = inputMap["Frequency"]
	outputMap["Transmit Frequency"] = fmt.Sprintf("%.5f", transmitFrequency)
	outputMap["Offset Frequency"] = offsetFrequency
	outputMap["Offset Direction"] = offsetDirection
	outputMap["Operating Mode"] = inputMap["Mode"]
	outputMap["TAG"] = "ON"
	outputMap["Name"] = "" // does not like it   inputMap["Name"] // [:6]
	outputMap["Tone Mode"] = ADMS13ToneValues[inputMap["Tone"]]
	outputMap["CTCSS Frequency"] = inputMap["rToneFreq"] + " Hz"
	outputMap["DCS"] = inputMap["DtcsCode"]
	outputMap["User CTCSS"] = "1500 Hz" //
	outputMap["Tx Power"] = "HIGH"      // "Tx Power"]
	outputMap["Scan"] = ADMS13ScanValues[inputMap["Skip"]]
	outputMap["Step"] = "5.0KHz" //    ??
	outputMap["Narrow"] = "OFF"  //   ?? narrow FM
	outputMap["Clock Shift"] = "OFF"
	outputMap["Comment"] = inputMap["Comment"]
	outputMap["Zero"] = "0"

	var concatenatestrings strings.Builder

	for _, keyName := range ADMS13OutputKeyOrder {
		concatenatestrings.WriteString(outputMap[keyName])
		if keyName != "Zero" {
			concatenatestrings.WriteString(",")
		}
	}
	concatenatestrings.WriteString("\n")
	outputRow := concatenatestrings.String()
	fmt.Printf("Output row %s\n", outputRow)
	return outputRow
}

// 	frequency, err := strconv.ParseFloat(inputRow["Receive Frequency"], 8)
// 	isHam := false
// 	if (frequency > 144.0 && frequency < 148.0) || (frequency > 420.0 && frequency < 450.0) {
// 		isHam = true
// 	}
// 	fmt.Println(frequency, err, isHam)

// 	isSplit := false
// 	if inputRow["Offset Direction"] == "Split" {
// 		isSplit = true
// 	}

// 	// Location
// 	outputRow[0] = inputRow["Channel Number"]
// 	// Name
// 	outputRow[1] = inputRow["Name"]
// 	// Frequency
// 	outputRow[2] = inputRow["Receive Frequency"]
// 	// Duplex - if not ham frequency, turn off transmitter
// 	if isHam == false {
// 		outputRow[3] = "off"			// Documented as turning off transmit.  will not import on Kenwood
// 	} else {
// 		if isSplit {
// 			outputRow[3] = "split"
// 		} else {
// 			outputRow[3] = offsetDirection[inputRow["Offset Direction"]]
// 		}
// 	}
// 	// Offset
// 	if isHam == false {
// 		outputRow[4] = "0.00000"
// 	} else {
// 		if isSplit {
// 			outputRow[4] = inputRow["Transmit Frequency"]
// 		} else {
// 			outputRow[4] = offsetValues[inputRow["Offset Frequency"]]
// 		}
// 	}

// 	// Tone
// 	outputRow[5] = toneValues[inputRow["Tone Mode"]]
// 	// rToneFreq
// 	outputRow[6] = strings.Fields(inputRow["CTCSS"])[0]
// 	// cToneFreq
// 	outputRow[7] = "88.5"
// 	// DtcsCode
// 	outputRow[8] = inputRow["DCS"]
// 	// DtcsPolarity
// 	outputRow[9] = "NN"
// 	// Mode
// 	outputRow[10] = "FM"
// 	// TStep
// 	outputRow[11] = "5.00"
// 	// Skip
// 	outputRow[12] = skipValues[inputRow["Skip"]]
// 	// Comment
// 	outputRow[13] = inputRow["Comment"]
// 	// URCALL
// 	outputRow[14] = ""
// 	// RPT1CALL
// 	outputRow[15] = ""
// 	// RPT2CALL
// 	outputRow[16] = ""

// 	return outputRow[:]
// }

// CSVFileToMap  reads csv file into slice of map
// slice is the line number
// map[string]string where key is column name
func CSVFileToMap(filePath string) (returnMap []map[string]string, err error) {

	// read csv file
	csvfile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	fmt.Printf("size is %d\n", len(rawCSVdata))

	header := []string{} // holds first row (header)
	for lineNum, record := range rawCSVdata {

		// for first row, build the header slice
		if lineNum == 0 {
			for i := 0; i < len(record); i++ {
				var columnName = strings.TrimSpace(record[i])
				if len(columnName) > 0 {
					header = append(header, columnName)
					fmt.Printf("header is /%s/ length is %d\n", columnName, len(columnName))
				}
			}
		} else {
			fmt.Printf("linenumber is %d\n", lineNum)
			// for each cell, map[string]string k=header v=value
			line := map[string]string{}
			for i := 0; i < len(record); i++ {
				if i < len(header) {
					line[header[i]] = record[i]
					fmt.Printf("header is %s value is %s\n", header[i], record[i])
				}
			}
			returnMap = append(returnMap, line)
		}
	}
	return
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
