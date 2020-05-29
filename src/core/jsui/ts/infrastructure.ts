export interface Server {
    Id: number
    OrgId: number
    Name: string
    Description: string
    OperatingSystem: string
    Location: string
    IpAddress: string
}

export interface ServerHandle {
    Id: number
    OrgId: number
}

export function extractServerHandle(f : Server) : ServerHandle {
    return {
        Id: f.Id,
        OrgId: f.OrgId,
    }
}

export interface ServerSSHConnectionGeneric {
    Id: number
    Username: string
}
