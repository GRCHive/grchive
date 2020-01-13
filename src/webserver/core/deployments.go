package core

const (
	KNoDeployment     int32 = -1
	KSelfDeployment   int32 = 0
	KVendorDeployment int32 = 1
)

type SelfDeployment struct {
	Servers []*Server
}

type VendorDeployment struct {
	Product *VendorProduct
}

type FullDeployment struct {
	Id               int64 `db:"id"`
	OrgId            int32 `db:"org_id"`
	DeploymentType   int32 `db:"deployment_type"`
	SelfDeployment   *SelfDeployment
	VendorDeployment *VendorDeployment
}

type StrippedSelfDeployment struct {
	Servers []*ServerHandle
}

type StrippedVendorDeployment struct {
	VendorId  int64
	ProductId int64
}

type StrippedFullDeployment struct {
	Id               int64
	OrgId            int32
	DeploymentType   int32
	SelfDeployment   *StrippedSelfDeployment
	VendorDeployment *StrippedVendorDeployment
}
