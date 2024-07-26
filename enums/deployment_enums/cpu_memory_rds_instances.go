package deployment_enums

type CpuMemoryRDSInstance uint

const (
	RdsT2Micro CpuMemoryRDSInstance = iota + 1
	RdsT3Micro
	RdsT2Small
	RdsT3Small
	RdsT2Medium
	RdsT3Medium
	RdsT2Large
	RdsT3Large
	RdsM5Large
	RdsT2XLarge
	RdsT3XLarge
	RdsM5XLarge
	RdsT22XLarge
	RdsT32XLarge
	RdsM52XLarge
	RdsM54XLarge
	RdsM58XLarge
	RdsM512XLarge
	RdsM516XLarge
	RdsM524XLarge

	MaxCpuMemoryRDSInstance
)

var CpuMemoryRDSInstanceToString = map[CpuMemoryRDSInstance]string{
	//RdsT2Micro:    "1 vCPU, 1 GB (db.t2.micro)",
	RdsT3Micro: "2 vCPU, 1 GB (db.t3.micro)",
	//RdsT2Small:    "1 vCPU, 2 GB (db.t2.small)",
	RdsT3Small: "2 vCPU, 2 GB (db.t3.small)",
	//RdsT2Medium:   "2 vCPU, 4 GB (db.t2.medium)",
	RdsT3Medium: "2 vCPU, 4 GB (db.t3.medium)",
	//RdsT2Large:    "2 vCPU, 8 GB (db.t2.large)",
	RdsT3Large:    "2 vCPU, 8 GB (db.t3.large)",
	RdsM5Large:    "2 vCPU, 8 GB (db.m5.large)",
	RdsT2XLarge:   "4 vCPU, 16 GB (db.t2.xlarge)",
	RdsT3XLarge:   "4 vCPU, 16 GB (db.t3.xlarge)",
	RdsM5XLarge:   "4 vCPU, 16 GB (db.m5.xlarge)",
	RdsT22XLarge:  "8 vCPU, 32 GB (db.t2.2xlarge)",
	RdsT32XLarge:  "8 vCPU, 32 GB (db.t3.2xlarge)",
	RdsM52XLarge:  "8 vCPU, 32 GB (db.m5.2xlarge)",
	RdsM54XLarge:  "16 vCPU, 64 GB (db.m5.4xlarge)",
	RdsM58XLarge:  "32 vCPU, 128 GB (db.m5.8xlarge)",
	RdsM512XLarge: "48 vCPU, 192 GB (db.m5.12xlarge)",
	RdsM516XLarge: "64 vCPU, 256 GB (db.m5.16xlarge)",
	RdsM524XLarge: "96 vCPU, 384 GB (db.m5.24xlarge)",
}

var cpuMemoryRDSInstanceToInstance = map[CpuMemoryRDSInstance]string{
	RdsT2Micro:    "db.t2.micro",
	RdsT3Micro:    "db.t3.micro",
	RdsT2Small:    "db.t2.small",
	RdsT3Small:    "db.t3.small",
	RdsT2Medium:   "db.t2.medium",
	RdsT3Medium:   "db.t3.medium",
	RdsT2Large:    "db.t2.large",
	RdsT3Large:    "db.t3.large",
	RdsM5Large:    "db.m5.large",
	RdsT2XLarge:   "db.t2.xlarge",
	RdsT3XLarge:   "db.t3.xlarge",
	RdsM5XLarge:   "db.m5.xlarge",
	RdsT22XLarge:  "db.t2.2xlarge",
	RdsT32XLarge:  "db.t3.2xlarge",
	RdsM52XLarge:  "db.m5.2xlarge",
	RdsM54XLarge:  "db.m5.4xlarge",
	RdsM58XLarge:  "db.m5.8xlarge",
	RdsM512XLarge: "db.m5.12xlarge",
	RdsM516XLarge: "db.m5.16xlarge",
	RdsM524XLarge: "db.m5.24xlarge",
}

func (c CpuMemoryRDSInstance) Instance() string {
	return cpuMemoryRDSInstanceToInstance[c]
}

func (c CpuMemoryRDSInstance) SupportsEncryption() bool {
	if c == RdsT2Micro {
		return false
	}
	return true
}

var InstanceToCpuMemoryRDSInstance = map[string]CpuMemoryRDSInstance{
	"db.t2.micro":    RdsT2Micro,
	"db.t3.micro":    RdsT3Micro,
	"db.t2.small":    RdsT2Small,
	"db.t3.small":    RdsT3Small,
	"db.t2.medium":   RdsT2Medium,
	"db.t3.medium":   RdsT3Medium,
	"db.t2.large":    RdsT2Large,
	"db.t3.large":    RdsT3Large,
	"db.m5.large":    RdsM5Large,
	"db.t2.xlarge":   RdsT2XLarge,
	"db.t3.xlarge":   RdsT3XLarge,
	"db.m5.xlarge":   RdsM5XLarge,
	"db.t2.2xlarge":  RdsT22XLarge,
	"db.t3.2xlarge":  RdsT32XLarge,
	"db.m5.2xlarge":  RdsM52XLarge,
	"db.m5.4xlarge":  RdsM54XLarge,
	"db.m5.8xlarge":  RdsM58XLarge,
	"db.m5.12xlarge": RdsM512XLarge,
	"db.m5.16xlarge": RdsM516XLarge,
	"db.m5.24xlarge": RdsM524XLarge,
}
