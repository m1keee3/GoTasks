package main

import (
	"fmt"
	"sync"
)

func main() {

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			out <- 1
		}),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, ok := dataRaw.(int)
			if !ok {
				fmt.Println("dataRaw is not a string")
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
	dataRaw := <-in
	data, ok := dataRaw.(string)
	var crcData, crcData2, mdData string
	if !ok {
		fmt.Println("dataRaw is not a string")
	}
	go func() {
		crcData = DataSignerCrc32(data)
	}()
	go func() {
		mdData = DataSignerMd5(data)
	}()
	go func() {
		crcData2 = DataSignerCrc32(mdData)
	}()

	out <- crcData + "~" + crcData2
}

func MultiHash(in, out chan interface{}) {
	dataRaw := <-in
	data, ok := dataRaw.(string)
	if !ok {
		fmt.Println("dataRaw is not a string")
	}
	resHashes := make([]string, 6)
	for i := 0; i < 6; i++ {
		go func() {
			resHashes[i] = DataSignerCrc32(string(rune(i)) + data)
		}()
	}

	var res string

	for _, hash := range resHashes {
		res += hash
	}
	out <- res

}

func CombineResults(in, out chan interface{}) {

}
