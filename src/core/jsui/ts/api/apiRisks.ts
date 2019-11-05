import axios from 'axios'
import * as qs from 'query-string'
import { postFormUrlEncoded } from '../http'
import { newRiskAPIUrl,
         deleteRiskAPIUrl,
         addExistingRiskAPIUrl,
         editRiskAPIUrl,
         allRiskAPIUrl,
         createSingleRiskAPIUrl } from '../url'
import { FullRiskData } from '../risks'
import { getAPIRequestConfig } from './apiUtility'

export interface TNewRiskInput {
    name: string
    description : string
    nodeId: number
    orgName: string
}

export interface TNewRiskOutput {
    data: ProcessFlowRisk
}

export function newRisk(inp : TNewRiskInput) : Promise<TNewRiskOutput> {
    return postFormUrlEncoded<TNewRiskOutput>(newRiskAPIUrl, inp, getAPIRequestConfig())
}

export interface TDeleteRiskInput {
    nodeId: number
    riskIds: number[]
    global: boolean
}

export interface TDeleteRiskOutput {
}

export function deleteRisk(inp : TDeleteRiskInput) : Promise<TDeleteRiskOutput> {
    return postFormUrlEncoded<TDeleteRiskOutput>(deleteRiskAPIUrl, inp, getAPIRequestConfig())
}

export interface TAddExistingRiskInput {
    nodeId: number
    riskIds: number[]
}

export interface TAddExistingRiskOutput {
}

export function addExistingRisk(inp : TAddExistingRiskInput) : Promise<TAddExistingRiskOutput> {
    return postFormUrlEncoded<TAddExistingRiskOutput>(addExistingRiskAPIUrl, inp, getAPIRequestConfig())
}

export interface TEditRiskInput {
    name: string
    description : string
    riskId: number
    orgName: string
}

export interface TEditRiskOutput {
    data: ProcessFlowRisk
}

export function editRisk(inp : TEditRiskInput) : Promise<TEditRiskOutput> {
    return postFormUrlEncoded<TEditRiskOutput>(editRiskAPIUrl, inp, getAPIRequestConfig())
}

export interface TAllRiskInput {
    orgName: string
}

export interface TAllRiskOutput {
    data: ProcessFlowRisk[]
}

export function getAllRisks(inp : TAllRiskInput) : Promise<TAllRiskOutput> {
    return axios.get(allRiskAPIUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TSingleRiskInput {
    riskId: number
}

export interface TSingleRiskOutput {
    data: FullRiskData
}

export function getSingleRisk(inp : TSingleRiskInput) : Promise<TSingleRiskOutput> {
    return axios.get(createSingleRiskAPIUrl(inp.riskId) + '?' + qs.stringify(inp), getAPIRequestConfig())
}
