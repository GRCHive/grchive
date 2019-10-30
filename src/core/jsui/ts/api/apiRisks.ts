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

export function newRisk(inp : TNewRiskInput) : Promise<TNewRiskOutput> {
    return postFormUrlEncoded<TNewRiskOutput>(newRiskAPIUrl, inp)
}

export function deleteRisk(inp : TDeleteRiskInput) : Promise<TDeleteRiskOutput> {
    return postFormUrlEncoded<TDeleteRiskOutput>(deleteRiskAPIUrl, inp)
}

export function addExistingRisk(inp : TAddExistingRiskInput) : Promise<TAddExistingRiskOutput> {
    return postFormUrlEncoded<TAddExistingRiskOutput>(addExistingRiskAPIUrl, inp)
}

export interface TEditRiskInput {
    csrf: string
    name: string
    description : string
    riskId: number
}

export interface TEditRiskOutput {
    data: ProcessFlowRisk
}

export function editRisk(inp : TEditRiskInput) : Promise<TEditRiskOutput> {
    return postFormUrlEncoded<TEditRiskOutput>(editRiskAPIUrl, inp)
}

export interface TAllRiskInput {
    csrf: string
}

export interface TAllRiskOutput {
    data: ProcessFlowRisk[]
}

export function getAllRisks(inp : TAllRiskInput) : Promise<TAllRiskOutput> {
    return axios.get(allRiskAPIUrl + '?' + qs.stringify(inp))
}

export interface TSingleRiskInput {
    csrf: string
    riskId: number
}

export interface TSingleRiskOutput {
    data: FullRiskData
}

export function getSingleRisk(inp : TSingleRiskInput) : Promise<TSingleRiskOutput> {
    return axios.get(createSingleRiskAPIUrl(inp.riskId) + '?' + qs.stringify(inp))
}
