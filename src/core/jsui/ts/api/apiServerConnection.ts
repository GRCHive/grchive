import axios from 'axios'
import * as qs from 'query-string'
import {
    postFormJson,
    putFormJson,
    deleteFormJson,
    postFormMultipart,
    putFormMultipart
} from '../http'
import {
    apiv2SingleServerConnectionSSHPassword,
    apiv2ServerConnectionSSHPassword,
    apiv2SingleServerConnectionSSHKey,
    apiv2ServerConnectionSSHKey,
    apiv2ServerConnection,
} from '../url'
import {
    ServerSSHConnectionGeneric,
    ServerSSHPasswordConnection,
    ServerSSHKeyConnection
} from '../infrastructure'
import { getAPIRequestConfig } from './apiUtility'

export interface TNewServerSSHPasswordConnectionInput {
    orgId: number
    serverId: number
    username: string
    password: string
}

export interface TNewServerSSHConnectionOutput {
    data: ServerSSHConnectionGeneric
}

export function newServerSSHPasswordConnection(inp : TNewServerSSHPasswordConnectionInput) : Promise<TNewServerSSHConnectionOutput> {
    return postFormJson(
        apiv2ServerConnectionSSHPassword(inp.orgId, inp.serverId),
        inp,
        getAPIRequestConfig(),
    )
}

export interface TNewServerSSHKeyConnectionInput {
    orgId: number
    serverId: number
    username: string
    file: File
}

export function newServerSSHKeyConnection(inp : TNewServerSSHKeyConnectionInput) : Promise<TNewServerSSHConnectionOutput> {
    let data = new FormData()
    data.set('file', inp.file)
    data.set('username', inp.username)

    return postFormMultipart(
        apiv2ServerConnectionSSHKey(inp.orgId, inp.serverId),
        data,
        getAPIRequestConfig(),
    )
}

export interface TNewServerSSHConnectionOutput {
    data: ServerSSHConnectionGeneric
}

export interface TEditServerSSHPasswordConnectionInput extends TNewServerSSHPasswordConnectionInput{
    connectionId : number
}

export function editServerSSHPasswordConnection(inp : TEditServerSSHPasswordConnectionInput) : Promise<TNewServerSSHConnectionOutput> {
    return putFormJson(
        apiv2SingleServerConnectionSSHPassword(inp.orgId, inp.serverId, inp.connectionId),
        inp,
        getAPIRequestConfig(),
    )
}

export interface TEditServerSSHKeyConnectionInput extends TNewServerSSHKeyConnectionInput{
    connectionId : number
}

export function editServerSSHKeyConnection(inp : TEditServerSSHKeyConnectionInput) : Promise<TNewServerSSHConnectionOutput> {
    let data = new FormData()
    data.set('file', inp.file)
    data.set('username', inp.username)

    return putFormMultipart(
        apiv2SingleServerConnectionSSHKey(inp.orgId, inp.serverId, inp.connectionId),
        data,
        getAPIRequestConfig(),
    )
}

export interface TDeleteServerConnectionInput {
    orgId: number
    serverId: number
    connectionId: number
}

export function deleteServerSSHPasswordConnection(inp : TDeleteServerConnectionInput) : Promise<void> {
    return deleteFormJson(
        apiv2SingleServerConnectionSSHPassword(inp.orgId, inp.serverId, inp.connectionId),
        inp,
        getAPIRequestConfig(),
    )
}

export function deleteServerSSHKeyConnection(inp : TDeleteServerConnectionInput) : Promise<void> {
    return deleteFormJson(
        apiv2SingleServerConnectionSSHKey(inp.orgId, inp.serverId, inp.connectionId),
        inp,
        getAPIRequestConfig(),
    )
}

export interface TGetServerConnectionInput {
    orgId: number
    serverId: number
    connectionId: number
}


export interface TGetServerConnectionOutput {
    data: ServerSSHPasswordConnection
}

export function getServerSSHPasswordConnection(inp : TGetServerConnectionInput) : Promise<TGetServerConnectionOutput> {
    return axios.get(
        apiv2SingleServerConnectionSSHPassword(inp.orgId, inp.serverId, inp.connectionId),
        getAPIRequestConfig(),
    )
}

export interface TGetServerKeyConnectionOutput {
    data: ServerSSHKeyConnection
}

export function getServerSSHKeyConnection(inp : TGetServerConnectionInput) : Promise<TGetServerKeyConnectionOutput> {
    return axios.get(
        apiv2SingleServerConnectionSSHKey(inp.orgId, inp.serverId, inp.connectionId),
        getAPIRequestConfig(),
    )
}

export interface TAllServerConnectionInput {
    orgId: number
    serverId: number
}

export interface TAllServerConnectionOutput {
    data: {
        SshPassword: ServerSSHConnectionGeneric | null
        SshKey: ServerSSHConnectionGeneric | null
    }
}

export function getAllServerConnections(inp : TAllServerConnectionInput) : Promise<TAllServerConnectionOutput> {
    return axios.get(
        apiv2ServerConnection(inp.orgId, inp.serverId),
        getAPIRequestConfig()
    )
}
