import { NumericFilterData, NullNumericFilterData } from './filters'

export interface FullRiskData {
    Risk: ProcessFlowRisk
    Flows: ProcessFlowBasicData[]
    Controls: ProcessFlowControl[]
}

export interface RiskFilterData {
    NumControls: NumericFilterData
}
export let NullRiskFilterData : RiskFilterData = {
    NumControls: JSON.parse(JSON.stringify(NullNumericFilterData))
}
