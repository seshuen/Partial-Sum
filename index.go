/*
* Name : Seshagopalan Narasimhan
* Description: GoLang multithreaded applications to compute the sum of integers stored in a file.
* Input: M and fname. M is an integer and fname is the pathname (relative or absolute) to the input data file.
* Output:  Sum
*/

package main

import (
	"os"
	"strconv"
	"math"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// *********************** All Struct ********************************
type worker struct {
	id int
	job chan []byte
	result chan []byte
}

type message struct {
	Fname string
	Start int
	End int
}

type result struct{
	Value uint64
	Prefix string
	Suffix string
	Error string
}

// ************************* Other Functions ***************************************
func checkError(e error){
	if e != nil{
		panic(e)
	}
}

func cleanChunk(chunk string) string{
	//Check for trailing whitespace  and Replace all whitespace by space
	ts := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	chunk_ := ts.ReplaceAllString(chunk, " ")

	//Check for trailing whitespace and Replace all whitespace by space
	ds := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	chunk_ = ds.ReplaceAllString(chunk_, " ")

	return chunk_
}


// ************************************* Co-ordinator and Worker ***********************************************
func coordinator(fname string, M float64){
	var resp = make([][]byte, int(M))
	var msg *message
	var frg map[string] interface{}
	var temp = ""
	var continue_ bool = false

	wg := new(sync.WaitGroup)

	job_ := make(chan []byte, int(M) * 10)
	result_ := make(chan []byte, int(M) * 10)

	sum := uint64(0)


//	open the file
	read, err := os.Open(fname)
	checkError(err)

	defer read.Close()

// compute the file size
	fs, err :=  read.Stat()

// Divide into M chunks
	frag := math.Floor(float64(fs.Size()-1) / M)

//	Get the last extra chunks to add to the last thread
	last := int64(fs.Size()) - ( int64(frag) * int64(M) )

// Call the thread creator and pass the filename with starting and ending byte
	fsize := 0

	fmt.Println("Total File Size: ", fs.Size())
	for i:=0; i < int(M); i++ {


		if i == int(M)-1 {
			msg = &message{
				fname,
				fsize,fsize + int(frag) + int(last) - 1,
			}
		} else {
			msg = &message{
				fname,
				fsize,
				fsize + int(frag) - 1,
			}
		}

		pckFile,_ := json.Marshal(msg)
		fmt.Printf("Following JSON sent to Worker %d: %s\n", i, pckFile)
		job_ <- pckFile

		wrk := &worker{i, job_ , result_}

		go wrk.compute(wg)
		wg.Add(1)

		resp[i] = <-result_
		fsize += int(frag)
	}

	for index, frag := range resp{
		err := json.Unmarshal(frag, &frg)
		checkError(err)
		fmt.Printf("Worker %d responded with: Value(%v), Prefix(%v), Suffix(%v), Error(%v)\n",index,frg["Value"],frg["Prefix"],frg["Suffix"],frg["Error"] )

		if(frg["Error"].(string) == ""){
			if(frg["Suffix"].(string) != "" && frg["Prefix"].(string) != ""){
				//fmt.Println("CASE 1")
				temp = fmt.Sprintf("%s%s", temp, frg["Prefix"].(string))
				newValue,_ := strconv.ParseInt(temp,10,64)
				sum = sum + uint64(newValue)
				temp = frg["Suffix"].(string)
				continue_ = true

			} else if(frg["Suffix"].(string) != ""){
				fmt.Println("CASE 2", continue_)
				if(continue_ == false){
					temp = frg["Suffix"].(string)
					continue_ = true
				} else{
					a := fmt.Sprintf("%s%s", frg["Prefix"].(string), temp)
					res,_ := strconv.ParseUint(a, 10, 64)
					sum += res
					temp = frg["Suffix"].(string)
				}
			} else if(frg["Prefix"].(string) != ""){

				//fmt.Println("CASE 3", continue_)
				if(continue_ == false){

					newValue,_ := strconv.ParseInt(frg["Prefix"].(string),10,64)
					sum = sum + uint64(newValue)

				} else{
					temp = fmt.Sprintf("%s%s", temp, frg["Prefix"].(string))
					newValue,_ := strconv.ParseInt(temp,10,64)
					sum = sum + uint64(newValue)
					temp = ""
					continue_ = false
				}
			} else{
				continue_ = false
			}

		} else{
			if(continue_ == false) {

				temp = frg["Error"].(string)
				continue_ = true

			}else{
				temp = fmt.Sprintf("%s%s", temp, frg["Error"].(string))
				//fmt.Println("Chunk Merged", temp)
			}

		}

		//Compute Sum
		if(temp != "" && continue_ == false) {
			newValue, _ := strconv.ParseInt(temp, 10, 64)
			fmt.Printf("When Temp-> sum + value + newvalue: %d + %d + %d\n", sum, uint64(frg["Value"].(float64)),uint64(newValue))
			sum = sum + uint64(frg["Value"].(float64)) +  uint64(newValue)
			temp = ""
			fmt.Printf(" sum + value: %d \n", sum)
		} else {
			fmt.Printf("When no Temp-> sum + value: %d + %d\n", sum, uint64(frg["Value"].(float64)))
			sum = sum + uint64(frg["Value"].(float64))
			fmt.Printf(" sum + value: %d \n", sum)
		}


	}//end of loop
	er := json.Unmarshal(resp[len(resp)-1], &frg)
	checkError(er)

	//lastIndex := len(resp) - 1
	temp = frg["Suffix"].(string)
	fmt.Println("Last Suffix", temp)
	newValue, _ := strconv.ParseUint(temp, 10, 64)
	sum = sum +  uint64(newValue)

	fmt.Printf("Total Sum = %d", sum)
	wg.Wait()

}

func (wrkr *worker) compute(wg *sync.WaitGroup){
	var msg_ map[string]interface{}
	value := uint64(0)
	error := ""
	prefix := ""
	suffix := ""

	//	get the chunk
	err := json.Unmarshal(<-wrkr.job, &msg_)
	checkError(err)

	//get the individual attributes from the fragment
	fname := msg_["Fname"].(string)
	start := msg_["Start"].(float64)
	end := msg_["End"].(float64)

	//create chunk
	read, err := os.Open(fname)
	checkError(err)

	defer read.Close()

	_, e_ := read.Seek(int64(start),0 )
	checkError(e_)

	chunk := make([]byte, byte(end - start + 1))

	_, e2 := read.Read(chunk)
	checkError(e2)

	chunk_ := string(chunk)
	fmt.Printf("[%s]\n",chunk_)
	chunkStr := cleanChunk(chunk_)

	if(strings.Contains(chunkStr, " ")){
		psum := uint64(0)
		chunks := strings.Split(string(chunkStr), " ")
		// if no starting or ending space send it as preffix or suffix
		prefix = chunks[0]
		suffix = chunks[len(chunks) - 1]

		// else calculate sum
		for _, element := range chunks[1:len(chunks) - 1]{

			val,_ := strconv.ParseInt(element, 10,64)
			psum += uint64(val)
		}

		// get partial sum
		value = psum

	} else{
		//	-- If no partial sum return in exception of json file
		error = chunkStr
	}

	res := &result{
		value,
		prefix,
		suffix,
		error,
	}
	pckRes,_ := json.Marshal(res)

	wg.Done()

	wrkr.result <- pckRes


}

// ************************************* Main() ***************************************************************
func main(){
	//Get the inputs
	M := os.Args[1]
	fname := os.Args[2]

	//Convert string to float
	frag, _ := strconv.ParseFloat(M, 64)

	//Pass the values to cordinator
	coordinator(fname, frag)

	//	Display the Total

}
