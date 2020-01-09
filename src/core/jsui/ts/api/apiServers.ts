import axios from 'axios'
import * as qs from 'query-string'
import { postFormJson } from '../http'
import { getAPIRequestConfig } from './apiUtility'
import { Server } from '../infrastructure'
import { 
    newServerUrl,
    allServersUrl,
    getServerUrl,
    updateServerUrl,
    deleteServerUrl,
} from '../url'

export interface TNewServerInput {
    orgId: number
    name: string
    description: string
    ip : string
    os : string
    location : string
}

export interface TNewServerOutput {
    data: Server
}

export function newServer(inp : TNewServerInput) : Promise<TNewServerOutput> {
    return postFormJson<TNewServerOutput>(newServerUrl, inp, getAPIRequestConfig())
}

export interface TUpdateServerInput extends TNewServerInput {
    serverId: number
}

export function updateServer(inp : TUpdateServerInput) : Promise<TNewServerOutput> {
    return postFormJson<TNewServerOutput>(updateServerUrl, inp, getAPIRequestConfig())
}

export interface TAllServerInput {
    orgId: number
}

export interface TAllServerOutput {
    data: Server[]
}

export function allServers(inp : TAllServerInput) : Promise<TAllServerOutput> {
    return axios.get(allServersUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TGetServerInput {
    serverId: number
    orgId: number
}

export interface TGetServerOutput {
    data: {
        Server: Server
    }
}

export function getServer(inp : TGetServerInput) : Promise<TGetServerOutput> {
    return axios.get(getServerUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TDeleteServerInput {
    serverId: number
    orgId: number
}

export function deleteServer(inp : TDeleteServerInput) : Promise<void> {
    return postFormJson<void>(deleteServerUrl, inp, getAPIRequestConfig())
}
