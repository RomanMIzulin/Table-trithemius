package main

import (
	"github.com/andlabs/ui"
	"fmt"
	"io/ioutil"
	"errors"
	"encoding/csv"
	"bytes"
	"io"
	"log"
	"os"
	"bufio"
	"strconv"
)

func main() {
	err := ui.Main(func(){

		boxH:=ui.NewHorizontalBox()
		box := ui.NewVerticalBox()

		combo := ui.NewCombobox()
		combo.Append("rus")
		combo.Append("eng")
		combo.SetSelected(0)
		// 0 is rus 1 is eng

		txtEnEntry := ui.NewEntry()
		buttonEn := ui.NewButton("Encrypt")
		output := ui.NewLabel(" ")
		PathKey := ui.NewEntry()
		txtEnEntry.SetText("Text to encode")
		PathKey.SetText("Path to key")
		box.Append(txtEnEntry,true)
		box.Append(PathKey,true)
		PathText:=ui.NewEntry()
		box.Append(PathText,true)
		box.Append(combo,true)

		box.Append(boxH,true)
		boxH.Append(buttonEn,true)
		butDecode:= ui.NewButton("Decrypt")
		boxH.Append(butDecode,true)

		PathText.SetText("path to text")


		statusKeyBar:=ui.NewLabel("output")
		box.Append(statusKeyBar,false)
		box.Append(output,true)

		pathMultiFile:=ui.NewEntry()
		pathMultiFile.SetText("txt_key_eng.csv")
		butForEnMul:=ui.NewButton("multi encrypt")
		butForDeMul:=ui.NewButton("multi decrypt")
		saveButton:=ui.NewButton("Save")
		box.Append(saveButton,true)
		box.Append(pathMultiFile,true)

		boxHBottom:=ui.NewHorizontalBox()
		box.Append(boxHBottom,true)
		boxHBottom.Append(butForEnMul,true)
		boxHBottom.Append(butForDeMul,true)

		//reference := ui.NewLabel(" ")

		window := ui.NewWindow("Table trithemius ", 100, 100, false)
		window.SetMargined(true)
		window.SetChild(box)

		//  /home/romanmatveev/goui/keyfile.txt
		buttonEn.OnClicked(func(*ui.Button) {
			output.SetText(string(*encrypt(txtEnEntry,PathKey,PathText,combo.Selected())))
			tmpText,err:=readFromFile((PathKey.Text()))
			if err == nil{
				statusKeyBar.SetText("encrypting with key: "+string(tmpText))
			}else{
				statusKeyBar.SetText("encrypting with default key")
			}
		})

		butDecode.OnClicked(func(*ui.Button){
			output.SetText(string(*textToDecode(txtEnEntry,PathKey,PathText,combo.Selected())))
			tmpText,err:=readFromFile((PathKey.Text()))
			if err == nil{
				statusKeyBar.SetText("encrypting with key: "+string(tmpText))
			}else{
				statusKeyBar.SetText("encrypting with default key")
			}
		})

		saveButton.OnClicked(func(button *ui.Button) {
			if output.Text()!= ""{
				ioutil.WriteFile("output.txt",[]byte(output.Text()),0777)
			}
		})

		butForEnMul.OnClicked(func(button2 *ui.Button ){
			counter := 0

			fi,err := os.Open("txt_key.csv")
			if err !=nil {
				log.Println("failed to open txt_key_eng.csv")
			}
			defer fi.Close()
			r:=bufio.NewReader(fi)
			r.ReadLine()

			rcsv := csv.NewReader(r)
			whoa,err:=os.Create("outputMulti.txt")
			defer whoa.Close()

			statusKeyBar.SetText(strconv.Itoa(counter))
			buf := make([]byte,1664)

			for {
				rec, err := rcsv.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatal(err)
				}

				txt:=([]byte(rec[0]))
				rtxt:=bufio.NewReader(bytes.NewReader(txt))
				key:=([]byte(rec[1]))
				for {
					n, err := rtxt.Read(buf)
					log.Println(n,buf[:n])
					if err != nil && err != io.EOF {
						panic(err)
					}
					if n == 0 {
						break
					}
					bf:= buf[:n]
					bb:=&bf
					if _, err := whoa.Write((encryptForMulti(bb, key, combo.Selected()))); err != nil {
						log.Println("wtf")
					}
				}
				whoa.Write([]byte("\n"))
				counter++
				statusKeyBar.SetText(strconv.Itoa(counter))
			}

		})

		butForDeMul.OnClicked(func(button2 *ui.Button ){
			counter := 0

			fi,err := os.Open("key_txt.csv")
			if err !=nil {
				log.Println("failed to open key_txt.csv")
			}
			defer fi.Close()
			r:=bufio.NewReader(fi)
			rcsv := csv.NewReader(r)

			statusKeyBar.SetText(strconv.Itoa(counter))
			buf := make([]byte,1664)

			outM,err:= os.Create("outputMulti.txt")
			if err != nil{
				log.Println("failed to open outputMulti.txt")
			}
			defer outM.Close()

			for {
				rec, err := rcsv.Read()

				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatal(err)
				}

				txt := ([]byte(rec[1]))
				log.Println(rec[1])
				rtxt := bufio.NewReader(bytes.NewReader(txt))
				key := ([]byte(rec[0]))
				for {
					n, err := rtxt.Read(buf)
					log.Println(n, buf[:n])
					if err != nil && err != io.EOF {
						panic(err)
					}
					if n == 0 {
						break
					}
					bf:= buf[:n]
					bb:=&bf
					if _, err := outM.Write((decryptForMulti(bb, key, combo.Selected()))); err != nil {
						log.Println("wtf")
					}
				}
				outM.Write([]byte("\n"))
				counter++
				statusKeyBar.SetText(strconv.Itoa(counter))
			}

		})

		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()

	})
	if err != nil {
		panic(err)
		}
}


func readFromFile(Path string) ([]byte,error) {
	key, err := ioutil.ReadFile(Path)
	if err != nil {
		fmt.Println("No such file")
		return key,errors.New("No")
	}
	return key,nil
}

func encrypt(userInput *ui.Entry,path_to_key *ui.Entry,pathToText *ui.Entry,mode int)  *[]byte {
	textToEncode := userInput.Text()
	if txt,err:=readFromFile(pathToText.Text());err==nil{
		textToEncode=string(txt)
	}
	lenOfText := len(textToEncode)
	keyTable, err := readFromFile(path_to_key.Text())
	textToEncodeBytes := []byte(textToEncode)
	if err != nil {
		//default encrypting
		if mode == 0 {

			for i := 1; i < lenOfText; i+=2 {
				getNumberOfRussianLetter(textToEncodeBytes[i])
				textToEncodeBytes[i],textToEncodeBytes[i-1]=getRussianCode(getNumberOfRussianLetter(textToEncodeBytes[i])+i/2)
			}
			return &textToEncodeBytes
		}
		if mode == 1 {

			for i := 0; i < lenOfText; i++ {
				if textToEncodeBytes[i]==0x7a{
					textToEncodeBytes[i]=byte(0x61)
				}

				textToEncodeBytes[i] += byte(i % 26)

			}
			return &textToEncodeBytes
		}
	} else {
		// /home/romanmatveev/goui/keyfile.txt
		if mode == 0 {
			fmt.Println(textToEncodeBytes)
			for k := 1; k < lenOfText; k+=2 {

					fmt.Println(k/2)


					if 175 < (textToEncodeBytes[k]) && (textToEncodeBytes[k]) < 192 {
						textToEncodeBytes[k] = keyTable[(k+2*(int(textToEncodeBytes[k]) - 176)) % 32 ]
						fmt.Println(textToEncodeBytes[k])
					}
					if 127 < (textToEncodeBytes[k]) && (textToEncodeBytes[k]) < 144 {
						if 175<keyTable[(k+2*(int(textToEncodeBytes[k])-128))%32 ]&& keyTable[(k+2*(int(textToEncodeBytes[k])-128))%32]<192 {
							textToEncodeBytes[k] = keyTable[(k+2*(int(textToEncodeBytes[k])-128))%32]
							textToEncodeBytes[k-1] = 208
						}else{
							textToEncodeBytes[k] = keyTable[(k+2*(int(textToEncodeBytes[k])-128))%32]
						}
						fmt.Println(textToEncodeBytes[k])
					}


			}
		return &textToEncodeBytes
		}
		if mode == 1 {
			for k := 0; k < lenOfText; k++ {

					if 96 < (textToEncodeBytes[k]) && (textToEncodeBytes[k]) < 123 {
						textToEncodeBytes[k] = keyTable[(k+(int(textToEncodeBytes[k])- 97))%26]

					}

			}
		return &textToEncodeBytes
		}

		}
		return &textToEncodeBytes
	}

// /home/romanmatveev/goui/keyfile.txt
// /home/romanmatveev/goui/deCode.txt

func textToDecode(userInput *ui.Entry,pathToKey *ui.Entry,pathToText *ui.Entry, mode int)  *[]byte {
	textToDecode := userInput.Text()
	textToDecodeBytes := []byte(textToDecode)
	text,err := readFromFile(pathToText.Text())
	if err == nil{
		textToDecodeBytes = text
	}
	lenOfText:= len(textToDecodeBytes)
	key, err := readFromFile(pathToKey.Text())
	keyBytes:= []byte(key)
	if err == nil{
		if mode ==0 {
			for i:=1;i<lenOfText;i+=2{
				if i!=1 {
					tmp := keyBytes[0:2]
					keyBytes = keyBytes[2:]
					keyBytes = append(keyBytes, tmp[0], tmp[1])
				}
				count:=0
				for j:=1;j<64 ;j+=2{

					if key[j]==textToDecodeBytes[i]{
						textToDecodeBytes[i],textToDecodeBytes[i-1] = getRussianCode(count/2)
						fmt.Println(getRussianCode(count))
					}
					count=count+1 % 32
				}
			}
			return &textToDecodeBytes
		}
		if mode ==1 {
			for i:=0;i<lenOfText;i++{
				if i!=0 {
					tmp := keyBytes[0]
					keyBytes = keyBytes[1:]
					keyBytes = append(keyBytes, tmp)
				}
				fmt.Println(keyBytes)
				count:=0
				for j:=0;j<26 ;j++{

					if keyBytes[j]==textToDecodeBytes[i]{
						textToDecodeBytes[i]=byte(0x61+count)

					}
					count=count+1 % 32
				}
			}
			fmt.Println(err)
		}
		}else{

		if mode ==0 {
			for i:=1;i<lenOfText;i+=2{
				textToDecodeBytes[i],textToDecodeBytes[i-1]=getRussianCode(i/2)
			}
		}
		if mode ==1 {
			for i := 0; i < lenOfText; i++ {
				if i%26!=0 && textToDecodeBytes[i]==0x61 {
					textToDecodeBytes[i]=byte(0x7a)
				}

				textToDecodeBytes[i] -= byte(i % 26)
			}
			return &textToDecodeBytes

		}
		}

	return &textToDecodeBytes
}

func getRussianCode(i int ) (byte,byte)  {
	i%=32
	if i< 16{
		return byte(176+i), 208
	}else{
		i%=16
		return byte(128+i), 209
	}

}
func getNumberOfRussianLetter(b byte) int{
	if b<176{
		return int(b)-128
	}else{
		return int(b)-176
	}
}

func encryptForMulti(txt *[]byte,key []byte,mode int) ([]byte){
	lenOfText:=len(*txt)
	text:=*txt
	keys:=key

// /home/romanmatveev/goui/keyfile.txt
	if mode == 0 {
		for k := 1; k < lenOfText; k+=2 {

		fmt.Println(k/2)

		if 175 < (text[k]) && (text[k]) < 192 {
		text[k] = keys[(k+2*(int(text[k]) - 176)) % 32 ]
		fmt.Println(text[k])
		}
		if 127 < (text[k]) && (text[k]) < 144 {
		if 175<keys[(k+2*(int(text[k])-128))%32 ]&& keys[(k+2*(int(text[k])-128))%32]<192 {
		text[k] = keys[(k+2*(int(text[k])-128))%32]
		text[k-1] = 208
		}else{
		text[k] = keys[(k+2*(int(text[k])-128))%32]
		}
		fmt.Println(text[k])
		}

		}
		return text
	}
	if mode == 1 {
		for k := 0; k < lenOfText; k++ {


		if 96 < (text[k]) && (text[k]) < 123 {
		text[k] = keys[(k+(int(text[k])- 97) ) % 26]

		}

		}
		return text
	}
	return text
}

func decryptForMulti(txt *[]byte,key []byte,mode int) ([]byte){
	text:=*txt
	lenOfText:= len(text)

		if mode ==0 {
			for i:=1;i<lenOfText;i+=2{
				if i!=1 {
					tmp := key[0:2]
					key = key[2:]
					key = append(key, tmp[0], tmp[1])
				}
				count:=0
				for j:=1;j<64 ;j+=2{

					if key[j]==text[i]{
						text[i],text[i-1] = getRussianCode(count/2)

						fmt.Println(getRussianCode(count))
					}
					count=count+1 % 32
				}
			}
			return text
		}
		if mode ==1 {
			for i:=0;i<lenOfText;i++{
				if i!=0 {
					tmp := key[0]
					key = key[1:]
					key = append(key, tmp)
				}
				fmt.Println(key)
				count:=0
				for j:=0;j<26 ;j++{

					if key[j]==text[i]{
						text[i]=byte(0x61+j )
						break
					}
					count=count+1 % 26
				}
			}

		}
	return text
}
