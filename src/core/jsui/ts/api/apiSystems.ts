import axios from 'axios'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import { postFormJson } from '../http'
import {
    newSystemUrl,
    allSystemsUrl,
    editSystemUrl,
    deleteSystemUrl,
    getSystemUrl
} from '../url'
import { System } from '../systems'

export interface TNewSystemInputs {
    orgId: number
    name: string
    purpose: string
    description: string
}

export interface TNewSystemOutputs {
    data: System
}

export function newSystem(inp : TNewSystemInputs) : Promise<TNewSystemOutputs> {
    return postFormJson<TNewSystemOutputs>(newSystemUrl, inp, getAPIRequestConfig())
}

export interface TAllSystemsInputs {
    orgId: number
}

export interface TAllSystemsOutputs {
    data: System[]
}

export function getAllSystems(inp : TAllSystemsInputs) : Promise<TAllSystemsOutputs> {
    return axios.get(allSystemsUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TEditSystemInputs extends TNewSystemInputs {
    sysId: number
}

export interface TEditSystemOutputs {
    data: System
}

export function editSystem(inp : TEditSystemInputs) : Promise<TEditSystemOutputs> {
    return postFormJson<TEditSystemOutputs>(editSystemUrl, inp, getAPIRequestConfig())
}

export interface TDeleteSystemInputs {
    sysId: number
    orgId: number
}

export interface TDeleteSystemOutputs {
}

export function deleteSystem(inp : TDeleteSystemInputs) : Promise<TDeleteSystemOutputs> {
    return postFormJson<TDeleteSystemOutputs>(deleteSystemUrl, inp, getAPIRequestConfig())
}

export interface TGetSystemInputs {
    sysId: number
    orgId: number
}

export interface TGetSystemOutputs {
    data: {
        System: System
    }
}

export function getSystem(inp : TGetSystemInputs) : Promise<TGetSystemOutputs> {
    return axios.get(getSystemUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
