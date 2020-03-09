export enum ComparisonOperators {
    Disabled = 0,
    Equal = 1,
    NotEqual = 2,
    Greater = 3,
    GreaterEqual = 4,
    Less = 5,
    LessEqual = 6
}

export let comparisonOperatorsSelectItems : any[] = [
    {
        text: "Disabled",
        value: ComparisonOperators.Disabled,
    },
    {
        text: "Equal To",
        value: ComparisonOperators.Equal,
    },
    {
        text: "Not Equal To",
        value: ComparisonOperators.NotEqual,
    },
    {
        text: "Greater Than (>)",
        value: ComparisonOperators.Greater,
    },
    {
        text: "Greater Than or Equal To (>=)",
        value: ComparisonOperators.GreaterEqual,
    },
    {
        text: "Less Than (<)",
        value: ComparisonOperators.Less,
    },
    {
        text: "Less Than or Equal To (<=)",
        value: ComparisonOperators.LessEqual,
    },
]

export interface NumericFilterData {
    Op          : ComparisonOperators
    Target      : number
}
export let NullNumericFilterData : NumericFilterData = {
    Op: ComparisonOperators.Disabled,
    Target: 0
}

export enum StringComparisonOperators {
    Disabled = 0,
    Equal = 1,
    NotEqual = 2,
    Contains = 3,
    Excludes = 4
}

export let stringComparisonOperatorsSelectItems : any[] = [
    {
        text: "Disabled",
        value: StringComparisonOperators.Disabled,
    },
    {
        text: "Equal To",
        value: StringComparisonOperators.Equal,
    },
    {
        text: "Not Equal To",
        value: StringComparisonOperators.NotEqual,
    },
    {
        text: "Contains",
        value: StringComparisonOperators.Contains,
    },
    {
        text: "Excludes",
        value: StringComparisonOperators.Excludes,
    },
]

export interface StringFilterData {
    Op      : StringComparisonOperators
    Target  : string
}
export let NullStringFilterData : StringFilterData = {
    Op: StringComparisonOperators.Disabled,
    Target: ""
}

export interface TimeRangeFilterData {
    Enabled : boolean
    Start   : Date
    End     : Date
}

export let NullTimeRangeFilterDate : TimeRangeFilterData = {
    Enabled: true,
    Start: (() => {
        let d : Date = new Date()
        return new Date(d.getTime() - 1000*60*60*24)
    })(),
    End: new Date()
}

export function cleanTimeRangeFilterDataFromJson(d : TimeRangeFilterData) : TimeRangeFilterData {
    d.Start = new Date(d.Start)
    d.End = new Date(d.End)
    return d
}
