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
    Start   : Date | null
    End     : Date | null
}

export let NullTimeRangeFilterDate : TimeRangeFilterData = {
    Enabled: true,
    Start: null,
    End: null,
}

export function copyTimeRangeFilterData(c : TimeRangeFilterData) : TimeRangeFilterData {
    let ret = JSON.parse(JSON.stringify(c))
    cleanTimeRangeFilterDataFromJson(ret)
    return ret
}

export function cleanTimeRangeFilterDataFromJson(d : TimeRangeFilterData) : TimeRangeFilterData {
    if (!!d.Start) {
        d.Start = new Date(d.Start)
    }

    if (!!d.End) {
        d.End = new Date(d.End)
    }
    return d
}

export interface UserFilterData {
    UserIds: (number | null)[]
}

export let NullUserFilterData : UserFilterData = {
    UserIds: [],
}

export function copyUserFilterData(c : UserFilterData) : UserFilterData {
    let ret = JSON.parse(JSON.stringify(c))
    return ret
}
