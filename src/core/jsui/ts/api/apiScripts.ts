import axios from 'axios'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import { postFormJson } from '../http'
import {
    ClientScript,
} from '../clientScripts'
import {
    ManagedCode
} from '../code'
import { 
    newClientScriptUrl,
    allClientScriptsUrl,
    deleteClientScriptUrl,
    getClientScriptUrl,
    updateClientScriptUrl,
    getClientScriptCodeLinkUrl
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

export interface TUpdateClientScriptInput extends TNewClientScriptInput {
    scriptId: number
}

export interface TUpdateClientScriptOutput {
    data: ClientScript
}

export function updateClientScript(inp : TUpdateClientScriptInput) : Promise<TUpdateClientScriptOutput> {
    return postFormJson<TUpdateClientScriptOutput>(updateClientScriptUrl, inp, getAPIRequestConfig())
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

export interface TGetClientScriptInput {
    orgId: number
    scriptId: number
}

export interface TGetClientScriptOutput {
    data: ClientScript
}

export function getClientScript(inp : TGetClientScriptInput) : Promise<TGetClientScriptOutput> {
    return axios.get(getClientScriptUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TGetClientScriptCodeFromLinkInput {
    orgId: number
    linkId: number
}

export interface TGetClientScriptCodeFromLinkOutput {
    data: {
        Script: ClientScript
        Code: ManagedCode
    }
}

export function getClientScriptCodeFromLink(inp : TGetClientScriptCodeFromLinkInput) : Promise<TGetClientScriptCodeFromLinkOutput> {
    return axios.get(getClientScriptCodeLinkUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
