package main

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

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
