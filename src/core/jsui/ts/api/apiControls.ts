import axios from 'axios'
import * as qs from 'query-string'
import { getControlTypesUrl,
         newControlUrl,
         addControlUrl,
         editControlUrl,
         deleteControlUrl,
         allControlAPIUrl,
         createSingleControlAPIUrl } from '../url'
import { postFormUrlEncoded } from '../http'
import { FullControlData } from '../controls'
import { getAPIRequestConfig } from './apiUtility'

export interface TGetControlTypesInput {
}

export interface TGetControlTypesOutput {
    data: ProcessFlowControlType[]
}

export function getControlTypes(inp : TGetControlTypesInput) : Promise<TGetControlTypesOutput> {
    return axios.get(getControlTypesUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TNewControlInput {
    name : string 
    description: string
    controlType : number
    frequencyType : number
    frequencyInterval : number
    ownerId : number
    nodeId:  number
    riskId : number
    orgName : string
}

export interface TNewControlOutput {
    data: ProcessFlowControl
}

export function newControl(inp: TNewControlInput) : Promise<TNewControlOutput> {
    return postFormUrlEncoded<TNewControlOutput>(newControlUrl, inp, getAPIRequestConfig())
}

export interface TDeleteControlInput {
    riskIds: number[]
    controlIds: number[]
    global: boolean
}

export interface TDeleteControlOutput {
}

export function deleteControls(inp : TDeleteControlInput): Promise<TDeleteControlOutput> {
    return postFormUrlEncoded<TDeleteControlOutput>(deleteControlUrl, inp, getAPIRequestConfig())
}


export interface TExistingControlInput {
    nodeId: number
    riskId: number
    controlIds: number[]
}

export interface TExistingControlOutput {
}


export function addExistingControls(inp : TExistingControlInput): Promise<TExistingControlOutput> {
    return postFormUrlEncoded<TExistingControlOutput>(addControlUrl, inp, getAPIRequestConfig())
}

export interface TEditControlInput extends TNewControlInput {
    controlId: number
}

export interface TEditControlOutput {
    data: ProcessFlowControl
}

export function editControl(inp: TEditControlInput) : Promise<TEditControlOutput> {
    return postFormUrlEncoded<TEditControlOutput>(editControlUrl, inp, getAPIRequestConfig())
}

export interface TAllControlInput {
    orgName: string
}

export interface TAllControlOutput {
    data: ProcessFlowControl[]
}

export function getAllControls(inp : TAllControlInput) : Promise<TAllControlOutput> {
    return axios.get(allControlAPIUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TSingleControlInput {
    controlId: number
}

export interface TSingleControlOutput {
    data: FullControlData
}

export function getSingleControl(inp : TSingleControlInput) : Promise<TSingleControlOutput> {
    return axios.get(createSingleControlAPIUrl(inp.controlId) + '?' + qs.stringify(inp), getAPIRequestConfig())
}
