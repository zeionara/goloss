package main

import (
	"bufio"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/itchyny/volume-go"
	"github.com/lawl/pulseaudio"
	"github.com/pieterclaerhout/go-log"
)

func main() {
	minAdjustmentDelay := 1

	if len(os.Args) > 1 {
		minAdjustmentDelay, _ = strconv.Atoi(os.Args[1])
	}
	client, err := pulseaudio.NewClient()
	chk(err)

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
			levels := makeVolumeLevelsParsingLine(line)

			if levels != nil {
				lastLevelMeasurements = append(lastLevelMeasurements, levels.mean())
			}

			if levels != nil {
				secondsPassed := time.Since(start).Seconds()

				if secondsPassed >= float64(minAdjustmentDelay) {
					currentVolume, err := volume.GetVolume()
					chk(err)
					currentNormalizedVolume := float64(currentVolume) / 100.0

					mean_levels := mean(lastLevelMeasurements)
					perfectVolume := 1.2 - 3.1*mean_levels + 2.79*mean_levels*mean_levels
					log.Info("Perfect volume is ", perfectVolume, "; current volume is ", currentNormalizedVolume)

					if perfectVolume != currentNormalizedVolume {
						if perfectVolume < 0 {
							client.SetVolume(0)
						} else if perfectVolume > 1 {
							client.SetVolume(1)
						} else {
							client.SetVolume(float32(perfectVolume))
						}

						start = time.Now()
						lastLevelMeasurements = nil
					}
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
