import axios from 'axios'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import { DatabaseType, Database} from '../databases'
import { postFormJson } from '../http'
import {
    newDatabaseUrl,
    allDatabaseUrl,
    typesDatabaseUrl,
    editDatabaseUrl,
    deleteDatabaseUrl,
    getDatabaseUrl,
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

export interface TEditDatabaseInputs extends TNewDatabaseInputs {
    dbId: number
}

export interface TEditDatabaseOutputs {
    data: Database
}

export function editDatabase(inp : TEditDatabaseInputs) : Promise<TEditDatabaseOutputs> {
    return postFormJson<TEditDatabaseOutputs>(editDatabaseUrl, inp, getAPIRequestConfig())
}

export interface TDeleteDatabaseInputs {
    dbId: number
    orgId: number
}

export interface TDeleteDatabaseOutputs {
}

export function deleteDatabase(inp : TDeleteDatabaseInputs) : Promise<TDeleteDatabaseOutputs> {
    return postFormJson<TDeleteDatabaseOutputs>(deleteDatabaseUrl, inp, getAPIRequestConfig())
}

export interface TGetDatabaseInputs {
    dbId: number
    orgId: number
}

export interface TGetDatabaseOutputs {
    data: {
        Database: Database
    }
}

export function getDatabase(inp : TGetDatabaseInputs) : Promise<TGetDatabaseOutputs> {
    return axios.get(getDatabaseUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
