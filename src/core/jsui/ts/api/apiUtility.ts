import { AxiosRequestConfig } from 'axios'
import { getCookie } from '../cookie'

export function getApiKey() : string {
    return getCookie('client-api-key')!
}

export function getAPIRequestConfig() : AxiosRequestConfig {
    return {
        headers: {
            'ApiKey': getApiKey()
        }
    }
}
