import axios from 'axios'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import { postFormJson } from '../http'
import {
    newSystemUrl,
    allSystemsUrl
} from '../url'
import { System } from '../systems'

export interface TNewSystemInputs {
    orgId: number
    name: string
    purpose: string
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
