package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/itchyny/volume-go"
	"github.com/lawl/pulseaudio"
	"github.com/pieterclaerhout/go-log"
)

type employee struct {
	firstName   string
	lastName    string
	totalLeaves int
	leavesTaken int
}

func News(firstName string, lastName string, totalLeave int, leavesTaken int) employee {
	e := employee{firstName, lastName, totalLeave, leavesTaken}
	return e
}

func (e employee) LeavesRemaining() {
	fmt.Printf("%s %s has %d leaves remaining\n", e.firstName, e.lastName, (e.totalLeaves - e.leavesTaken))
}

type VolumeLevels struct {
	left  float64
	right float64
}

func makeVolumeLevels(left string, right string) VolumeLevels {
	left_, err := strconv.ParseFloat(left, 32)
	chk(err)

	right_, err := strconv.ParseFloat(right, 32)
	chk(err)

	levels := VolumeLevels{
		left:  left_,
		right: right_,
	}

	return levels
}

func (levels VolumeLevels) mean() float64 {
	return (levels.left + levels.right) / 2.0
}

func makeVolumeLevelsParsingLine(line string, r *regexp.Regexp) *VolumeLevels { // VolumeLevels { (levels VolumeLevels)
	// matched, err := regexp.MatchString(line, "[0-9.]+\\s+[0-9.]")
	// chk(err)
	// print(matched.)
	submatch := r.FindStringSubmatch(line)
	// fmt.Printf("%#v\n", submatch)
	// fmt.Printf("%#v\n", r.SubexpNames())

	if len(submatch) > 2 {
		value := makeVolumeLevels(submatch[1], submatch[2])
		return &value
	}
	return nil
	// return levels.New(submatch[0], submatch[1])
}

func main() {
	// bar := employee.News("2.2", "2.3")
	reg := regexp.MustCompile(`(?P<LeftLevel>[0-9.]+)\s+(?P<RightLevel>[0-9.]+)`)
	minAdjustmentDelay := 5
	if len(os.Args) > 1 {
		minAdjustmentDelay, _ = strconv.Atoi(os.Args[1])
	}
	client, err := pulseaudio.NewClient()
	chk(err)
	// out, err := exec.Command("./run.sh").Output()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("The date is %s\n", out)

	// Print the log timestamps
	log.PrintTimestamp = true

	// The command you want to run along with the argument
	cmd := exec.Command("pavumeterc")

	// Get a pipe to read from standard out
	r, _ := cmd.StdoutPipe()

	// Use the same pipe for standard error
	cmd.Stderr = cmd.Stdout

	// Make a new channel which will be used to ensure we get all output
	done := make(chan struct{})

	// Create a scanner which scans r in a line-by-line fashion
	scanner := bufio.NewScanner(r)

	// Use the scanner to scan the output line by line and log it
	// It's running in a goroutine so that it doesn't block
	go func() {

		// Read line by line and process it
		start := time.Now()

		for scanner.Scan() {
			line := scanner.Text()
			levels := makeVolumeLevelsParsingLine(line, reg)
			if levels != nil {
				// fmt.Printf("%f\n", levels.mean())
				log.Info(levels.mean())
				secondsPassed := time.Since(start).Seconds()

				if secondsPassed >= float64(minAdjustmentDelay) {
					currentVolume, err := volume.GetVolume()
					chk(err)
					currentNormalizedVolume := float64(currentVolume) / 100.0

					if levels.mean() < 0.5 && currentNormalizedVolume != 0.8 {
						client.SetVolume(0.8)
						start = time.Now()
					} else if currentNormalizedVolume != 0.3 {
						client.SetVolume(0.3)
						start = time.Now()
					}
				}
			}
		}

		// We're all done, unblock the channel
		done <- struct{}{}

	}()

	// Start the command and check for errors
	err = cmd.Start()
	log.CheckError(err)

	// Wait for all output to be processed
	<-done

	// Wait for the command to finish
	err = cmd.Wait()
	log.CheckError(err)
}

// func main() {
// 	portaudio.Initialize()
// 	defer portaudio.Terminate()
// 	h, err := portaudio.DefaultHostApi()
// 	chk(err)
// 	var n int32
// 	n = 0
// 	stream, err := portaudio.OpenStream(portaudio.HighLatencyParameters(nil, h.DefaultOutputDevice), func(out []int32) {
// 		for i := range out {
// 			out[i] = int32(rand.Uint32())
// 		}
// 		println(n)
// 		// n += 10000
// 	})
// 	chk(err)
// 	defer stream.Close()
// 	chk(stream.Start())
// 	time.Sleep(30 * time.Second)
// 	chk(stream.Stop())
// }

// func main() {
// 	portaudio.Initialize()
// 	defer portaudio.Terminate()
// 	h, err := portaudio.DefaultHostApi()
// 	chk(err)
// 	var s []float64
// 	stream, err := portaudio.OpenStream(portaudio.HighLatencyParameters(h.DefaultOutputDevice, nil), func(in []int32) {
// 		// print(in[0])
// 		for j := range in {
// 			// println(in[i])
// 			// fmt.Printf(float64(int32(rand.Uint32())) / 2147483647)
// 			// fmt.Printf("%f\n", float64(int32(rand.Uint32()))/2147483647)
// 			// in[i] = 20 // int32(rand.Uint32())
// 			s = append(s, float64(in[j])/2147483647)
// 			if len(s) >= 10000 {
// 				// print(s)

// 				sum := 0.0

// 				// traversing through the
// 				// array using for loop
// 				for i := 0; i < len(s); i++ {

// 					// adding the values of
// 					// array to the variable sum
// 					sum += (s[i])
// 				}

// 				// declaring a variable
// 				// avg to find the average
// 				avg := (float64(sum)) / (float64(len(s)))

// 				fmt.Printf("%f\n", avg)
// 				s = make([]float64, 0)
// 			}
// 			// print(s)
// 			// fmt.Printf("%f\n", float64(in[i])/2147483647)
// 		}
// 	})
// 	chk(err)
// 	defer stream.Close()
// 	chk(stream.Start())
// 	time.Sleep(30 * time.Second)
// 	chk(stream.Stop())
// }

// func main() {
// 	fmt.Println("Hello, world.")

// 	vol, err := volume.GetVolume()
// 	if err != nil {
// 		log.Fatalf("get volume failed: %+v", err)
// 	}
// 	fmt.Printf("current volume: %d\n", vol)

// 	err = volume.SetVolume(10)
// 	if err != nil {
// 		log.Fatalf("set volume failed: %+v", err)
// 	}
// 	fmt.Printf("set volume success\n")

// 	// client, err := pulseaudio.NewClient()
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// client.Mute()
// 	// sinks, err := client.Sinks()
// 	// mainSink := sinks[0]

// 	// client.SetSinkVolume(mainSink.Name, 0.2)
// 	// client.Mute()
// 	// client.SetVolume(0.8)

// 	// r := strings.NewReader("Hello, Reader! XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")

// 	// b := make([]byte, 8)
// 	// for {
// 	// 	n, err := r.Read(b)
// 	// 	fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
// 	// 	fmt.Printf("b[:n] = %q\n", b[:n])
// 	// 	if err == io.EOF {
// 	// 		break
// 	// 	}
// 	// }

// 	// info, err := client.ServerInfo()
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// print(info.DefaultSink)

// 	// mainSink.ReadFrom(r)

// 	// defer client.Close()

// 	//
// 	//
// 	//

// 	// portaudio.Initialize()
// 	// defer portaudio.Terminate()
// 	// in := make([]int32, 64)
// 	// stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
// 	// chk(err)
// 	// defer stream.Close()

// 	// chk(stream.Start())
// 	// for {
// 	// 	chk(stream.Read())
// 	// 	print(in[:2])
// 	// 	// chk(binary.Write(f, binary.BigEndian, in))
// 	// 	// nSamples += len(in)
// 	// 	// select {
// 	// 	// case <-sig:
// 	// 	// 	return
// 	// 	// default:
// 	// 	// }
// 	// }
// 	// chk(stream.Stop())

// 	newEcho(22)
// }

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

// func newEcho(delay time.Duration) *echo {
// 	h, err := portaudio.DefaultHostApi()
// 	chk(err)
// 	p := portaudio.LowLatencyParameters(h.DefaultInputDevice, h.DefaultOutputDevice)
// 	p.Input.Channels = 1
// 	p.Output.Channels = 1
// 	e := &echo{buffer: make([]float32, int(p.SampleRate*delay.Seconds()))}
// 	e.Stream, err = portaudio.OpenStream(p, e.processAudio)
// 	chk(err)
// 	return e
// }
