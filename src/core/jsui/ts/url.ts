// TODO: Replace this whole file with some way of getting metadata from the server
// and just using stuff from the metadata...
export const getStartedUrl : string = "/getting-started"
export const contactUsUrl : string = "/contact-us"
export const homePageUrl : string = "/"
export const loginPageUrl : string = "/login"

export function createAssetUrl(asset : string) : string {
    return "/static/assets/" + asset
}

export const learnMoreUrl : string = "/learn"
export const dashboardUrl : string = "/dashboard"

export function createMailtoUrl(user : string, domain : string) : Object {
    const email = createEmailAddress(user, domain)
    return Object.freeze({
        mailto: "mailto:" + email,
        email: email
    })
}

export function createEmailAddress(user : string, domain : string) : string {
    return user + "@" + domain
}

export const baseLogoutUrl : string = "/logout";
export function createLogoutUrl(csrf : string) : string {
    return baseLogoutUrl + "?csrf=" + csrf
}

export const myAccountBaseUrl : string = "/dashboard/user";
export function createMyAccountUrl(email: string) : string {
    return myAccountBaseUrl + "/" + email
}

export function createUserProfileEditAPIUrl(email : string) : string {
    return "/api/user/" + email + "/profile"
}

export const newProcessFlowAPIUrl : string = "/api/flows/new"
export const getAllProcessFlowAPIUrl: string = "/api/flows/"
export const getAllProcessFlowNodeTypesAPIUrl: string = "/api/flownodes/types"
export const getAllProcessFlowIOTypesAPIUrl: string = "/api/flowio/types"
export const newProcessFlowIOAPIUrl: string = "/api/flowio/new"
export const deleteProcessFlowIOAPIUrl: string = "/api/flowio/delete"
export const editProcessFlowIOAPIUrl: string = "/api/flowio/edit"
export function createGetProcessFlowFullDataUrl(id : number) : string {
    return "/api/flows/" + id.toString() + "/full"
}

export const newProcessFlowNodeAPIUrl: string = "/api/flownodes/new"
export const editProcessFlowNodeAPIUrl: string = "/api/flownodes/edit"
export const deleteProcessFlowNodeAPIUrl: string = "/api/flownodes/delete"

export function createUpdateProcessFlowApiUrl(id : number) : string {
    return "/api/flows/" + id.toString() + "/update"
}

export const newProcessFlowEdgeAPIUrl: string = "/api/flowedges/new"
export const deleteProcessFlowEdgeAPIUrl: string = "/api/flowedges/delete"

export function createProcessFlowNodeDisplaySettingsWebsocket(host : string, csrf : string, flowId: number) : string {
    return `ws://${host}/ws/flownodedisp/${flowId.toString()}?csrf=${csrf}`
}

export const newRiskAPIUrl : string = "/api/risk/new"
