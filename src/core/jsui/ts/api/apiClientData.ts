import axios from 'axios'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import { postFormJson } from '../http'
import {
    FullClientDataWithLink,
} from '../clientData'
import { 
    newClientDataUrl,
    updateClientDataUrl,
    allClientDataUrl,
    getClientDataUrl,
    deleteClientDataUrl,
} from '../url'

export interface TNewClientDataInput {
    orgId: number
    name: string
    description: string
    sourceId : number
    sourceTarget: Record<string, any>
}

export interface TNewClientDataOutput {
    data: FullClientDataWithLink
}

export function newClientData(inp : TNewClientDataInput) : Promise<TNewClientDataOutput> {
    return postFormJson<TNewClientDataOutput>(newClientDataUrl, inp, getAPIRequestConfig())
}

export interface TUpdateClientDataInput extends TNewClientDataInput {
    dataId: number
}

export interface TUpdateClientDataOutput {
    data: FullClientDataWithLink
}

export function updateClientData(inp : TUpdateClientDataInput) : Promise<TUpdateClientDataOutput> {
    return postFormJson<TUpdateClientDataOutput>(updateClientDataUrl, inp, getAPIRequestConfig())
}

export interface TAllClientDataInput {
    orgId: number
}

export interface TAllClientDataOutput {
    data: FullClientDataWithLink[]
}

export function allClientData(inp : TAllClientDataInput) : Promise<TAllClientDataOutput> {
    return axios.get(allClientDataUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TGetClientDataInput {
    orgId: number
    dataId: number
}

export interface TGetClientDataOutput {
    data: FullClientDataWithLink
}

export function getClientData(inp : TGetClientDataInput) : Promise<TGetClientDataOutput> {
    return axios.get(getClientDataUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TDeleteClientDataInput {
    orgId: number
    dataId: number
}

export function deleteClientData(inp : TDeleteClientDataInput) : Promise<void> {
    return postFormJson<void>(deleteClientDataUrl, inp, getAPIRequestConfig())
}
