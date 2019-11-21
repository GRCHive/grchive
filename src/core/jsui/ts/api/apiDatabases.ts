import axios from 'axios'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import { DatabaseType, Database} from '../databases'
import { postFormJson } from '../http'
import {
    newDatabaseUrl,
    allDatabaseUrl,
    typesDatabaseUrl,
} from '../url'

export interface TDbTypeOutputs {
    data: DatabaseType[]
}

export function getAllDatabaseTypes(): Promise<TDbTypeOutputs> {
    return axios.get(typesDatabaseUrl, getAPIRequestConfig())
}

export interface TNewDatabaseInputs {
    name: string
    orgId: number
    typeId: number
    otherType: string
    version: string
}

export interface TNewDatabaseOutputs {
    data: Database
}

export function newDatabase(inp : TNewDatabaseInputs) : Promise<TNewDatabaseOutputs> {
    return postFormJson<TNewDatabaseOutputs>(newDatabaseUrl, inp, getAPIRequestConfig())
}

export interface TAllDatabaseInputs {
    orgId: number
}

export interface TAllDatabaseOutputs {
    data: Database[]
}

export function allDatabase(inp : TAllDatabaseInputs) : Promise<TAllDatabaseOutputs> {
    return axios.get(allDatabaseUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
