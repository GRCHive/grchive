import axios from 'axios'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import { postFormJson } from '../http'
import { Vendor } from '../vendors'
import { 
    newVendorUrl,
    allVendorsUrl,
    getVendorUrl,
    updateVendorUrl,
    deleteVendorUrl,
} from '../url'

export interface TNewVendorInput {
    orgId: number
    name: string
    description: string
    url: string
}

export interface TNewVendorOutput {
    data: Vendor
}

export function newVendor(inp : TNewVendorInput) : Promise<TNewVendorOutput> {
    return postFormJson<TNewVendorOutput>(newVendorUrl, inp, getAPIRequestConfig())
}

export interface TAllVendorInput {
    orgId: number
}

export interface TAllVendorOutput {
    data: Vendor[]
}

export function allVendors(inp : TAllVendorInput) : Promise<TAllVendorOutput> {
    return axios.get(allVendorsUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TGetVendorInput {
    vendorId: number
    orgId: number
}

export interface TGetVendorOutput {
    data: Vendor
}

export function getVendor(inp : TGetVendorInput) : Promise<TGetVendorOutput> {
    return axios.get(getVendorUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TDeleteVendorInput {
    vendorId: number
    orgId: number
}

export function deleteVendor(inp : TDeleteVendorInput) : Promise<void> {
    return postFormJson<void>(deleteVendorUrl, inp, getAPIRequestConfig())
}

export interface TUpdateVendorInput extends TNewVendorInput {
    vendorId: number
}

export function updateVendor(inp : TUpdateVendorInput) : Promise<TNewVendorOutput> {
    return postFormJson<TNewVendorOutput>(updateVendorUrl, inp, getAPIRequestConfig())
}
