package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	programName = "blindTest"
	rate        = 48000
	nperiods    = 2
	period      = 256
	jackEnabled = false
	jackProc    *os.Process
	conv        = []string{
		"opusenc --bitrate 64.0 '{}' '{}'",
		"opusenc --bitrate 128.0 '{}' '{}'",
		"ffmpeg -i '{}' -b:a 128k '{}'",
		"ffmpeg -i '{}' -b:a 196k '{}'",
	}
)

func song2opus(song string) {
	fmt.Printf("Converting: %s\n", song)
	for i, j := range conv {
		out := fmt.Sprintf("%s-%d.opus", strings.TrimSuffix(song, filepath.Ext(song)), i)
		cmd := fmt.Sprintf(j, song, out)
		if *debug {
			fmt.Println(cmd)
		}
		cmdSplit := strings.Split(cmd, " ")
		subprocess := exec.Command(cmdSplit[0], cmdSplit[1:]...)
		subprocess.Run()
	}
}

func play(order []int) {
	if !*debug {
		clearScreen()
	}
	fmt.Printf("\n  > %s -- can *you* hear the difference?\n\n", programName)
	for i, j := range order {
		if j+1 > len(conv) {
			j = j - len(conv)
			if !jackEnabled {
				enableJack()
			}
		} else {
			if jackEnabled {
				disableJack()
			}
		}
		time.Sleep(1 * time.Second)
		fmt.Printf("Song no: %d (make a note!)\n", i+1)
		audioType := "pulse"
		if jackEnabled {
			audioType = "jack"
		}
		cmd := fmt.Sprintf("mpv -ao '%s' --start=1:00 --length=3 '%s-%d.opus'", audioType, strings.TrimSuffix(*song, filepath.Ext(*song)), j)
		if *debug {
			fmt.Println(cmd)
		}
		cmdSplit := strings.Split(cmd, " ")
		subprocess := exec.Command(cmdSplit[0], cmdSplit[1:]...)
		subprocess.Run()
		fmt.Print("Press <enter> to continue")
		var input string
		fmt.Scanln(&input)
	}
}

func showlist(order []int) {
	if jackEnabled {
		disableJack()
	}
	fmt.Println("\nThose were all the songs.")
	fmt.Println("\nPress <enter> to continue")
	for i, j := range order {
		moreJack := ""
		if j+1 > len(conv) {
			j = j - len(conv)
			moreJack = "- (played using jack)"
		}
		line := fmt.Sprintf("- Song no %d encoded with: %s %s", i+1, conv[j], moreJack)
		fmt.Println(strings.ReplaceAll(strings.ReplaceAll(line, "'{}' ", ""), "-i", ""))
	}
}

func enableJack() {
	cmd := fmt.Sprintf("pasuspender -- /usr/bin/jackd -R -P89 -dalsa -dhw:0 -r%d -p%d -n%d", rate, period, nperiods)
	if *debug {
		fmt.Println(cmd)
	}
	cmdSplit := strings.Split(cmd, " ")
	subprocess := exec.Command(cmdSplit[0], cmdSplit[1:]...)
	err := subprocess.Start()
	if err != nil {
		fmt.Println("Error starting jackd:", err)
	}
	jackProc = subprocess.Process
	jackEnabled = true
}

func disableJack() {
	if *debug {
		fmt.Println("Disable jack")
	}
	err := jackProc.Signal(os.Interrupt)
	if err != nil {
		fmt.Println("Error terminating jackd:", err)
	}
	jackEnabled = false
}

func run() {
	if *convert {
		song2opus(*song)
	}
	playlistStart := len(conv) * *jack
	if *jack == 2 {
		playlistStart = len(conv)
	}

	playlistLength := len(conv) * (2 - *jack)
	if *jack == 1 {
		playlistLength *= 2
	}
	order := make([]int, playlistLength)
	for i := range order {
		order[i] = (i + playlistStart) % len(conv)
	}

	play(order)
	showlist(order)
}

var (
	convert = flag.Bool("c", false, "Convert song from flac to opus")
	jack    = flag.Int("j", 0, "listen using with (and without) jack")
	debug   = flag.Bool("d", false, "enable debugging")
	song    = flag.String("song", "", "Reference song (should probably be .flac)")
)

func main() {
	flag.Parse()
	run()
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}