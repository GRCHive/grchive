import axios from 'axios'
import * as qs from 'query-string'
import { getControlTypesUrl,
         newControlUrl,
         addControlUrl,
         editControlUrl,
         deleteControlUrl,
         allControlAPIUrl,
         createSingleControlAPIUrl,
         linkCatControlUrl,
         unlinkCatControlUrl } from '../url'
import { postFormUrlEncoded, postFormJson } from '../http'
import { FullControlData, ControlFilterData } from '../controls'
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
    frequencyOther: string
    ownerId : number
    nodeId:  number
    riskId : number
    orgName : string
    manual: boolean
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
    filter: ControlFilterData
}

export interface TAllControlOutput {
    data: ProcessFlowControl[]
}

export function getAllControls(inp : TAllControlInput) : Promise<TAllControlOutput> {
    let passData : any = {
        orgName: inp.orgName,
        filter: JSON.stringify(inp.filter),
    }
    return axios.get(allControlAPIUrl + '?' + qs.stringify(passData), getAPIRequestConfig())
}

export interface TSingleControlInput {
    controlId: number
    orgId: number
}

export interface TSingleControlOutput {
    data: FullControlData
}

export function getSingleControl(inp : TSingleControlInput) : Promise<TSingleControlOutput> {
    return axios.get(createSingleControlAPIUrl(inp.controlId) + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TLinkToDocCatInput {
    controlId: number
    orgId: number
    catId: number
    isInput: boolean
}

export function linkControlToDocumentCategory(inp : TLinkToDocCatInput) : Promise<void> {
    return postFormJson<void>(linkCatControlUrl, inp, getAPIRequestConfig())
}

export interface TUnlinkFromDocCatInput {
    controlId: number
    orgId: number
    catId: number
    isInput: boolean
}

export function unlinkControlFromDocumentCategory(inp : TUnlinkFromDocCatInput) : Promise<void> {
    return postFormJson<void>(unlinkCatControlUrl, inp, getAPIRequestConfig())
}
