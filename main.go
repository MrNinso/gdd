package main

import (
	"github.com/cheggaaa/pb/v3"
	. "github.com/logrusorgru/aurora"
	"github.com/urfave/cli/v2"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type MyFile struct {
	*os.File
	size int64
}

var HideBar bool

var Notation = []string{"", "k", "m", "g", "t", "p", "e", "z" }

func main() {
	app := cli.NewApp()

	app.Name = "gdd"
	app.EnableBashCompletion = true
	app.Usage = "A dd but in GOLANG !!"

	app.Flags = []cli.Flag {
		&cli.PathFlag{
			Name:      "input",
			Aliases:   []string{ "i" },
			Usage:     "Input File",
			Required:  true,
			TakesFile: true,
		},
		&cli.PathFlag{
			Name:      "output",
			Aliases:   []string{ "o" },
			Usage:     "Output File",
			Required:  true,
			TakesFile: true,
		},
		&cli.BoolFlag{
			Name:        "hide-progress",
			Aliases:     []string{ "hp" },
			Usage:       "Remove Progress bar",
			Destination: &HideBar,
			HasBeenSet:  false,
		},
		&cli.StringFlag{
			Name:    "block-size",
			Aliases: []string { "bs" },
			Usage:   "Size of copyFile2File block in bytes",
			Value:   "512",
		},
		&cli.StringFlag{
			Name:    "block-count",
			Aliases: []string { "count", "c" },
			Usage:   "Size of copyFile2File block in bytes",
			Value:   "-1",
		},
	}

	app.Action = func(context *cli.Context) error {
		input := getFile(context.Path("input"), false)
		output := getFile(context.Path("output"), true)

		var bar *pb.ProgressBar

		defer input.Close()
		defer output.Close()

		block := make([]byte, getInt64(context.String("block-size")))
		bc := context.String("block-count")
		var blockCount int64

		if bc == "-1" {
			blockCount = -1
		} else {
			blockCount = getInt64(bc)
		}

		var pos int64 = 0
		var count int64 = 0

		if !HideBar {
			if blockCount != -1 {
				bar = pb.Full.Start64(blockCount)
			} else {
				bar = pb.Full.Start64(input.size)
			}
			bar.StartTime()
		}

		if blockCount != -1 {
			for copyFile2File(input, output, &block, &pos) {
				count++

				if bar != nil {
					bar.SetCurrent(count)
				}

				if count >= blockCount {
					break
				}

			}
		} else {
			for copyFile2File(input, output, &block, &pos) {
				if bar != nil {
					bar.SetCurrent(pos)
				}
			}
		}

		if bar != nil {
			bar.Finish()
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		logErrorFatal(err)
	}
}

func copyFile2File(input, output *MyFile, block *[]byte, pos *int64) (bool) {
	c,_ := input.ReadAt(*block, *pos)

	if c == 0 {
		return false
	}

	_, err := output.WriteAt((*block)[0:c], *pos)

	if err != nil {
		logErrorFatal(err)
		return false
	}

	*pos += int64(c)

	return true
}

func getInt64(s string) int64 {
	m := 0
	n := strings.ToLower(s)
	for i := 1; i < len(Notation); i++ {
		if strings.Contains(n, Notation[i]) {
			m = i*3
			n = strings.Replace(n, Notation[i], "", -1)
			break
		}
	}

	coefficient, err := strconv.ParseInt(n, 10, 64)

	if err != nil {
		logErrorFatal("cant convert ", s)
	}

	result := coefficient * int64(math.Pow10(m))

	if result < 0 {
		logErrorFatal(s, "is negative")
	}

	return result
}

func getFile(path string, isOutput bool) (file *MyFile) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) && !isOutput  {
		logErrorFatal(err)
		os.Exit(404)
	}

	var f *os.File
	var size int64

	if isOutput {
		f, err = os.OpenFile(path, os.O_WRONLY | os.O_CREATE, 0666)

		if f == nil {
			logErrorFatal("cant open or create the file", path)
			os.Exit(1)
		}

		if info != nil && info.Size() != 0 {
			errTruncate := f.Truncate(0)
			if errTruncate != nil {
				logErrorFatal(errTruncate)
				os.Exit(1)
			}
		}

		size = 0
	} else {
		f, err = os.OpenFile(path, os.O_RDWR, os.ModePerm)
		size = info.Size()
	}

	if err != nil {
		logErrorFatal(err)
		os.Exit(1)
	}
	
	return &MyFile{
		f,
		size,
	}
}

func logErrorFatal(v ...interface{}) {
	log.Fatal(Red(Bold(v)))
}