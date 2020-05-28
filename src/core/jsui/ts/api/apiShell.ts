import axios from 'axios'
import {
    postFormJson,
    deleteFormJson
} from '../http'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import { 
    createOrgApiv2Url,
    allShellScriptsUrl,
    allShellScriptVersionsUrl,
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


export interface TAllShellScriptVersionsInput {
    orgId: number
    shellId: number
}

export interface TAllShellScriptVersionsOutput {
    data: ShellScriptVersion[]
}

export function allShellScriptVersions(inp : TAllShellScriptVersionsInput) : Promise<TAllShellScriptVersionsOutput> {
    return axios.get(
        createOrgApiv2Url(inp.orgId, allShellScriptVersionsUrl) + '?' + qs.stringify(inp),
        getAPIRequestConfig()
    ).then((resp : TAllShellScriptVersionsOutput) => {
        resp.data.forEach(cleanShellScriptVersionFromJson)
        return resp
    })
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
        createOrgApiv2Url(inp.orgId, allShellScriptsUrl + `/${inp.shellId}`),
        inp,
        getAPIRequestConfig()
    )
}
