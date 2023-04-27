package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	pause      = flag.Int("pause", 2, "Time in seconds to pause between a words. Use number from 1 to 5")
	file       = flag.String("file", "words.txt", "File with words.")
	rhvProfile = flag.String("profile", "", "Specify a voice profile 'Elena' for Russian or 'slt' for English")
)

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func execute(s string) {
	_, err := exec.Command("bash", "-c", s).Output()
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	strPause := fmt.Sprintf("|pause%d.mp3|", *pause)
	strPauseLast := fmt.Sprintf("|pause%d.mp3", *pause)

	var toSplit string
	var splitted string

	lines, err := readLines(*file)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	for i, line := range lines {
		makeMp3 := fmt.Sprintf(`echo "%s" | RHVoice-test --profile "%s" -o outrhv%v.mp3`, line, *rhvProfile, i)
		execute(makeMp3)
		repair := fmt.Sprintf(`ffmpeg -y -i outrhv%d.mp3 -ar 24000 -ac 1 -b:a 64k outff%d.mp3`, i, i)
		execute(repair)

		fileName := fmt.Sprintf("outff%d.mp3", i)
		if i == len(lines)-1 {
			toSplit = fileName + strPauseLast
			splitted += toSplit
			break
		}
		toSplit = fileName + strPause
		splitted += toSplit
	}

	filenameWithoutExtension := strings.TrimSuffix(*file, ".txt")
	split := fmt.Sprintf(`ffmpeg -y -i "concat:%s" -ar 24000 -ac 1 -b:a 64k %s.mp3`, splitted, filenameWithoutExtension)
	execute(split)

	for j, _ := range lines {
		del := fmt.Sprintf(`rm outff%d.mp3 outrhv%d.mp3`, j, j)
		execute(del)
	}
}
