package vpcs

import (
	"github.com/deployment-io/deployment-runner-kit/enums/region_enums"
	"github.com/deployment-io/deployment-runner-kit/enums/vpc_enums"
	"github.com/deployment-io/deployment-runner-kit/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func GetSubnetsMapFromDto(subnetsFromDto []SubnetDtoV1) map[string]string {
	subnetsMap := make(map[string]string)
	for _, subnet := range subnetsFromDto {
		subnetsMap[subnet.Name] = subnet.ID
	}
	return subnetsMap
}

func GetPrivateSubnetIdsFromDto(subnetsFromDto []SubnetDtoV1) primitive.A {
	var privateSubnets primitive.A
	for _, subnet := range subnetsFromDto {
		if subnet.IsPrivate == True {
			privateSubnets = append(privateSubnets, subnet.ID)
		}
	}
	return privateSubnets
}

func GetPublicSubnetIdsFromDto(subnetsFromDto []SubnetDtoV1) primitive.A {
	var publicSubnets primitive.A
	for _, subnet := range subnetsFromDto {
		if subnet.IsPrivate == False {
			publicSubnets = append(publicSubnets, subnet.ID)
		}
	}
	return publicSubnets
}

func GetRouteTablesMapFromDto(routeTablesFromDB []RouteTableDtoV1) map[string]string {
	routeTablesMap := make(map[string]string)
	for _, routeTable := range routeTablesFromDB {
		routeTablesMap[routeTable.Name] = routeTable.ID
	}
	return routeTablesMap
}
