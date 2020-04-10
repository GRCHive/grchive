import axios from 'axios'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import { postFormJson } from '../http'
import {
    ClientScript,
} from '../clientScripts'
import { 
    newClientScriptUrl,
    allClientScriptsUrl,
    deleteClientScriptUrl,
} from '../url'

export interface TNewClientScriptInput {
    orgId: number
    name: string
    description: string
}

export interface TNewClientScriptOutput {
    data: ClientScript
}

export function newClientScript(inp : TNewClientScriptInput) : Promise<TNewClientScriptOutput> {
    return postFormJson<TNewClientScriptOutput>(newClientScriptUrl, inp, getAPIRequestConfig())
}

export interface TAllClientScriptsInput {
    orgId: number
}
export interface TAllClientScriptsOutput {
    data: ClientScript[]
}

export function allClientScripts(inp : TAllClientScriptsInput) : Promise<TAllClientScriptsOutput> {
    return axios.get(allClientScriptsUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TDeleteClientScriptInput {
    orgId: number
    scriptId: number
}

export function deleteClientScript(inp : TDeleteClientScriptInput) : Promise<void> {
    return postFormJson<void>(deleteClientScriptUrl, inp, getAPIRequestConfig())
}
