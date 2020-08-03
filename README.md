# gdd - Golang Diskdump [![Go Report Card](https://goreportcard.com/badge/github.com/MrNinso/gdd)](https://goreportcard.com/report/github.com/MrNinso/gdd)  

This is a dd copy in golang

````
NAME:
   gdd - A dd but in GOLANG !!

USAGE:
   gdd [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --input value, -i value                       Input File
   --output value, -o value                      Output File
   --hide-progress, --hp                         Remove Progress bar (default: false)
   --block-size value, --bs value                Size of copyFile2File block in bytes (default: "512")
   --block-count value, --count value, -c value  Size of copyFile2File block in bytes (default: "-1")
   --help, -h                                    show help (default: false)

````

## Build

````
git clone https://github.com/MrNinso/gdd.git
cd gdd
./buid.sh
cd build
````

