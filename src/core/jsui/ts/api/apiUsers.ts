import axios from 'axios'
import * as qs from 'query-string'
import { createGetAllOrgUsersAPIUrl } from '../url'

export function getAllOrgUsers(inp : TGetAllOrgUsersInput) : Promise<TGetAllOrgUsersOutput> {
    return axios.get(createGetAllOrgUsersAPIUrl(inp.org) + '?' + qs.stringify(inp))
}
