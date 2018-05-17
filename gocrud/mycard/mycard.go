package mycard

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/ebfe/scard"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
	//"github.com/varokas/tis620"
	"unicode/utf8"
)

const OFFSET = 0xd60
const WIDTH = 3

type Person struct {
	ID        string   `json:"id"`
	THprefix  string   `json:"thprefix"`
	THfname   string   `json:"thfname"`
	THlname   string   `json:"thlname"`
	ENprefix  string   `json:"enprefix"`
	ENfname   string   `json:"enfname"`
	ENlname   string   `json:"enlname"`
	Addr      []string `json:"addr"`
	Birthdate string   `json:"birthdate"`
	Age       int      `json:"age"`
	Sex       string   `json:"sex"`
}

func ToUTF8(tis620bytes []byte) []byte {
	l := findOutputLength(tis620bytes)
	output := make([]byte, l)

	index := 0
	buffer := make([]byte, WIDTH)
	for _, c := range tis620bytes {
		if !isThaiChar(c) {
			output[index] = c

			index++
		} else {
			utf8.EncodeRune(buffer, int32(c)+OFFSET)
			output[index] = buffer[0]
			output[index+1] = buffer[1]
			output[index+2] = buffer[2]

			index += 3
		}
	}
	return output
}

func findOutputLength(tis620bytes []byte) int {
	outputLen := 0
	for i, _ := range tis620bytes {
		if isThaiChar(tis620bytes[i]) {
			outputLen += WIDTH //always 3 bytes for thai char
		} else {
			outputLen += 1
		}
	}
	return outputLen
}

func isThaiChar(c byte) bool {
	return (c >= 0xA1 && c <= 0xDA) || (c >= 0xDF && c <= 0xFB)
}
func Decode(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, traditionalchinese.Big5.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}
func CToGoString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

/* func GetCardData(p *person){
	_thname := split(p.name[0:100],"#")
} */
func getAge(b string) int {
	var year, month, day int
	//year, err := strconv.ParseInt(b[0:4], 10, 0)
	//year = int16(year) - 543
	year, _ = strconv.Atoi(b[0:4])
	year = year - 543
	month, _ = strconv.Atoi(b[4:6])
	day, _ = strconv.Atoi(b[6:8])
	birthday := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	now := time.Now()
	age := now.Year() - birthday.Year()
	if now.YearDay() < birthday.YearDay() {
		age--
	}
	fmt.Println(year, month, day)
	return age
}
func ReadCard() Person {
	// Establish a PC/SC context
	//var p Person
	context, err := scard.EstablishContext()
	if err != nil {
		fmt.Println("Error EstablishContext:", err)

	}

	// Release the PC/SC context (when needed)
	defer context.Release()

	// List available readers
	readers, err := context.ListReaders()
	if err != nil {
		fmt.Println("Error ListReaders:", err)

	}

	// Use the first reader
	reader := readers[0]
	fmt.Println("Using reader:", reader)

	// Connect to the card
	card, err := context.Connect(reader, scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		fmt.Println("Error Connect:", err)

	}

	// Disconnect (when needed)
	defer card.Disconnect(scard.LeaveCard)

	// Send select APDU
	var cmd_select = []byte{0x00, 0xA4, 0x04, 0x00, 0x08, 0xA0, 0x00, 0x00, 0x00, 0x54, 0x48, 0x00, 0x01}
	rsp, err := card.Transmit(cmd_select)
	if err != nil {
		fmt.Println("Error Transmit:", err)

	}
	fmt.Println(rsp)

	// Send command APDU
	var cmd_name1 = []byte{0x80, 0xb0, 0x00, 0x11, 0x02, 0x00, 0xd1}
	var cmd_name2 = []byte{0x00, 0xc0, 0x00, 0x00, 0xd1}

	var cmd_id1 = []byte{0x80, 0xb0, 0x00, 0x04, 0x02, 0x00, 0x0d}
	var cmd_id2 = []byte{0x00, 0xc0, 0x00, 0x00, 0x0d}

	var cmd_addr1 = []byte{0x80, 0xb0, 0x15, 0x79, 0x02, 0x00, 0x64}
	var cmd_addr2 = []byte{0x00, 0xc0, 0x00, 0x00, 0x64}
	//get id
	rsp, err = card.Transmit(cmd_id1)
	if err != nil {
		fmt.Println("Error Transmit:", err)

	}

	rsp, err = card.Transmit(cmd_id2)
	if err != nil {
		fmt.Println("Error Transmit:", err)

	}
	id := ToUTF8(rsp)
	//fmt.Println("ID:", id)
	//get name
	rsp, err = card.Transmit(cmd_name1)
	if err != nil {
		fmt.Println("Error Transmit:", err)

	}
	//fmt.Println(rsp)
	rsp, err = card.Transmit(cmd_name2)
	if err != nil {
		fmt.Println("Error Transmit:", err)

	}
	name := ToUTF8(rsp)
	//fmt.Println("name:", name)

	//-----------------get addr
	rsp, err = card.Transmit(cmd_addr1)
	if err != nil {
		fmt.Println("Error Transmit:", err)

	}
	//fmt.Println(rsp)
	rsp, err = card.Transmit(cmd_addr2)
	if err != nil {
		fmt.Println("Error Transmit:", err)

	}
	addr := ToUTF8(rsp)
	//fmt.Println("addr:", addr)
	//fmt.Println(rsp)
	//tt := []byte(rsp)
	//fmt.Println(tt)
	//t,err := Decode(rsp)
	//fmt.Println(string(b))
	//fmt.Fprintf(w, "<h1>%s</h1>",tt)
	//fmt.Println(string(b))
	//fmt.Println(err)
	//str := CToGoString(rsp[:])

	//x := ToUTF8(rsp)
	_id := strings.TrimSpace(string(id))

	thname := strings.Split(string(name)[0:100], "#")
	_th_prefix := strings.TrimSpace(thname[0])
	_th_fname := strings.TrimSpace(thname[1])
	_th_lname := strings.TrimSpace(thname[3])

	enname := strings.Split(string(name)[100:200], "#")
	_en_prefix := strings.TrimSpace(enname[0])
	_en_fname := strings.TrimSpace(enname[1])
	_en_lname := strings.TrimSpace(enname[3])
	//fmt.Println("name len:%d", len(thname))
	//_name = string(name)[0:100]
	_addr := strings.Split(string(addr), "#")
	/*
		[0]house no.
		[1]village no.addr
		[2]lane
		[3]road
		[4]
		[5]tambol
		[6]amphur
		[7]province

	*/
	b := string(name)

	b = strings.TrimSpace(b[200:247])
	fmt.Println("len:", len(b))
	fmt.Println("birthday:", b)
	birth := b[0:8]
	_sex := b[8:9]

	fmt.Println("sex:", _sex)
	_age := getAge(birth)
	fmt.Println("age:", _age)
	p := Person{_id, _th_prefix, _th_fname, _th_lname, _en_prefix, _en_fname, _en_lname, _addr, birth, _age, _sex}
	//fmt.Println(p)
	//a := Person{"a", "b", "c", "d", "e", "g", "g", "h", "i", "10", "1"}
	return p
	//fmt.Println(string(x))
	//fmt.Fprintf(w, "<h1>%s</h1>", x)

	//fmt.Println()
}
