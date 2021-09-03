package main

import (
	"math/rand"
	"time"

	"github.com/gordonklaus/portaudio"
)

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()
	h, err := portaudio.DefaultHostApi()
	chk(err)
	var n int32
	n = 0
	stream, err := portaudio.OpenStream(portaudio.HighLatencyParameters(nil, h.DefaultOutputDevice), func(out []int32) {
		for i := range out {
			out[i] = int32(rand.Uint32())
		}
		println(n)
		// n += 10000
	})
	chk(err)
	defer stream.Close()
	chk(stream.Start())
	time.Sleep(30 * time.Second)
	chk(stream.Stop())
}

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
