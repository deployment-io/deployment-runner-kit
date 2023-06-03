package vpcs

import (
	"github.com/deployment-io/deployment-runner-kit/enums/region_enums"
	"github.com/deployment-io/deployment-runner-kit/enums/vpc_enums"
	"github.com/deployment-io/deployment-runner-kit/types"
)

const (
	True  = "true"
	False = "false"
)

type SubnetDtoV1 struct {
	Name      string
	ID        string
	Cidr      string
	Az        string
	IsPrivate string
}

type RouteTableDtoV1 struct {
	Name      string
	ID        string
	IsPrivate string
}

type UpsertVpcDtoV1 struct {
	Type                  vpc_enums.Type
	Region                region_enums.Type
	VpcId                 string
	VpcCidr               string
	InternetGatewayId     string
	ElasticIPAllocationId string
	NatGatewayId          string
	Subnets               []SubnetDtoV1
	RouteTables           []RouteTableDtoV1
}

type UpsertVpcsArgsV1 struct {
	types.AuthArgsV1
	Vpcs []UpsertVpcDtoV1
}

type UpsertVpcsReplyV1 struct {
	Done bool
}

func GetSubnetsMapFromDto(subnetsFromDB []SubnetDtoV1) map[string]string {
	subnetsMap := make(map[string]string)
	for _, subnet := range subnetsFromDB {
		subnetsMap[subnet.Name] = subnet.ID
	}
	return subnetsMap
}

func GetRouteTablesMapFromDto(routeTablesFromDB []RouteTableDtoV1) map[string]string {
	routeTablesMap := make(map[string]string)
	for _, routeTable := range routeTablesFromDB {
		routeTablesMap[routeTable.Name] = routeTable.ID
	}
	return routeTablesMap
}
