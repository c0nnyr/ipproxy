package common

import (
	"log"
	"runtime"
    "time"
)

func ParrallelRun(job_func func(chan<- int, int), start_ind int) {
	PARRAL_COUNT := runtime.NumCPU()
	ch := make(chan int, PARRAL_COUNT)

	working_tasks := 0
	all_works_is_done := false
	ind := start_ind
	job_func_with_pack := func() {
		if all_works_is_done {
			return
		}
		working_tasks++
		ind++
		go job_func(ch, ind)
	}
	for i := 0; i < PARRAL_COUNT; i++ {
		job_func_with_pack()
	}
	for {
		select {
		case status := <-ch:
			working_tasks--
			if status == 0 {
				all_works_is_done = true
			} else {
				job_func_with_pack()
			}
		}
		if all_works_is_done && working_tasks == 0 {
			break
		}
	}
	log.Println("All are done")
}

func run_timeout(ch chan interface{}, param interface{}) {
	time.Sleep(60e9) //60s
	ch <- param
}
