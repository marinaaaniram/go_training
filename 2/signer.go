package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// func main() {
// 	inputData := []int{0, 1}

// 	hashSignJobs := []job{
// 		job(func(in, out chan interface{}) {
// 			for _, fibNum := range inputData {
// 				out <- fibNum
// 			}
// 		}),
// 		SingleHash,
// 		MultiHash,
// 		CombineResults,
// 		job(func(in, out chan interface{}) {
// 			dataRaw := <-in
// 			data, ok := dataRaw.(string)
// 			if !ok {
// 				fmt.Println("cant convert result data to string")
// 			}
// 			fmt.Println(data)
// 		}),
// 	}

// 	ExecutePipeline(hashSignJobs...)

// }

func ExecutePipeline(freeFlowJobs ...job) {

	channels_jobs := make([]chan interface{}, len(freeFlowJobs)+1)
	for i := range channels_jobs {
		channels_jobs[i] = make(chan interface{}, 2)
	}

	wg := &sync.WaitGroup{}
	for i, j := range freeFlowJobs {
		wg.Add(1)
		go func(wgg *sync.WaitGroup, j job, in, out chan interface{}) {
			defer wgg.Done()
			j(in, out)
			close(out)
		}(wg, j, channels_jobs[i], channels_jobs[i+1])
	}

	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	for inputData := range in {
		data := strconv.Itoa(inputData.(int))

		chan_step1 := make(chan string)
		chan_step2 := make(chan string)
		chan_step3 := make(chan string)

		go func(data string) {
			chan_step1 <- DataSignerMd5(data)
		}(data)

		go func(data string) {
			chan_step3 <- DataSignerCrc32(data)
		}(data)

		step1 := <-chan_step1
		close(chan_step1)

		go func(data string) {
			chan_step2 <- DataSignerCrc32(data)
		}(step1)

		step3 := <-chan_step3
		close(chan_step3)

		step2 := <-chan_step2
		close(chan_step2)

		fmt.Printf("%s SingleHash md5(data) %s\n", data, step1)
		fmt.Printf("%s SingleHash crc32(data) %s\n", data, step3)
		fmt.Printf("%s SingleHash crc32(md5(data)) %s\n", data, step2)

		step4 := step3 + "~" + step2
		fmt.Printf("%s SingleHash result %s\n", data, step4)
		out <- step4
	}
}

func MultiHash(in, out chan interface{}) {
	thLen := 6

	for inputData := range in {
		data, ok := inputData.(string)
		if !ok {
			panic("Can't convert result data to string")
		}

		result := make([]string, thLen)
		wg := &sync.WaitGroup{}
		for i := 0; i < thLen; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, i int, data string) {
				defer wg.Done()
				result[i] = DataSignerCrc32(strconv.Itoa(i) + data)
				fmt.Printf("%s MultiHash: crc32(th+step1) %d %s\n", data, i, result[i])
			}(wg, i, data)
		}

		wg.Wait()

		splitResult := strings.Join(result, "")
		fmt.Printf("%s MultiHash result %s\n", data, splitResult)
		out <- splitResult

	}
}

func CombineResults(in, out chan interface{}) {
	var resultSlice []string

	for i := range in {
		data, ok := i.(string)
		if !ok {
			panic("Can't convert result data to string")
		}
		resultSlice = append(resultSlice, data)
	}

	sort.Strings(resultSlice)

	result := strings.Join(resultSlice, "_")

	fmt.Printf("CombineResults %s\n", result)

	out <- result

}
