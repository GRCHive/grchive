import axios from 'axios';
import * as qs from 'query-string';

export function postFormUrlEncoded(url: string, data: Object) : Promise<void> {
    const config = {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
    }
    return axios.post(url, qs.stringify(data), config)
}
