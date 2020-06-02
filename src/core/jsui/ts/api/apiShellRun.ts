import axios from 'axios'
import * as qs from 'query-string'
import {
    postFormJson,
} from '../http'
import { 
    createOrgApiv2Url,
    allShellScriptRunsUrl,
    singleShellScriptVersionRunUrl,
    singleShellRunUrl,
} from '../url'
import { getAPIRequestConfig } from './apiUtility'
import {
    ShellScript,
    ShellScriptVersion,
    ShellScriptRun,
    ShellScriptRunPerServer,
    cleanShellScriptRunFromJson,
    cleanShellScriptVersionFromJson,
    cleanShellScriptRunPerServerFromJson,
} from '../shell'

export interface TRequestRunShellScriptInput {
    orgId: number
    shellId: number
    versionId : number
    servers : number[]
}

export interface TRequestRunShellScriptOutput {
    data: {
        RunId : number
        RequestId?: number
    }
}

export function requestRunShellScript(inp : TRequestRunShellScriptInput) : Promise<TRequestRunShellScriptOutput> {
    return postFormJson(
        singleShellScriptVersionRunUrl(inp.orgId, inp.shellId, inp.versionId),
        inp,
        getAPIRequestConfig()
    )
}

export interface TAllShellRunInput {
    orgId: number
    shellId? : number
    serverId? : number
}

export interface TAllShellRunOutput {
    data: ShellScriptRun[]
}

export function allShellRuns(inp : TAllShellRunInput): Promise<TAllShellRunOutput> {
    return axios.get(
        createOrgApiv2Url(inp.orgId, allShellScriptRunsUrl) + '?' + qs.stringify(inp),
        getAPIRequestConfig()
    ).then((resp : TAllShellRunOutput) => {
        resp.data.forEach(cleanShellScriptRunFromJson)
        return resp
    })
}

export interface TGetShellRunInput {
    orgId: number
    runId: number
    includeLogs: boolean
}

export interface TGetShellRunOutput {
    data: {
        Run: ShellScriptRun
        Script: ShellScript
        Version?: ShellScriptVersion
        VersionNum: number
        ServerRuns: ShellScriptRunPerServer[]
    }
}

export function getShellRunInformation(inp : TGetShellRunInput): Promise<TGetShellRunOutput> {
    return axios.get(
        singleShellRunUrl(inp.orgId, inp.runId) + '?'+ qs.stringify(inp),
        getAPIRequestConfig(),
    ).then((resp : TGetShellRunOutput) => {
        cleanShellScriptRunFromJson(resp.data.Run)
        if (!!resp.data.Version) {
            cleanShellScriptVersionFromJson(resp.data.Version)
        }
        resp.data.ServerRuns.forEach(cleanShellScriptRunPerServerFromJson)
        return resp
    })
}
