# diskbench

##### Multi threaded disk benchmarking tool written in go

### Usage

```bash
mkdir /abcd/tmp #abcd is the tested drive mount
./diskbench -dir /abcd/tmp -files 1000 -size 8

./diskbench --help                          
Usage of ./diskbench:
  -dir string
        directory path to read (default ".")
  -files int
        amount of files to write
  -readers int
        amount of reader threads (default 1)
  -size int
        file size in MB (to write) (default 10)
```

