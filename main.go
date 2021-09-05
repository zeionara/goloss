package main

import (
	"bufio"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/itchyny/volume-go"
	"github.com/lawl/pulseaudio"
	"github.com/pieterclaerhout/go-log"
)

func mean(measurements []float64) float64 {
	if len(measurements) == 0 {
		return 0
	}

	sum := float64(0)

	for i := 0; i < len(measurements); i++ {
		sum += (measurements[i])
	}

	return float64(sum) / float64(len(measurements))
}

func main() {
	reg := regexp.MustCompile(`(?P<LeftLevel>[0-9.]+)\s+(?P<RightLevel>[0-9.]+)`)
	minAdjustmentDelay := 1
	// windowSize := 10
	// targetRatio := 0.5
	if len(os.Args) > 1 {
		minAdjustmentDelay, _ = strconv.Atoi(os.Args[1])
	}
	client, err := pulseaudio.NewClient()
	chk(err)
	// print(client)

	log.PrintTimestamp = true

	cmd := exec.Command("pavumeterc")

	r, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	done := make(chan struct{})
	scanner := bufio.NewScanner(r)

	var lastLevelMeasurements []float64

	go func() {
		start := time.Now()

		for scanner.Scan() {
			line := scanner.Text()
			levels := makeVolumeLevelsParsingLine(line, reg)
			// mean_levels := levels.mean()
			lastLevelMeasurements = append(lastLevelMeasurements, levels.mean())

			if levels != nil {
				// fmt.Printf("%f\n", levels.mean())
				// log.Info(levels.mean())
				secondsPassed := time.Since(start).Seconds()

				if secondsPassed >= float64(minAdjustmentDelay) {
					currentVolume, err := volume.GetVolume()
					chk(err)
					currentNormalizedVolume := float64(currentVolume) / 100.0

					// ratio := levels.mean() / float64(currentNormalizedVolume)
					// print("Ratio = ", ratio, " levels = ",  " volume = ", )
					// client.ServerInfo()
					// log.Info(" levels = ", levels.mean(), " volume = ", currentNormalizedVolume)
					// start = time.Now()
					mean_levels := mean(lastLevelMeasurements)
					perfectVolume := 1.2 - 3.1*mean_levels + 2.79*mean_levels*mean_levels // levels.mean() / targetRatio
					log.Info("Perfect volume is ", perfectVolume, "; current volume is ", currentNormalizedVolume)

					if perfectVolume != currentNormalizedVolume {
						if perfectVolume < 0 {
							client.SetVolume(0)
						} else if perfectVolume > 1 {
							client.SetVolume(1)
						} else {
							client.SetVolume(float32(perfectVolume))
						}

						// if perfectVolume < currentNormalizedVolume {
						start = time.Now()
						lastLevelMeasurements = nil
						// }
					}

					// if levels.mean() < 0.7 && currentNormalizedVolume != 0.8 {
					// 	if false {
					// 		client.SetVolume(0.8)
					// 	}
					// 	start = time.Now() // it's preferred to decrase volume even if minimum amount of time has not yet passed
					// } else if currentNormalizedVolume != 0.3 {
					// 	// client.SetVolume(0.3)
					// 	start = time.Now()
					// }
				}
			}
		}

		done <- struct{}{}
	}()

	err = cmd.Start()
	log.CheckError(err)
	<-done
	err = cmd.Wait()
	log.CheckError(err)
}
