package clusters

import (
	"github.com/deployment-io/deployment-runner-kit/enums/cluster_enums"
	"github.com/deployment-io/deployment-runner-kit/enums/region_enums"
	"github.com/deployment-io/deployment-runner-kit/types"
)

type UpsertClusterDtoV1 struct {
	Type                    cluster_enums.Type
	Region                  region_enums.Type
	Name                    string
	ClusterArn              string
	EcsTaskExecutionRoleArn string
}

type UpsertClustersArgsV1 struct {
	types.AuthArgsV1
	Clusters []UpsertClusterDtoV1
}

type UpsertClustersReplyV1 struct {
	Done bool
}

type UpsertClustersArgsV2 struct {
	types.AuthArgsV2
	Clusters []UpsertClusterDtoV1
}
