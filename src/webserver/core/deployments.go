package core

const (
	KNoDeployment     int32 = -1
	KSelfDeployment   int32 = 0
	KVendorDeployment int32 = 1
)

type SelfDeployment struct {
}

type VendorDeployment struct {
	VendorName    string `db:"vendor_name"`
	VendorProduct string `db:"vendor_product"`
	SocFiles      []*ControlDocumentationFile
}

type FullDeployment struct {
	Id               int64 `db:"id"`
	OrgId            int32 `db:"org_id"`
	DeploymentType   int32 `db:"deployment_type"`
	SelfDeployment   *SelfDeployment
	VendorDeployment *VendorDeployment
}

type StrippedSelfDeployment struct {
}

type StrippedVendorDeployment struct {
	VendorName    string
	VendorProduct string
	SocFiles      []*ControlDocumentationFileHandle
}

type StrippedFullDeployment struct {
	Id               int64
	OrgId            int32
	DeploymentType   int32
	SelfDeployment   *StrippedSelfDeployment
	VendorDeployment *StrippedVendorDeployment
}
