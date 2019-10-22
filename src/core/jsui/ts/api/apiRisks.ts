import { postFormUrlEncoded } from '../http'
import { newRiskAPIUrl, deleteRiskAPIUrl  } from '../url'

export function newRisk(inp : TNewRiskInput) : Promise<TNewRiskOutput> {
    return postFormUrlEncoded<TNewRiskOutput>(newRiskAPIUrl, inp)
}

export function deleteRisk(inp : TDeleteRiskInput) : Promise<TDeleteRiskOutput> {
    return postFormUrlEncoded<TDeleteRiskOutput>(deleteRiskAPIUrl, inp)
}
