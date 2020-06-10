export interface SapErpIntegrationSetup {
    Client : string
    SysNr: string
    Host: string
    RealHostname: string | null
    Username: string
    Password: string
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
