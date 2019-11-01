import axios from 'axios';
import * as qs from 'query-string';

export function postFormUrlEncoded<T=void>(url: string, data: Object) : Promise<T> {
    const config = {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
    }
    return axios.post(url, qs.stringify(data), config)
}

export function postFormMultipart<T=void>(url : string, data : FormData) : Promise<T> {
    const config = {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    }
    return axios.post(url, data, config)
}
