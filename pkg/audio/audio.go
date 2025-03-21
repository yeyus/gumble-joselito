package audio

import "log"

var ULawDecode = []int16{
	-32124, -31100, -30076, -29052, -28028, -27004, -25980, -24956,
	-23932, -22908, -21884, -20860, -19836, -18812, -17788, -16764,
	-15996, -15484, -14972, -14460, -13948, -13436, -12924, -12412,
	-11900, -11388, -10876, -10364, -9852, -9340, -8828, -8316,
	-7932, -7676, -7420, -7164, -6908, -6652, -6396, -6140,
	-5884, -5628, -5372, -5116, -4860, -4604, -4348, -4092,
	-3900, -3772, -3644, -3516, -3388, -3260, -3132, -3004,
	-2876, -2748, -2620, -2492, -2364, -2236, -2108, -1980,
	-1884, -1820, -1756, -1692, -1628, -1564, -1500, -1436,
	-1372, -1308, -1244, -1180, -1116, -1052, -988, -924,
	-876, -844, -812, -780, -748, -716, -684, -652,
	-620, -588, -556, -524, -492, -460, -428, -396,
	-372, -356, -340, -324, -308, -292, -276, -260,
	-244, -228, -212, -196, -180, -164, -148, -132,
	-120, -112, -104, -96, -88, -80, -72, -64,
	-56, -48, -40, -32, -24, -16, -8, 0,
	32124, 31100, 30076, 29052, 28028, 27004, 25980, 24956,
	23932, 22908, 21884, 20860, 19836, 18812, 17788, 16764,
	15996, 15484, 14972, 14460, 13948, 13436, 12924, 12412,
	11900, 11388, 10876, 10364, 9852, 9340, 8828, 8316,
	7932, 7676, 7420, 7164, 6908, 6652, 6396, 6140,
	5884, 5628, 5372, 5116, 4860, 4604, 4348, 4092,
	3900, 3772, 3644, 3516, 3388, 3260, 3132, 3004,
	2876, 2748, 2620, 2492, 2364, 2236, 2108, 1980,
	1884, 1820, 1756, 1692, 1628, 1564, 1500, 1436,
	1372, 1308, 1244, 1180, 1116, 1052, 988, 924,
	876, 844, 812, 780, 748, 716, 684, 652,
	620, 588, 556, 524, 492, 460, 428, 396,
	372, 356, 340, 324, 308, 292, 276, 260,
	244, 228, 212, 196, 180, 164, 148, 132,
	120, 112, 104, 96, 88, 80, 72, 64,
	56, 48, 40, 32, 24, 16, 8, 0}

var firCoeffs = []float64{
	-2.1734e-04, -5.2746e-04, -9.3654e-04, -1.4945e-03, -2.2148e-03, -3.0577e-03,
	-3.9204e-03, -4.6356e-03, -4.9800e-03, -4.6927e-03, -3.5020e-03, -1.1580e-03,
	2.5329e-03, 7.6751e-03, 1.4258e-02, 2.2142e-02, 3.1050e-02, 4.0587e-02,
	5.0259e-02, 5.9510e-02, 6.7773e-02, 7.4514e-02, 7.9283e-02, 8.1753e-02,
	8.1753e-02, 7.9283e-02, 7.4514e-02, 6.7773e-02, 5.9510e-02, 5.0259e-02,
	4.0587e-02, 3.1050e-02, 2.2142e-02, 1.4258e-02, 7.6751e-03, 2.5329e-03,
	-1.1580e-03, -3.5020e-03, -4.6927e-03, -4.9800e-03, -4.6356e-03, -3.9204e-03,
	-3.0577e-03, -2.2148e-03, -1.4945e-03, -9.3654e-04, -5.2746e-04, -2.1734e-04,
}

func UpsampleAndFilter(input []int16) []int16 {
	L := 6 // Upsampling factor
	output := make([]int16, len(input)*L)

	// Step 1: Zero-stuffing
	for i, sample := range input {
		output[i*L] = sample // Insert original sample
		// All other values are already zero (default in Go)
	}

	// Step 2: FIR Filtering
	N := len(firCoeffs)
	filtered := make([]int16, len(output))

	for i := range output {
		var acc float64
		for j := 0; j < N; j++ {
			if i-j >= 0 {
				acc += firCoeffs[j] * float64(output[i-j])
			}
		}

		// magic number amplification ¯\_(ツ)_/¯
		acc *= 10
		// Clip and convert back to int16
		if acc > 32767 {
			log.Printf("clipping")
			acc = 32767
		} else if acc < -32768 {
			log.Printf("clipping")
			acc = -32768
		}
		filtered[i] = int16(acc)
	}

	return filtered
}
