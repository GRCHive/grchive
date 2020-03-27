import axios from 'axios'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import { postFormJson } from '../http'
import {
    FullClientDataWithLink,
    DataSourceOption  
} from '../clientData'
import { 
    newClientDataUrl,
    allClientDataUrl,
    allDataSourceUrl,
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

export interface TAllClientDataInput {
    orgId: number
}

export interface TAllClientDataOutput {
    data: FullClientDataWithLink[]
}

export function allClientData(inp : TAllClientDataInput) : Promise<TAllClientDataOutput> {
    return axios.get(allClientDataUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TAllDataSourceOutput {
    data: DataSourceOption[]
}

export function allSupportedDataSources() : Promise<TAllDataSourceOutput> {
    return axios.get(allDataSourceUrl, getAPIRequestConfig())
}

export interface TDeleteClientDataInput {
    orgId: number
    dataId: number
}

export function deleteClientData(inp : TDeleteClientDataInput) : Promise<void> {
    return postFormJson<void>(deleteClientDataUrl, inp, getAPIRequestConfig())
}
