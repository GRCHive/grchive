import { postFormUrlEncoded } from '../http'
import { newRiskAPIUrl } from '../url'

export function newRisk(inp : TNewRiskInput) : Promise<TNewRiskOutput> {
    return postFormUrlEncoded<TNewRiskOutput>(newRiskAPIUrl, inp)
}
