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

