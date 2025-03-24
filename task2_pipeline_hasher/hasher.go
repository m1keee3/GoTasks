package main

import (
	"fmt"
	"strconv"
	"sync"
)

func main() {

	inputData := []int{0, 1}

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			dataRaw, ok := <-in
			data, ok := dataRaw.(string)
			if !ok {
				fmt.Println("cant convert result data to string")
			}
			fmt.Println(data)
		}),
	}

	ExecutePipeline(hashSignJobs...)
}

func ExecutePipeline(jobs ...job) {

	chans := make([]chan interface{}, len(jobs)-1)

	wg := sync.WaitGroup{}
	wg.Add(len(jobs))

	for i, job := range jobs {

		if i == 0 {
			chans[i] = make(chan interface{})
			go func() {
				defer wg.Done()
				job(nil, chans[i])
				close(chans[i])
			}()

		} else if i == len(jobs)-1 {
			go func() {
				defer wg.Done()
				job(chans[i-1], nil)
			}()
		} else {
			chans[i] = make(chan interface{})
			go func() {
				defer wg.Done()
				job(chans[i-1], chans[i])
				close(chans[i])
			}()
		}
	}

	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	mu := &sync.Mutex{}
	wg := sync.WaitGroup{}
	for dataRaw := range in {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, ok := dataRaw.(int)
			if !ok {
				fmt.Println("dataRaw is not a string, SingleHash")
			}
			data := strconv.Itoa(dataRaw.(int))

			crcChan1 := make(chan string)
			crcChan2 := make(chan string)
			go func() {
				crcChan1 <- DataSignerCrc32(data)
			}()
			go func() {
				mu.Lock()
				mdData := DataSignerMd5(data)
				mu.Unlock()
				crcChan2 <- DataSignerCrc32(mdData)
			}()
			crc1 := <-crcChan1
			crc2 := <-crcChan2
			fmt.Println(crc1 + "~" + crc2)
			out <- crc1 + "~" + crc2
		}()
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	for dataRaw := range in {
		wg.Add(1)
		go func() {
			defer wg.Done()
			wg := &sync.WaitGroup{}
			wg.Add(6)
			data, ok := dataRaw.(string)
			if !ok {
				fmt.Println("dataRaw is not a string, MultiHash")
			}
			resHashs := make([]string, 6)
			for i := 0; i < 6; i++ {
				go func() {
					defer wg.Done()
					resHashs[i] = DataSignerCrc32(strconv.Itoa(i) + data)
				}()
			}
			wg.Wait()

			var res string
			for _, hash := range resHashs {
				res += hash
			}
			out <- res
		}()
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var res string
	for dataRow1 := range in {
		data, ok := dataRow1.(string)
		if !ok {
			fmt.Println("dataRow1 is not a string, CombineResults")
		}

		if res != "" {
			res += "_"
		}
		res += data
	}
	out <- res
}
