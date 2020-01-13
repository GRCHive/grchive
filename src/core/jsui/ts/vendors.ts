export interface Vendor {
    Id : number
    OrgId : number
    Name : string
    Description : string
    Url : string
    DocCatId: number
}

export interface VendorProduct {
    Id : number
    VendorId : number
    OrgId : number
    Name : string
    Description:  string
    Url : string
}
