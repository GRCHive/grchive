import axios from 'axios'
import {
    postFormJson,
    putFormJson,
    deleteFormJson,
} from '../../http'
import { getAPIRequestConfig } from '../apiUtility'
import { GenericIntegration } from '../../integrations/integration'
import { SapErpIntegrationSetup } from '../../integrations/sap'
import { 
    sapErpIntegrationBaseUrl,
    apiv2SingleSapErpIntegrationUrl,
    apiv2SingleSapErpIntegrationRfcUrl,
    apiv2SingleSapErpIntegrationSingleRfcUrl,
    apiv2SingleSapErpIntegrationSingleRfcVersionsUrl,
    apiv2SingleSapErpIntegrationSingleRfcSingleVersionUrl,
} from '../../url'
import { TIntegrationLink, createApiv2IntegrationUrl } from './apiIntegrations'
import {
    SapErpRfcMetadata,
    SapErpRfcVersion,
    SapErpRfcSettings,
    cleanSapErpRfcVersionFromJson,
} from '../../integrations/sap'

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

export interface TAllSapErpRfcInput {
    orgId : number
    integrationId: number
}

export interface TAllSapErpRfcOutput {
    data: SapErpRfcMetadata[]
}

export function allSapErpRfc(inp : TAllSapErpRfcInput) : Promise<TAllSapErpRfcOutput> {
    return axios.get(
        apiv2SingleSapErpIntegrationRfcUrl(inp.orgId, inp.integrationId),
        getAPIRequestConfig(),
    )
}

export interface TNewSapErpRfcInput {
    orgId : number
    integrationId: number
    function : string
}

export interface TNewSapErpRfcOutput {
    data: SapErpRfcMetadata
}

export function newSapErpRfc(inp : TNewSapErpRfcInput) : Promise<TNewSapErpRfcOutput> {
    return postFormJson(
        apiv2SingleSapErpIntegrationRfcUrl(inp.orgId, inp.integrationId),
        inp,
        getAPIRequestConfig(),
    )
}

export interface TDeleteSapErpRfcInput {
    orgId : number
    integrationId: number
    rfcId: number
}

export function deleteSapErpRfc(inp : TDeleteSapErpRfcInput) : Promise<void> {
    return deleteFormJson(
        apiv2SingleSapErpIntegrationSingleRfcUrl(inp.orgId, inp.integrationId, inp.rfcId),
        inp,
        getAPIRequestConfig(),
    )
}

export interface TAllSapErpRfcVersionsInput {
    orgId : number
    integrationId: number
    rfcId: number
}

export interface TAllSapErpRfcVersionsOutput {
    data: SapErpRfcVersion[]
}

export function allSapErpRfcVersions(inp : TAllSapErpRfcVersionsInput) : Promise<TAllSapErpRfcVersionsOutput> {
    return axios.get(
        apiv2SingleSapErpIntegrationSingleRfcVersionsUrl(inp.orgId, inp.integrationId, inp.rfcId),
        getAPIRequestConfig(),
    ).then((resp : TAllSapErpRfcVersionsOutput) => {
        resp.data.forEach(cleanSapErpRfcVersionFromJson)
        return resp
    })
}

export interface TNewSapErpRfcVersionInput {
    orgId : number
    integrationId: number
    rfcId: number
}

export interface TNewSapErpRfcVersionOutput {
    data: SapErpRfcVersion
}

export function newSapErpRfcVersion(inp : TNewSapErpRfcVersionInput) : Promise<TNewSapErpRfcVersionOutput> {
    return postFormJson<TNewSapErpRfcVersionOutput>(
        apiv2SingleSapErpIntegrationSingleRfcVersionsUrl(inp.orgId, inp.integrationId, inp.rfcId),
        {},
        getAPIRequestConfig(),
    ).then((resp : TNewSapErpRfcVersionOutput) => {
        cleanSapErpRfcVersionFromJson(resp.data)
        return resp
    })
}

export interface TGetSapErpRfcVersionInput {
    orgId : number
    integrationId: number
    rfcId: number
    versionId: number
}

export interface TGetSapErpRfcVersionOutput {
    data: SapErpRfcVersion
}

export function getSapErpRfcVersion(inp : TGetSapErpRfcVersionInput) : Promise<TGetSapErpRfcVersionOutput> {
    return axios.get(
        apiv2SingleSapErpIntegrationSingleRfcSingleVersionUrl(inp.orgId, inp.integrationId, inp.rfcId, inp.versionId),
        getAPIRequestConfig(),
    ).then((resp : TGetSapErpRfcVersionOutput) => {
        cleanSapErpRfcVersionFromJson(resp.data)
        return resp
    })
}
