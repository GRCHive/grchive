import axios from 'axios'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import { postFormJson } from '../http'
import { VendorProduct } from '../vendors'
import { 
    newVendorProductUrl,
    allVendorProductsUrl,
    getVendorProductUrl,
    updateVendorProductUrl,
    deleteVendorProductUrl,
    newVendorProductSocLinkUrl,
    deleteVendorProductSocLinkUrl,
} from '../url'
import { ControlDocumentationFile, ControlDocumentationFileHandle, cleanJsonControlDocumentationFile } from '../controls'

export interface TNewVendorProductInput {
    orgId: number
    vendorId: number
    name: string
    description: string
    url: string
}

export interface TNewVendorProductOutput {
    data: VendorProduct
}

export function newVendorProduct(inp : TNewVendorProductInput) : Promise<TNewVendorProductOutput> {
    return postFormJson<TNewVendorProductOutput>(newVendorProductUrl, inp, getAPIRequestConfig())
}

export interface TAllVendorProductInput {
    orgId: number
    vendorId: number
}

export interface TAllVendorProductOutput {
    data: VendorProduct[]
}

export function allVendorProducts(inp : TAllVendorProductInput) : Promise<TAllVendorProductOutput> {
    return axios.get(allVendorProductsUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TGetVendorProductInput {
    productId: number
    vendorId: number
    orgId: number
}

export interface TGetVendorProductOutput {
    data: {
        Product: VendorProduct
        SocFiles: ControlDocumentationFile[]
    }
}

export function getVendorProduct(inp : TGetVendorProductInput) : Promise<TGetVendorProductOutput> {
    return axios.get(getVendorProductUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TGetVendorProductOutput) => {
        resp.data.SocFiles.forEach(cleanJsonControlDocumentationFile)
        return resp
    })
}

export interface TDeleteVendorProductInput {
    productId: number
    vendorId: number
    orgId: number
}

export function deleteVendorProduct(inp : TDeleteVendorProductInput) : Promise<void> {
    return postFormJson<void>(deleteVendorProductUrl, inp, getAPIRequestConfig())
}

export interface TUpdateVendorProductInput extends TNewVendorProductInput {
    productId: number
}

export function updateVendorProduct(inp : TUpdateVendorProductInput) : Promise<TNewVendorProductOutput> {
    return postFormJson<TNewVendorProductOutput>(updateVendorProductUrl, inp, getAPIRequestConfig())
}

export interface TLinkProductSocInput {
    productId: number
    vendorId: number
    orgId: number
    socFiles: ControlDocumentationFileHandle[]
}

export function linkVendorProductSocFiles(inp : TLinkProductSocInput) : Promise<void> {
    return postFormJson<void>(newVendorProductSocLinkUrl, inp, getAPIRequestConfig())
}

export interface TUnlinkProductSocInput {
    productId: number
    vendorId: number
    orgId: number
    socFiles: ControlDocumentationFileHandle[]
}

export function unlinkVendorProductSocFiles(inp : TUnlinkProductSocInput) : Promise<void> {
    return postFormJson<void>(deleteVendorProductSocLinkUrl, inp, getAPIRequestConfig())
}
