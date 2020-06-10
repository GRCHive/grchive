export interface SapErpIntegrationSetup {
    Client : string
    SysNr: string
    Host: string
    RealHostname: string | null
    Username: string
    Password: string
}

export interface SapErpRfcMetadata {
    Id : number
    IntegrationId : number
    Function : string
}

export interface SapErpRfcVersion {
    Id: number
    RfcId : number
    CreatedTime: Date
    FinishedTime: Date | null
    Data: Object | null
    Success: boolean
    Logs: string | null
}

export function cleanSapErpRfcVersionFromJson(v : SapErpRfcVersion) {
    v.CreatedTime = new Date(v.CreatedTime)
    if (!!v.FinishedTime) {
        v.FinishedTime = new Date(v.FinishedTime)
    }
}

export interface SapErpRfcSettings {
}

export function emptySapErpIntegrationSetup() : SapErpIntegrationSetup {
    return {
        Client: "",
        SysNr: "",
        Host: "",
        RealHostname: null,
        Username: "",
        Password: "",
    }
}
