import axios from 'axios'
import * as qs from 'query-string'
import { getControlTypesUrl,
         newControlUrl } from '../url'
import { postFormUrlEncoded } from '../http'

export function getControlTypes(inp : TGetControlTypesInput) : Promise<TGetControlTypesOutput> {
    return axios.get(getControlTypesUrl + '?' + qs.stringify(inp))
}

export function newControl(inp: TNewControlInput) : Promise<TNewControlOutput> {
    return postFormUrlEncoded<TNewControlOutput>(newControlUrl, inp)
}
