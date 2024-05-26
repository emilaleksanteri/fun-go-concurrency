package workerpool

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

type Work struct {
	file   string
	result chan Result
}

type Result struct {
	file string
	text string
}

type WorkerPool struct {
	poolSize int
	fileDir  string
}

func NewWorkerPool(poolSize int, fileDir string) WorkerPool {
	return WorkerPool{
		poolSize: poolSize,
		fileDir:  fileDir,
	}
}

func (wp *WorkerPool) DoPool() {
	jobs := make(chan Work)
	allResults := make(chan chan Result)
	wg := sync.WaitGroup{}
	for i := 0; i < wp.poolSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(jobs)
		}()
	}
	go func() {
		defer close(allResults)
		filepath.Walk(wp.fileDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				ch := make(chan Result)
				jobs <- Work{file: path, result: ch}
				allResults <- ch
			}
			return nil
		})

	}()

	for resultCh := range allResults {
		for result := range resultCh {
			fmt.Printf("%s result: %s\n", result.file, result.text)
		}
	}
	close(jobs)
	wg.Wait()

}

func worker(jobs chan Work) {
	for work := range jobs {
		file, err := os.Open(work.file)
		if err != nil {
			fmt.Printf("err reading file: %v\n", err)
			continue
		}
		defer file.Close()

		scan := bufio.NewScanner(file)
		lineNum := 0
		content := []byte{}
		for scan.Scan() {
			content = append(content, scan.Bytes()...)
			lineNum++
		}
		work.result <- Result{
			file: work.file,
			text: string(content),
		}
		close(work.result)
	}
}
