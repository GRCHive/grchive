import axios from 'axios'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import { postFormJson } from '../http'
import { 
    saveCodeUrl,
    allCodeUrl,
    getCodeUrl,
    getCodeBuildStatusUrl,
    runCodeUrl,
    allCodeRunsUrl,
    getCodeRunUrl,
    getCodeLinkUrl,
} from '../url'
import { 
    ManagedCode,
    cleanManagedCodeFromJson,
    CodeParamType,
    DroneCiStatus,
    cleanDroneCiStatusFromJson,
} from '../code'
import {
    FullClientDataWithLink,
    ClientData
} from '../clientData'
import { ClientScript } from '../clientScripts'
import { 
    ScriptRun, cleanScriptRunFromJson
} from '../code'
import { ScheduledEvent } from '../event'

export interface TSaveCodeInput {
    orgId: number
    code : string
    dataId? : number
    scriptId? : number
    scriptData? : {
        params : (CodeParamType | null)[],
        clientDataId: number[],
    }
}

export interface TSaveCodeOutput {
    data : ManagedCode
}

export function saveCode(inp : TSaveCodeInput) : Promise<TSaveCodeOutput> {
    return postFormJson<TSaveCodeOutput>(saveCodeUrl, inp, getAPIRequestConfig()).then((resp : TSaveCodeOutput) => {
        cleanManagedCodeFromJson(resp.data)
        return resp
    })
}

export interface TAllCodeInput {
    orgId: number
    dataId? : number
    scriptId? : number
}

export interface TAllCodeOutput {
    data : ManagedCode[]
}

export function allCode(inp : TAllCodeInput) : Promise<TAllCodeOutput> {
    return axios.get(allCodeUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TAllCodeOutput) => {
        resp.data.forEach(cleanManagedCodeFromJson)
        return resp
    })
}

export interface TGetCodeInput {
    orgId: number
    dataId? : number
    scriptId? : number

    codeId? : number
    codeCommit? : string
}

export interface TGetCodeOutput {
    data: {
        Code: string
        ScriptData? : {
            Params : (CodeParamType | null)[],
            ClientData: FullClientDataWithLink[],
        },
        Full: ManagedCode
    }
}

export function getCode(inp : TGetCodeInput) : Promise<TGetCodeOutput> {
    return axios.get(getCodeUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TGetCodeOutput) => {
        cleanManagedCodeFromJson(resp.data.Full)
        return resp
    })
}


export interface TGetCodeBuildStatusInput {
    orgId: number
    commitHash : string
}

export interface TGetCodeBuildStatusOutput {
    data: DroneCiStatus | null
}

export function getCodeBuildStatus(inp : TGetCodeBuildStatusInput) : Promise<TGetCodeBuildStatusOutput> {
    return axios.get(getCodeBuildStatusUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TGetCodeBuildStatusOutput) => {
        if (!!resp.data) {
            cleanDroneCiStatusFromJson(resp.data)
        }
        return resp
    })
}

export interface TRunCodeInput {
    orgId: number
    codeId: number
    latest: boolean
    params: Record<number, any>
    schedule : ScheduledEvent | null
}

export interface TRunCodeOutput {
    data: number
}

export function runCode(inp : TRunCodeInput) : Promise<TRunCodeOutput> {
    return postFormJson<TRunCodeOutput>(runCodeUrl, inp, getAPIRequestConfig())
}

export interface TAllCodeRunsInput {
    orgId: number
    scriptId? : number
}

export interface TAllCodeRunsOutput {
    data: ScriptRun[]
}

export function allCodeRuns(inp : TAllCodeRunsInput) : Promise<TAllCodeRunsOutput> {
    return axios.get(allCodeRunsUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TAllCodeRunsOutput) => {
        resp.data.forEach(cleanScriptRunFromJson)
        return resp
    })
}

export interface TGetCodeRunInput {
    orgId: number
    runId : number
}

export interface TGetCodeRunOutput {
    data: ScriptRun
}

export function getCodeRun(inp : TGetCodeRunInput) : Promise<TGetCodeRunOutput> {
    return axios.get(getCodeRunUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TGetCodeRunOutput) => {
        cleanScriptRunFromJson(resp.data)
        return resp
    })
}

export interface TGetCodeLinkInput {
    orgId: number
    codeId : number
}

export interface TGetCodeLinkOutput {
    data: {
        Data: ClientData | null,
        Script: ClientScript | null
    }
}

export function getCodeLink(inp : TGetCodeLinkInput) : Promise<TGetCodeLinkOutput> {
    return axios.get(getCodeLinkUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
