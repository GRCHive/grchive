import axios from 'axios'
import { 
    createOrgApiv2Url,
    integrationBaseUrl,
    apiv2SingleIntegrationUrl,
} from '../../url'
import {
    putFormJson,
    deleteFormJson,
} from '../../http'
import { GenericIntegration } from '../../integrations/integration'
import { getAPIRequestConfig } from '../apiUtility'

export interface TIntegrationLink {
    systemId: number | undefined
}

export function createApiv2IntegrationUrl(orgId : number, link : TIntegrationLink, endpoint : string) : string {
    let linkStr : string = ""
    if (!!link.systemId) {
        linkStr = `system/${link.systemId}`
    }

    return createOrgApiv2Url(
        orgId,
        `${linkStr}/${endpoint}`)
}

export interface TAllGenericIntegrationsInput {
    orgId: number
    link : TIntegrationLink
}

export interface TAllGenericIntegrationsOutput {
    data: GenericIntegration[]
}

export function allGenericIntegrations(inp : TAllGenericIntegrationsInput) : Promise<TAllGenericIntegrationsOutput> {
    return axios.get(
        createApiv2IntegrationUrl(inp.orgId, inp.link, integrationBaseUrl),
        getAPIRequestConfig(),
    )
}

export interface TEditGenericIntegrationInput {
    orgId: number
    integrationId: number
    data: GenericIntegration
}

export interface TEditGenericIntegrationOutput {
    data: GenericIntegration
}

export function editGenericIntegration(inp : TEditGenericIntegrationInput) : Promise<TEditGenericIntegrationOutput> {
    return putFormJson(
        apiv2SingleIntegrationUrl(inp.orgId, inp.integrationId),
        inp,
        getAPIRequestConfig()
    )
}

export interface TDeleteGenericIntegrationInput {
    orgId: number
    integrationId: number
}

export function deleteGenericIntegration(inp : TDeleteGenericIntegrationInput) : Promise<void> {
    return deleteFormJson(
        apiv2SingleIntegrationUrl(inp.orgId, inp.integrationId),
        {},
        getAPIRequestConfig()
    )
}
