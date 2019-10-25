import axios from 'axios'
import * as qs from 'query-string'
import { getControlTypesUrl,
         newControlUrl,
         addControlUrl,
         editControlUrl,
         deleteControlUrl } from '../url'
import { postFormUrlEncoded } from '../http'

export function getControlTypes(inp : TGetControlTypesInput) : Promise<TGetControlTypesOutput> {
    return axios.get(getControlTypesUrl + '?' + qs.stringify(inp))
}

export function newControl(inp: TNewControlInput) : Promise<TNewControlOutput> {
    return postFormUrlEncoded<TNewControlOutput>(newControlUrl, inp)
}

export interface TDeleteControlInput {
    csrf: string
    nodeId: number
    riskIds: number[]
    controlIds: number[]
    global: boolean
}

export interface TDeleteControlOutput {
}

export function deleteControls(inp : TDeleteControlInput): Promise<TDeleteControlOutput> {
    return postFormUrlEncoded<TDeleteControlOutput>(deleteControlUrl, inp)
}


export interface TExistingControlInput {
    csrf: string
    nodeId: number
    riskId: number
    controlIds: number[]
}

export interface TExistingControlOutput {
}


export function addExistingControls(inp : TExistingControlInput): Promise<TExistingControlOutput> {
    return postFormUrlEncoded<TExistingControlOutput>(addControlUrl, inp)
}

export interface TEditControlInput extends TNewControlInput {
    controlId: number
}

export interface TEditControlOutput {
    data: ProcessFlowControl
}

export function editControl(inp: TEditControlInput) : Promise<TEditControlOutput> {
    return postFormUrlEncoded<TEditControlOutput>(editControlUrl, inp)
}
