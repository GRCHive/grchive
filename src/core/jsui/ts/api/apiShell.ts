import axios from 'axios'
import {
    postFormJson,
    deleteFormJson,
    putFormJson,
} from '../http'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import { 
    createOrgApiv2Url,
    allShellScriptsUrl,
    singleShellScriptUrl,
    singleShellScriptVersionUrl,
} from '../url'
import {
    ShellScript,
    ShellScriptVersion,
    cleanShellScriptVersionFromJson
} from '../shell'

export interface TAllShellScriptsInput {
    orgId: number
    shellType : number
}

export interface TAllShellScriptsOutput {
    data: ShellScript[]
}

export function allShellScripts(inp : TAllShellScriptsInput) : Promise<TAllShellScriptsOutput> {
    return axios.get(
        createOrgApiv2Url(inp.orgId, allShellScriptsUrl) + '?' + qs.stringify(inp),
        getAPIRequestConfig()
    )
}


export interface TNewShellScriptInput {
    orgId: number
    shellType : number
    name: string
    description: string
    script : string
}


export interface TNewShellScriptOutput {
    data: ShellScript
}

export function newShellScript(inp : TNewShellScriptInput) : Promise<TNewShellScriptOutput> {
    return postFormJson(
        createOrgApiv2Url(inp.orgId, allShellScriptsUrl),
        inp,
        getAPIRequestConfig()
    )
}

export interface TDeleteShellScriptInput {
    orgId: number
    shellId: number
}

export function deleteShellScript(inp : TDeleteShellScriptInput) : Promise<void> {
    return deleteFormJson(
        singleShellScriptUrl(inp.orgId, inp.shellId),
        inp,
        getAPIRequestConfig()
    )
}

export interface TGetShellScriptInput {
    orgId: number
    shellId: number
}


export interface TGetShellScriptOutput {
    data: ShellScript
}

export function getShellScript(inp : TGetShellScriptInput) : Promise<TGetShellScriptOutput> {
    return axios.get(
        singleShellScriptUrl(inp.orgId, inp.shellId) + '?' + qs.stringify(inp),
        getAPIRequestConfig(),
    )
}

export interface TEditShellScriptInput {
    orgId: number
    shellId: number
    name : string
    description: string
}

export function editShellScript(inp : TEditShellScriptInput) : Promise<TNewShellScriptOutput> {
    return putFormJson(
        singleShellScriptUrl(inp.orgId, inp.shellId),
        inp,
        getAPIRequestConfig(),
    )
}

export interface TAllShellScriptVersionsInput {
    orgId: number
    shellId: number
}

export interface TAllShellScriptVersionsOutput {
    data: ShellScriptVersion[]
}

export function allShellScriptVersions(inp : TAllShellScriptVersionsInput) : Promise<TAllShellScriptVersionsOutput> {
    return axios.get(
        singleShellScriptUrl(inp.orgId, inp.shellId) + '/version?' + qs.stringify(inp),
        getAPIRequestConfig()
    ).then((resp : TAllShellScriptVersionsOutput) => {
        resp.data.forEach(cleanShellScriptVersionFromJson)
        return resp
    })
}

export interface TGetShellScriptVersionInput {
    orgId: number
    shellId: number
    version: number
}

export interface TGetShellScriptVersionOutput {
    data: string
}

export function getShellScriptVersion(inp : TGetShellScriptVersionInput) : Promise<TGetShellScriptVersionOutput> {
    return axios.get(
        singleShellScriptVersionUrl(inp.orgId, inp.shellId, inp.version),
        getAPIRequestConfig()
    )
}

export interface TNewShellScriptVersionInput {
    orgId: number
    shellId: number
    script : string
}

export interface TNewShellScriptVersionOutput {
    data: ShellScriptVersion
}

export function newShellScriptVersion(inp : TNewShellScriptVersionInput) : Promise<TNewShellScriptVersionOutput> {
    return postFormJson<TNewShellScriptVersionOutput>(
        singleShellScriptUrl(inp.orgId, inp.shellId) + '/version',
        inp,
        getAPIRequestConfig(),
    ).then((resp : TNewShellScriptVersionOutput) => {
        cleanShellScriptVersionFromJson(resp.data)
        return resp
    })
}

