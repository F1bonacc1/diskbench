package main

import (
	"bufio"
	"crypto/rand"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

type arguments struct {
	dirPath      string
	readers      int
	filesToWrite int
	fileSize     int
}

func getArgs() arguments {
	args := arguments{}
	flag.IntVar(&args.readers, "readers", 1, "amount of reader threads")
	flag.IntVar(&args.filesToWrite, "files", 0, "amount of files to write")
	flag.IntVar(&args.fileSize, "size", 10, "file size in MB (to write)")
	flag.StringVar(&args.dirPath, "dir", ".", "directory path to read")
	flag.Parse()

	return args
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// Benchmark provides disk performance measurement capabilities
type Benchmark struct {
	readDuration  float64
	writeDuration float64
	dataInBytes   int64
	dirPath       string
	filesToWrite  int
	filesWritten  int
	fileSizeMB    int
	randomData    []byte
}

func (b Benchmark) readFile(filename string) (int64, error) {
	file, err := os.Open(filename)

	if err != nil {
		return 0, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return 0, statsErr
	}

	size := stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	return size, err
}

func (b *Benchmark) writeFile(filename string) error {
	err := ioutil.WriteFile(filename, b.randomData, 0644)
	return err
}

func (b *Benchmark) iterrateDir(dirPath string) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	// defer timeTrack(time.Now(), "read")
	b.dataInBytes = 0
	start := time.Now()
	for i, file := range files {
		fileFullPath := path.Join(dirPath, file.Name())
		size, _ := b.readFile(fileFullPath)
		b.dataInBytes += size
		if i%5 == 0 {
			b.readDuration = time.Since(start).Seconds()
			b.printReadResults()
		}
	}

	b.readDuration = time.Since(start).Seconds()
}

func (b *Benchmark) generateFiles(dirPath string, filesToWrite, fileSizeMB int) {
	b.randomData = make([]byte, fileSizeMB*1024*1024)
	rand.Read(b.randomData)
	start := time.Now()
	for i := 1; i <= filesToWrite; i++ {
		fileName := fmt.Sprintf("file_%06d.dat", i)
		fileFullPath := path.Join(dirPath, fileName)
		b.writeFile(fileFullPath)
		b.filesWritten = i
		if i%5 == 0 {
			b.writeDuration = time.Since(start).Seconds()
			b.printWriteResults()
		}
	}

	b.writeDuration = time.Since(start).Seconds()
}

func (b *Benchmark) run(wg *sync.WaitGroup) {
	defer wg.Done()
	b.generateFiles(b.dirPath, b.filesToWrite, b.fileSizeMB)
	b.iterrateDir(b.dirPath)
}

func (b Benchmark) printReadResults() {
	dataInMB := float64(b.dataInBytes) / (1024 * 1024)
	log.Printf("Read: %6.3f MB in %6.3f seconds\n", dataInMB, b.readDuration)
	readSpeed := dataInMB / b.readDuration
	log.Printf("Read Speed is: %6.3f MB/s\n", readSpeed)
}

func (b Benchmark) printWriteResults() {
	wroteMB := float64(b.fileSizeMB * b.filesWritten)
	log.Printf("Wrote: %6.3f MB in %6.3f seconds\n", wroteMB, b.writeDuration)
	writeSpeed := wroteMB / b.writeDuration
	log.Printf("Write Speed is: %6.3f MB/s\n", writeSpeed)
}

func main() {

	args := getArgs()
	log.SetOutput(os.Stdout)

	var benchArr []*Benchmark
	var wg sync.WaitGroup
	for i := 1; i <= args.readers; i++ {
		bench := Benchmark{
			dirPath:      args.dirPath,
			filesToWrite: args.filesToWrite,
			fileSizeMB:   args.fileSize,
		}
		wg.Add(1)
		go bench.run(&wg)
		benchArr = append(benchArr, &bench)
	}

	wg.Wait()

	for _, bench := range benchArr {
		bench.printWriteResults()
		bench.printReadResults()
	}

}
