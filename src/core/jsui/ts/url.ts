// TODO: Replace this whole file with some way of getting metadata from the server
// and just using stuff from the metadata...Maybe?
export const getStartedUrl : string = "/getting-started"
export const contactUsUrl : string = "/contact-us"
export const homePageUrl : string = "/"
export const loginPageUrl : string = "/login"
export const registerPageUrl : string = "/register"


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

export function createMyAccountUrl(id : number) : string {
    return `${myAccountBaseUrl}/${id}`
}

export function createMyProfileUrl(id : number) : string {
    return `${myAccountBaseUrl}/${id}/profile`
}

export function createMyOrgsUrl(id : number) : string {
    return `${myAccountBaseUrl}/${id}/orgs`
}

export function createOrgUrl(org : string) : string {
    return `/dashboard/org/${org}`
}

export function createOrgGLUrl(org : string ): string {
    return `${createOrgUrl(org)}/gl`
}

export function createOrgSystemUrl(org : string) : string {
    return `${createOrgUrl(org)}/it/systems`
}

export function createOrgDatabaseUrl(org : string) : string {
    return `${createOrgUrl(org)}/it/databases`
}

export function createSingleSystemUrl(org: string, sys : number) : string {
    return `/dashboard/org/${org}/it/systems/${sys}`
}

export function createSingleDbUrl(org: string, db : number) : string {
    return `/dashboard/org/${org}/it/databases/${db}`
}

export function createSingleInfraUrl(org: string, infra : number) : string {
    return `/dashboard/org/${org}/it/infrastructure/${infra}`
}

export function createFlowUrl(org : string, flow : number) : string {
    return `/dashboard/org/${org}/flows/${flow}`
}

export function createRiskUrl(org : string, risk : number) : string {
    return `/dashboard/org/${org}/risks/${risk}`
}

export function createControlUrl(org : string, control : number) : string {
    return `/dashboard/org/${org}/controls/${control}`
}

export function createOrgAllRolesUrl(org : string) : string {
    return `/dashboard/org/${org}/settings/roles`
}

export function createOrgRoleUrl(org : string, role : number) : string {
    return `/dashboard/org/${org}/settings/roles/${role}`
}

export function createSingleDocCatUrl(org : string, id: number) : string {
    return `/dashboard/org/${org}/documentation/cat/${id}`
}

export function createUserProfileEditAPIUrl(id : number) : string {
    return `/api/user/${id}/profile`
}


export function createUserGetOrgsAPIUrl(id : number) : string {
    return `/api/user/${id}/orgs`
}

export const newProcessFlowAPIUrl : string = "/api/flows/new"
export const deleteProcessFlowAPIUrl : string = "/api/flows/delete"
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
export const deleteRiskAPIUrl : string = "/api/risk/delete"
export const editRiskAPIUrl : string = "/api/risk/edit"
export const addExistingRiskAPIUrl : string = "/api/risk/add"
export const allRiskAPIUrl : string = "/api/risk/"

export function createSingleRiskAPIUrl(riskId : number) : string {
    return `/api/risk/${riskId}`
}

export function createGetAllOrgUsersAPIUrl(org : string) : string {
    return `/api/org/${org}/users`
}

export const getControlTypesUrl : string = "/api/control/types"
export const newControlUrl : string = "/api/control/new"
export const deleteControlUrl : string = "/api/control/delete"
export const addControlUrl : string = "/api/control/add"
export const editControlUrl : string = "/api/control/edit"
export const linkCatControlUrl : string = "/api/control/linkCat"
export const unlinkCatControlUrl : string = "/api/control/unlinkCat"
export const allControlAPIUrl : string = "/api/control/"
export function createSingleControlAPIUrl(controlId : number) : string {
    return `/api/control/${controlId}`
}

export const newControlDocCatUrl : string = "/api/documentation/cat/new"
export const editControlDocCatUrl : string = "/api/documentation/cat/edit"
export const deleteControlDocCatUrl : string = "/api/documentation/cat/delete"
export const allControlDocCatUrl : string = "/api/documentation/cat/all"
export const getControlDocCatUrl : string = "/api/documentation/cat/get"

export const uploadControlDocUrl : string = "/api/documentation/file/upload"
export const getControlDocUrl : string = "/api/documentation/file/get"
export const deleteControlDocUrl : string = "/api/documentation/file/delete"
export const downloadControlDocUrl : string = "/api/documentation/file/download"

export const requestVerificationEmailUrl : string = "/api/verification/resend"
export const inviteUsersToOrgUrl: string = "/api/invite/send"

export const getOrgRolesUrl : string = "/api/roles/all"
export const getSingleOrgRoleUrl : string = "/api/roles/get"
export const newRoleUrl : string = "/api/roles/new"
export const editRoleUrl : string = "/api/roles/edit"
export const deleteRoleUrl : string = "/api/roles/delete"
export const addUsersToRoleUrl : string = "/api/roles/addUsers"

export const getGLUrl : string = "/api/gl/get"
export const createNewGLCatUrl : string = "/api/gl/cat/new"
export const editGLCatUrl : string = "/api/gl/cat/edit"
export const deleteGLCatUrl : string = "/api/gl/cat/delete"
export const createNewGLAccUrl : string = "/api/gl/acc/new"
export const editGLAccUrl : string = "/api/gl/acc/edit"
export const getGLAccUrl : string = "/api/gl/acc/get"
export const deleteGLAccUrl : string = "/api/gl/acc/delete"
export function createSingleGLAccountUrl(org : string, accId : number) : string {
    return `/dashboard/org/${org}/gl/acc/${accId}`
}

export const newSystemUrl : string = "/api/it/systems/new"
export const allSystemsUrl : string = "/api/it/systems/all"
export const editSystemUrl : string = "/api/it/systems/edit"
export const deleteSystemUrl : string = "/api/it/systems/delete"
export const getSystemUrl : string = "/api/it/systems/get"

export const newDatabaseUrl : string = "/api/it/db/new"
export const allDatabaseUrl : string = "/api/it/db/all"
export const typesDatabaseUrl : string = "/api/it/db/types"
export const editDatabaseUrl : string = "/api/it/db/edit"
export const deleteDatabaseUrl : string = "/api/it/db/delete"
export const getDatabaseUrl : string = "/api/it/db/get"

export const newDbConnUrl : string = "/api/it/db/connection/new"
export const deleteDbConnUrl : string = "/api/it/db/connection/delete"

export const linkDbsToSystemUrl : string = "/api/it/systems/linkDb"
export const linkSystemsToDbUrl : string = "/api/it/db/linkSys"
export const deleteDbSysLinkUrl : string = "/api/it/deleteDbSysLink"
