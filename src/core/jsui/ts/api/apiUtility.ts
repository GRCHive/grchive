import { AxiosRequestConfig } from 'axios'
import { getCookie } from '../cookie'

export function getAPIRequestConfig() : AxiosRequestConfig {
    return {
        headers: {
            'ApiKey': getCookie('client-api-key')!
        }
    }
}
