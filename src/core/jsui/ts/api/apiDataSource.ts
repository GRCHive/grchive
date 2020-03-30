import axios from 'axios'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import {
    DataSourceOption  
} from '../clientData'
import { 
    allDataSourceUrl,
    getDataSourceUrl,
} from '../url'
import { DataSourceLink } from '../clientData'
import { ResourceHandle } from '../resourceUtils'

export interface TAllDataSourceOutput {
    data: DataSourceOption[]
}

export function allSupportedDataSources() : Promise<TAllDataSourceOutput> {
    return axios.get(allDataSourceUrl, getAPIRequestConfig())
}

export interface TGetDataSourceInput {
    source : DataSourceLink
}

export interface TGetDataSourceOutput {
    data: ResourceHandle
}

export function getDataSource(inp : TGetDataSourceInput) : Promise<TGetDataSourceOutput> {
    let params = {
        source : JSON.stringify(inp.source)
    }
    return axios.get(getDataSourceUrl + '?' + qs.stringify(params), getAPIRequestConfig())
}
