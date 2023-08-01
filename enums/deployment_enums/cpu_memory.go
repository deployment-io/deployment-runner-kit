package deployment_enums

type CpuMemory uint

//add more based on
//https://docs.aws.amazon.com/AmazonECS/latest/developerguide/AWS_Fargate.html

const (
	PointTwoFiveVCpuPointFiveGB CpuMemory = iota + 1
	PointTwoFiveVCpuOneGB
	PointTwoFiveVCpuTwoGB

	PointFiveVCpuOneGB
	PointFiveVCpuTwoGB
	PointFiveVCpuThreeGB
	PointFiveVCpuFourGB

	OneVCpuTwoGB
	OneVCpuThreeGB
	OneVCpuFourGB
	OneVCpuFiveGB
	OneVCpuSixGB
	OneVCpuSevenGB
	OneVCpuEightGB

	TwoVCpuFourGB
	TwoVCpuFiveGB
	TwoVCpuSixGB
	TwoVCpuSevenGB
	TwoVCpuEightGB
	TwoVCpuNineGB
	TwoVCpuTenGB
	TwoVCpuElevenGB
	TwoVCpuTwelveGB
	TwoVCpuThirteenGB
	TwoVCpuFourteenGB
	TwoVCpuFifteenGB
	TwoVCpuSixteenGB

	FourVCpuEightGB
	FourVCpuNineGB
	FourVCpuTenGB
	FourVCpuElevenGB
	FourVCpuTwelveGB
	FourVCpuThirteenGB
	FourVCpuFourteenGB
	FourVCpuFifteenGB
	FourVCpuSixteenGB
	FourVCpuSeventeenGB
	FourVCpuEighteenGB
	FourVCpuNineteenGB
	FourVCpuTwentyGB
	FourVCpuTwentyOneGB
	FourVCpuTwentyTwoGB
	FourVCpuTwentyThreeGB
	FourVCpuTwentyFourGB
	FourVCpuTwentyFifeGB
	FourVCpuTwentySixGB
	FourVCpuTwentySevenGB
	FourVCpuTwentyEightGB
	FourVCpuTwentyNineGB
	FourVCpuThirtyGB

	MaxCpuMemory //always add new cpu and memory combinations before MaxCpuMemory
)

var CpuMemoryToString = map[CpuMemory]string{
	PointTwoFiveVCpuPointFiveGB: ".25 vCPU, 0.5 GB",
	PointTwoFiveVCpuOneGB:       ".25 vCPU, 1 GB",
	PointTwoFiveVCpuTwoGB:       ".25 vCPU, 2 GB",

	PointFiveVCpuOneGB:   ".5 vCPU, 1 GB",
	PointFiveVCpuTwoGB:   ".5 vCPU, 2 GB",
	PointFiveVCpuThreeGB: ".5 vCPU, 3 GB",
	PointFiveVCpuFourGB:  ".5 vCPU, 4 GB",

	OneVCpuTwoGB:   "1 vCPU, 2 GB",
	OneVCpuThreeGB: "1 vCPU, 3 GB",
	OneVCpuFourGB:  "1 vCPU, 4 GB",
	OneVCpuFiveGB:  "1 vCPU, 5 GB",
	OneVCpuSixGB:   "1 vCPU, 6 GB",
	OneVCpuSevenGB: "1 vCPU, 7 GB",
	OneVCpuEightGB: "1 vCPU, 8 GB",

	TwoVCpuFourGB:     "2 vCPU, 4 GB",
	TwoVCpuFiveGB:     "2 vCPU, 5 GB",
	TwoVCpuSixGB:      "2 vCPU, 6 GB",
	TwoVCpuSevenGB:    "2 vCPU, 7 GB",
	TwoVCpuEightGB:    "2 vCPU, 8 GB",
	TwoVCpuNineGB:     "2 vCPU, 9 GB",
	TwoVCpuTenGB:      "2 vCPU, 10 GB",
	TwoVCpuElevenGB:   "2 vCPU, 11 GB",
	TwoVCpuTwelveGB:   "2 vCPU, 12 GB",
	TwoVCpuThirteenGB: "2 vCPU, 13 GB",
	TwoVCpuFourteenGB: "2 vCPU, 14 GB",
	TwoVCpuFifteenGB:  "2 vCPU, 15 GB",
	TwoVCpuSixteenGB:  "2 vCPU, 16 GB",

	FourVCpuEightGB:       "4 vCPU, 8 GB",
	FourVCpuNineGB:        "4 vCPU, 9 GB",
	FourVCpuTenGB:         "4 vCPU, 10 GB",
	FourVCpuElevenGB:      "4 vCPU, 11 GB",
	FourVCpuTwelveGB:      "4 vCPU, 12 GB",
	FourVCpuThirteenGB:    "4 vCPU, 13 GB",
	FourVCpuFourteenGB:    "4 vCPU, 14 GB",
	FourVCpuFifteenGB:     "4 vCPU, 15 GB",
	FourVCpuSixteenGB:     "4 vCPU, 16 GB",
	FourVCpuSeventeenGB:   "4 vCPU, 17 GB",
	FourVCpuEighteenGB:    "4 vCPU, 18 GB",
	FourVCpuNineteenGB:    "4 vCPU, 19 GB",
	FourVCpuTwentyGB:      "4 vCPU, 20 GB",
	FourVCpuTwentyOneGB:   "4 vCPU, 21 GB",
	FourVCpuTwentyTwoGB:   "4 vCPU, 22 GB",
	FourVCpuTwentyThreeGB: "4 vCPU, 23 GB",
	FourVCpuTwentyFourGB:  "4 vCPU, 24 GB",
	FourVCpuTwentyFifeGB:  "4 vCPU, 25 GB",
	FourVCpuTwentySixGB:   "4 vCPU, 26 GB",
	FourVCpuTwentySevenGB: "4 vCPU, 27 GB",
	FourVCpuTwentyEightGB: "4 vCPU, 28 GB",
	FourVCpuTwentyNineGB:  "4 vCPU, 29 GB",
	FourVCpuThirtyGB:      "4 vCPU, 30 GB",
}

var cpuMemoryToCpu = map[CpuMemory]string{
	PointTwoFiveVCpuPointFiveGB: ".25 vCPU",
	PointTwoFiveVCpuOneGB:       ".25 vCPU",
	PointTwoFiveVCpuTwoGB:       ".25 vCPU",

	PointFiveVCpuOneGB:   ".5 vCPU",
	PointFiveVCpuTwoGB:   ".5 vCPU",
	PointFiveVCpuThreeGB: ".5 vCPU",
	PointFiveVCpuFourGB:  ".5 vCPU",

	OneVCpuTwoGB:   "1 vCPU",
	OneVCpuThreeGB: "1 vCPU",
	OneVCpuFourGB:  "1 vCPU",
	OneVCpuFiveGB:  "1 vCPU",
	OneVCpuSixGB:   "1 vCPU",
	OneVCpuSevenGB: "1 vCPU",
	OneVCpuEightGB: "1 vCPU",

	TwoVCpuFourGB:     "2 vCPU",
	TwoVCpuFiveGB:     "2 vCPU",
	TwoVCpuSixGB:      "2 vCPU",
	TwoVCpuSevenGB:    "2 vCPU",
	TwoVCpuEightGB:    "2 vCPU",
	TwoVCpuNineGB:     "2 vCPU",
	TwoVCpuTenGB:      "2 vCPU",
	TwoVCpuElevenGB:   "2 vCPU",
	TwoVCpuTwelveGB:   "2 vCPU",
	TwoVCpuThirteenGB: "2 vCPU",
	TwoVCpuFourteenGB: "2 vCPU",
	TwoVCpuFifteenGB:  "2 vCPU",
	TwoVCpuSixteenGB:  "2 vCPU",

	FourVCpuEightGB:       "4 vCPU",
	FourVCpuNineGB:        "4 vCPU",
	FourVCpuTenGB:         "4 vCPU",
	FourVCpuElevenGB:      "4 vCPU",
	FourVCpuTwelveGB:      "4 vCPU",
	FourVCpuThirteenGB:    "4 vCPU",
	FourVCpuFourteenGB:    "4 vCPU",
	FourVCpuFifteenGB:     "4 vCPU",
	FourVCpuSixteenGB:     "4 vCPU",
	FourVCpuSeventeenGB:   "4 vCPU",
	FourVCpuEighteenGB:    "4 vCPU",
	FourVCpuNineteenGB:    "4 vCPU",
	FourVCpuTwentyGB:      "4 vCPU",
	FourVCpuTwentyOneGB:   "4 vCPU",
	FourVCpuTwentyTwoGB:   "4 vCPU",
	FourVCpuTwentyThreeGB: "4 vCPU",
	FourVCpuTwentyFourGB:  "4 vCPU",
	FourVCpuTwentyFifeGB:  "4 vCPU",
	FourVCpuTwentySixGB:   "4 vCPU",
	FourVCpuTwentySevenGB: "4 vCPU",
	FourVCpuTwentyEightGB: "4 vCPU",
	FourVCpuTwentyNineGB:  "4 vCPU",
	FourVCpuThirtyGB:      "4 vCPU",
}

var cpuMemoryToMemory = map[CpuMemory]string{
	PointTwoFiveVCpuPointFiveGB: "0.5 GB",
	PointTwoFiveVCpuOneGB:       "1 GB",
	PointTwoFiveVCpuTwoGB:       "2 GB",

	PointFiveVCpuOneGB:   "1 GB",
	PointFiveVCpuTwoGB:   "2 GB",
	PointFiveVCpuThreeGB: "3 GB",
	PointFiveVCpuFourGB:  "4 GB",

	OneVCpuTwoGB:   "2 GB",
	OneVCpuThreeGB: "3 GB",
	OneVCpuFourGB:  "4 GB",
	OneVCpuFiveGB:  "5 GB",
	OneVCpuSixGB:   "6 GB",
	OneVCpuSevenGB: "7 GB",
	OneVCpuEightGB: "8 GB",

	TwoVCpuFourGB:     "4 GB",
	TwoVCpuFiveGB:     "5 GB",
	TwoVCpuSixGB:      "6 GB",
	TwoVCpuSevenGB:    "7 GB",
	TwoVCpuEightGB:    "8 GB",
	TwoVCpuNineGB:     "9 GB",
	TwoVCpuTenGB:      "10 GB",
	TwoVCpuElevenGB:   "11 GB",
	TwoVCpuTwelveGB:   "12 GB",
	TwoVCpuThirteenGB: "13 GB",
	TwoVCpuFourteenGB: "14 GB",
	TwoVCpuFifteenGB:  "15 GB",
	TwoVCpuSixteenGB:  "16 GB",

	FourVCpuEightGB:       "8 GB",
	FourVCpuNineGB:        "9 GB",
	FourVCpuTenGB:         "10 GB",
	FourVCpuElevenGB:      "11 GB",
	FourVCpuTwelveGB:      "12 GB",
	FourVCpuThirteenGB:    "13 GB",
	FourVCpuFourteenGB:    "14 GB",
	FourVCpuFifteenGB:     "15 GB",
	FourVCpuSixteenGB:     "16 GB",
	FourVCpuSeventeenGB:   "17 GB",
	FourVCpuEighteenGB:    "18 GB",
	FourVCpuNineteenGB:    "19 GB",
	FourVCpuTwentyGB:      "20 GB",
	FourVCpuTwentyOneGB:   "21 GB",
	FourVCpuTwentyTwoGB:   "22 GB",
	FourVCpuTwentyThreeGB: "23 GB",
	FourVCpuTwentyFourGB:  "24 GB",
	FourVCpuTwentyFifeGB:  "25 GB",
	FourVCpuTwentySixGB:   "26 GB",
	FourVCpuTwentySevenGB: "27 GB",
	FourVCpuTwentyEightGB: "28 GB",
	FourVCpuTwentyNineGB:  "29 GB",
	FourVCpuThirtyGB:      "30 GB",
}

func (c CpuMemory) String() string {
	return CpuMemoryToString[c]
}

func (c CpuMemory) Cpu() string {
	return cpuMemoryToCpu[c]
}

func (c CpuMemory) Memory() string {
	return cpuMemoryToMemory[c]
}
