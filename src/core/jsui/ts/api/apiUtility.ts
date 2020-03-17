import { AxiosRequestConfig } from 'axios'
import { getCookie } from '../cookie'
import axios from 'axios'

export function getJsonApiKey() : any {
    return JSON.parse(atob(getCookie('client-api-key')!))
}

export function getApiKey() : string {
    return getJsonApiKey()['Key']
}

export function getApiKeyExpirationTime() : Date {
    return new Date(getJsonApiKey()['Expiration'])
}

export function getAPIRequestConfig() : AxiosRequestConfig {
    return {
        headers: {
            'ApiKey': getApiKey()
        }
    }
}

export function startTemporaryApiKeyRefresh() {
    let expirationTime : Date = getApiKeyExpirationTime()

    // Refresh the API key ~10 min beforehand.
    let timeoutMs : number = Math.max(expirationTime.getTime() - new Date().getTime() - 10 * 60 * 1000, 1)

    setTimeout(() => {
        axios.get('/dashboard').then(startTemporaryApiKeyRefresh)
    }, timeoutMs)
}
