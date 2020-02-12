import axios from 'axios'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import { DatabaseType, Database, DatabaseConnection } from '../databases'
import { postFormJson } from '../http'
import {
    newDatabaseUrl,
    allDatabaseUrl,
    typesDatabaseUrl,
    editDatabaseUrl,
    deleteDatabaseUrl,
    getDatabaseUrl,
    newDbConnUrl,
    deleteDbConnUrl,
    linkSystemsToDbUrl
} from '../url'
import { System } from '../systems'
import { FullDeployment } from '../deployments'

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
    deploymentType? : number
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
        Connection: DatabaseConnection
        RelevantSystemIds: number[]
        AllSystems: System[]
        Deployment: FullDeployment
    }
}

export function getDatabase(inp : TGetDatabaseInputs) : Promise<TGetDatabaseOutputs> {
    return axios.get(getDatabaseUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TNewDbConnInputs {
    dbId: number
    orgId: number
    host: string
    port: number
    dbName: string
    parameters: Record<string, string>
    username: string
    password: string
}

export interface TNewDbConnOutputs {
    data: DatabaseConnection
}

export function newDatabaseConnection(inp : TNewDbConnInputs) : Promise<TNewDbConnOutputs> {
    return postFormJson<TNewDbConnOutputs>(newDbConnUrl, inp, getAPIRequestConfig())
}

export interface TDeleteDbConnInputs {
    connId: number
    dbId: number
    orgId: number
}

export interface TDeleteDbConnOutputs {
}

export function deleteDatabaseConnection(inp : TDeleteDbConnInputs) : Promise<TDeleteDbConnOutputs> {
    return postFormJson<TDeleteDbConnOutputs>(deleteDbConnUrl, inp, getAPIRequestConfig())
}

export interface TLinkSystemInputs {
    dbId: number
    orgId: number
    sysIds: number[]
}

export interface TLinkSystemOutputs {
}

export function linkSystemsToDatabase(inp : TLinkSystemInputs) : Promise<TLinkSystemOutputs> {
    return postFormJson<TLinkSystemOutputs>(linkSystemsToDbUrl, inp, getAPIRequestConfig())
}
