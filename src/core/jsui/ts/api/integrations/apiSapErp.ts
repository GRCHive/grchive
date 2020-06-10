import axios from 'axios'
import {
    postFormJson,
    putFormJson,
} from '../../http'
import { getAPIRequestConfig } from '../apiUtility'
import { GenericIntegration } from '../../integrations/integration'
import { SapErpIntegrationSetup } from '../../integrations/sap'
import { 
    sapErpIntegrationBaseUrl,
    apiv2SingleSapErpIntegrationUrl,
} from '../../url'
import { TIntegrationLink, createApiv2IntegrationUrl } from './apiIntegrations'

export interface TNewSapErpIntegrationInput {
    orgId : number
    setup : SapErpIntegrationSetup
    integration : GenericIntegration
    link : TIntegrationLink
}

export interface TNewSapErpIntegrationOutput {
    data: GenericIntegration
}

export function newSapErpIntegration(inp  : TNewSapErpIntegrationInput) : Promise<TNewSapErpIntegrationOutput> {
    return postFormJson(
        createApiv2IntegrationUrl(inp.orgId, inp.link, sapErpIntegrationBaseUrl),
        inp,
        getAPIRequestConfig()
    )
}

export interface TGetSapErpIntegrationInput {
    orgId : number
    integrationId: number
}

export interface TGetSapErpIntegrationOutput {
    data: SapErpIntegrationSetup
}

export function getSapErpIntegration(inp : TGetSapErpIntegrationInput) : Promise<TGetSapErpIntegrationOutput> {
    return axios.get(
        apiv2SingleSapErpIntegrationUrl(inp.orgId, inp.integrationId),
        getAPIRequestConfig(),
    )
}

export interface TEditSapErpIntegrationInput {
    orgId : number
    integrationId: number
    setup : SapErpIntegrationSetup
}

export interface TEditSapErpIntegrationOutput {
    data: SapErpIntegrationSetup
}

export function editSapErpIntegration(inp : TEditSapErpIntegrationInput) : Promise<TEditSapErpIntegrationOutput> {
    return putFormJson(
        apiv2SingleSapErpIntegrationUrl(inp.orgId, inp.integrationId),
        inp,
        getAPIRequestConfig(),
    )
}
