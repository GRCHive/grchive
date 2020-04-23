import { Days } from './time'
import moment from 'moment-timezone'

export enum CronFrequency {
    Daily,
    Weekly,
    Monthly
}

export enum CronWeekdayHash {
    First,
    Second,
    Third,
    Fourth,
    Last
}

export let CronWeekdayHashItems = Object.keys(CronWeekdayHash)
    .filter((key : any) => !isNaN(Number(CronWeekdayHash[key])))
    .map((key : any) => ({
        text: key,
        value: CronWeekdayHash[key]
    }))

export interface DailyCron {
    Times : Date[]
}

export interface WeeklyCron {
    Days : Days[]
    Time: Date
}

export interface MonthlyCron {
    UseDate: boolean
    // Run every month on these dates
    Dates: number[]

    // Run every month on the [Nth] [Day] of the month.
    Nth: CronWeekdayHash
    Day: Days

    Time: Date
}

export interface ScheduledEvent {
    Repeat: boolean
    OneTimeDate: Date | null
    Frequency: CronFrequency
    Name : string
    Description: string
    Daily: DailyCron | null
    Weekly: WeeklyCron | null
    Monthly: MonthlyCron | null
    Timezone: string
}

export function createEmptyScheduledEvent() : ScheduledEvent {
    return {
        Repeat: false,
        OneTimeDate: new Date(),
        Name: "",
        Description: "",
        Frequency: CronFrequency.Daily,
        Daily: createEmptyDailyCron(),
        Weekly: null,
        Monthly: null,
        Timezone: moment.tz.guess(),
    }
}

export function cleanScheduledEventFromJson(e : ScheduledEvent) {
    if (!!e.OneTimeDate) {
        e.OneTimeDate = new Date(e.OneTimeDate)
    }

    if (!!e.Daily) {
        cleanDailyCronFromJson(e.Daily)
    }

    if (!!e.Weekly) {
        cleanWeeklyCronFromJson(e.Weekly)
    }

    if (!!e.Monthly) {
        cleanMonthlyCronFromJson(e.Monthly)
    }
}

export function createEmptyDailyCron() : DailyCron {
    return {
        Times: [],
    }
}

export function cleanDailyCronFromJson(e : DailyCron) {
    e.Times = e.Times.map((ele : Date) => new Date(ele))
}

export function createEmptyWeeklyCron() : WeeklyCron {
    return {
        Days: [],
        Time: new Date()
    }
}

export function cleanWeeklyCronFromJson(e : WeeklyCron) {
    e.Time = new Date(e.Time)
}

export function createEmptyMonthlyCron() : MonthlyCron {
    return {
        UseDate: true,
        Dates: [],
        Nth: CronWeekdayHash.First,
        Day: Days.Monday,
        Time: new Date(),
    }
}

export function cleanMonthlyCronFromJson(e : MonthlyCron) {
    e.Time = new Date(e.Time)
}
