import axios from 'axios';
import {AxiosRequestConfig } from 'axios';
import * as qs from 'query-string';

export function postFormJson<T=void>(url: string, data: Object, config : AxiosRequestConfig) : Promise<T> {
    if (!config.headers) {
        config.headers = {}
    }
    config.headers['Content-Type'] = 'application/json'
    return axios.post(url, JSON.stringify(data), config)
}

export function putFormJson<T=void>(url: string, data: Object, config : AxiosRequestConfig) : Promise<T> {
    if (!config.headers) {
        config.headers = {}
    }
    config.headers['Content-Type'] = 'application/json'
    return axios.put(url, JSON.stringify(data), config)
}


export function deleteFormJson<T=void>(url: string, data: Object, config : AxiosRequestConfig) : Promise<T> {
    if (!config.headers) {
        config.headers = {}
    }
    config.headers['Content-Type'] = 'application/json'
    config.data = data
    return axios.delete(url, config)
}

export function postFormUrlEncoded<T=void>(url: string, data: Object, config : AxiosRequestConfig) : Promise<T> {
    if (!config.headers) {
        config.headers = {}
    }
    config.headers['Content-Type'] = 'application/x-www-form-urlencoded'
    return axios.post(url, qs.stringify(data), config)
}

export function postFormMultipart<T=void>(url : string, data : FormData, config : AxiosRequestConfig) : Promise<T> {
    if (!config.headers) {
        config.headers = {}
    }
    config.headers['Content-Type'] = 'multipart/form-data'
    return axios.post(url, data, config)
}

export function putFormMultipart<T=void>(url: string, data: FormData, config : AxiosRequestConfig) : Promise<T> {
    if (!config.headers) {
        config.headers = {}
    }
    config.headers['Content-Type'] = 'multipart/form-data'
    return axios.put(url, data, config)
}
