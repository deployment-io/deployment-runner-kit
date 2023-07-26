package vpcs

import (
	"github.com/deployment-io/deployment-runner-kit/enums/region_enums"
	"github.com/deployment-io/deployment-runner-kit/enums/vpc_enums"
	"github.com/deployment-io/deployment-runner-kit/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type NatGatewayDtoV1 struct {
	Name                    string
	ID                      string
	ElasticIPAllocationName string
	ElasticIPAllocationId   string
}

type DefaultSecurityIngressRuleDtoV1 struct {
	ID string
}

type UpsertVpcDtoV1 struct {
	Type                        vpc_enums.Type
	Region                      region_enums.Type
	VpcId                       string
	VpcCidr                     string
	InternetGatewayId           string
	NatGateways                 []NatGatewayDtoV1
	Subnets                     []SubnetDtoV1
	RouteTables                 []RouteTableDtoV1
	DefaultSecurityIngressRules []DefaultSecurityIngressRuleDtoV1
	DefaultSecurityGroupId      string
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
		if subnet.IsPrivate == types.True {
			privateSubnets = append(privateSubnets, subnet.ID)
		}
	}
	return privateSubnets
}

func GetPublicSubnetIdsFromDto(subnetsFromDto []SubnetDtoV1) primitive.A {
	var publicSubnets primitive.A
	for _, subnet := range subnetsFromDto {
		if subnet.IsPrivate == types.False {
			publicSubnets = append(publicSubnets, subnet.ID)
		}
	}
	return publicSubnets
}

func GetRouteTablesMapFromDto(routeTablesFromDto []RouteTableDtoV1) map[string]string {
	routeTablesMap := make(map[string]string)
	for _, routeTable := range routeTablesFromDto {
		routeTablesMap[routeTable.Name] = routeTable.ID
	}
	return routeTablesMap
}

func GetNatGatewaysMapFromDto(natGatewaysFromDto []NatGatewayDtoV1) map[string]string {
	natGatewaysMap := make(map[string]string)
	for _, n := range natGatewaysFromDto {
		natGatewaysMap[n.Name] = n.ID
	}
	return natGatewaysMap
}

func GetAllocationsMapFromDto(natGatewaysFromDto []NatGatewayDtoV1) map[string]string {
	allocationsMap := make(map[string]string)
	for _, n := range natGatewaysFromDto {
		allocationsMap[n.ElasticIPAllocationName] = n.ElasticIPAllocationId
	}
	return allocationsMap
}
