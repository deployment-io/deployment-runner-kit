package runner_enums

import "fmt"

type TargetCloud uint

const (
	AwsCloud       TargetCloud = iota + 1
	MaxTargetCloud             //always add new target clouds before MaxTargetCloud
)

var targetCloudsToString = map[TargetCloud]string{
	AwsCloud: "aws",
}

func (t TargetCloud) String() string {
	return targetCloudsToString[t]
}

func GetTargetCloudFromString(targetCloudStr string) (TargetCloud, error) {
	switch targetCloudStr {
	case "aws":
		return AwsCloud, nil
	}
	return 0, fmt.Errorf("invalid target cloud: %s", targetCloudStr)
}
