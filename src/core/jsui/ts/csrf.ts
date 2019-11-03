import { getCookie } from './cookie'

export function getCurrentCSRF() : string {
    return getCookie('client-csrf')!
}
