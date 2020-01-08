import { postFormJson } from '../http'
import { getAPIRequestConfig } from './apiUtility'
import {
    newDeploymentUrl,
    updateDeploymentUrl
} from '../url'
import { FullDeployment, deepCopyFullDeployment, createStrippedDeployment } from '../deployments'

export interface TNewDeploymentInput {
    orgId: number
    systemId: number | null
    dbId: number | null
}

export interface TNewDeploymentOutput {
    data : FullDeployment
}

export function newDeployment(inp : TNewDeploymentInput) : Promise<TNewDeploymentOutput> {
    return postFormJson<TNewDeploymentOutput>(newDeploymentUrl, inp, getAPIRequestConfig())
}

export interface TUpdateDeploymentInput {
    deployment: FullDeployment
}

export interface TUpdateDeploymentOutput {
    data : FullDeployment
}

export function updateDeployment(inp : TUpdateDeploymentInput) : Promise<TUpdateDeploymentOutput> {
    return postFormJson<TUpdateDeploymentOutput>(
        updateDeploymentUrl,
        {
            deployment: createStrippedDeployment(inp.deployment)
        },
        getAPIRequestConfig()).then(
    (resp : TUpdateDeploymentOutput) => {
        // primarily used to convert the dates from strings to Date objects
        resp.data = deepCopyFullDeployment(resp.data)
        return resp
    })
}
