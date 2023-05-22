package deployment_enums

type CpuMemory uint

//add more based on
//https://docs.aws.amazon.com/AmazonECS/latest/developerguide/AWS_Fargate.html

const (
	TwoFiveSixCpuFiveOneTwoMB CpuMemory = iota + 1
	MaxCpuMemory                        //always add new cpu and memory combinations before MaxCpuMemory
)

var cpuMemoryToString = map[CpuMemory]string{
	TwoFiveSixCpuFiveOneTwoMB: ".25 vCPU, 0.5 GB",
}

var cpuMemoryToCpu = map[CpuMemory]string{
	TwoFiveSixCpuFiveOneTwoMB: ".25 vCPU",
}

var cpuMemoryToMemory = map[CpuMemory]string{
	TwoFiveSixCpuFiveOneTwoMB: "0.5 GB",
}

func (c CpuMemory) String() string {
	return cpuMemoryToString[c]
}

func (c CpuMemory) Cpu() string {
	return cpuMemoryToCpu[c]
}

func (c CpuMemory) Memory() string {
	return cpuMemoryToMemory[c]
}
