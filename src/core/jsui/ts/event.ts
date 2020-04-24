import { Days, DaysKey } from './time'
import moment from 'moment-timezone'
import { RRule, RRuleSet, Weekday, WeekdayStr, rrulestr } from 'rrule'

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

export function createScheduledEventFromRRule(rrule : string | null | undefined) : (ScheduledEvent | null) {
    if (!rrule) {
        return null
    }
    
    let rule = rrulestr(rrule)
    let schedule : ScheduledEvent = createEmptyScheduledEvent()

    schedule.Repeat = true
    schedule.OneTimeDate = null

    let convertRRuleDayToHash = (d : Weekday | WeekdayStr | number) : CronWeekdayHash => {
        if (typeof d === 'string') {
            let newD = Weekday.fromStr(<WeekdayStr>d)
            return convertRRuleDayToHash(newD)
        } else if (typeof d === 'number') {
            // ????
            return CronWeekdayHash.First
        } else {
            switch (d.n!) {
            case 1:
                return CronWeekdayHash.First
            case 2:
                return CronWeekdayHash.Second
            case 3:
                return CronWeekdayHash.Third
            case 4:
                return CronWeekdayHash.Fourth
            case 5:
                return CronWeekdayHash.Last
            case -1:
                return CronWeekdayHash.Last
            }
        }

        return CronWeekdayHash.First

    }

    let convertRRuleDay = (d : Weekday | WeekdayStr | number) : Days => {
        if (typeof d === 'string') {
            let newD = Weekday.fromStr(<WeekdayStr>d)
            return convertRRuleDay(newD)
        } else if (typeof d === 'number') {
            let newD : number = d - 1
            if (newD == -1) {
                newD = 6
            }
            return Days[<DaysKey>Days[newD]]
        } else {
            switch (d) {
            case RRule.MO:
                return Days.Monday
            case RRule.TU:
                return Days.Tuesday
            case RRule.WE:
                return Days.Wednesday
            case RRule.TH:
                return Days.Thursday
            case RRule.FR:
                return Days.Friday
            case RRule.SA:
                return Days.Saturday
            case RRule.SU:
                return Days.Sunday
            }
        }

        return Days.Sunday
    }

    let parseRrule = (r : RRule) => {
        if (r.origOptions.freq == RRule.DAILY) {
            schedule.Frequency = CronFrequency.Daily
            if (!schedule.Daily) {
                schedule.Daily = createEmptyDailyCron()
                schedule.Weekly = null
                schedule.Monthly = null
            }

            schedule.Daily.Times.push(r.origOptions.dtstart!)
        } else if (r.origOptions.freq == RRule.WEEKLY) {
            schedule.Frequency = CronFrequency.Weekly
            if (!schedule.Weekly) {
                schedule.Daily = null
                schedule.Weekly = createEmptyWeeklyCron()
                schedule.Monthly = null
            }

            if (Array.isArray(r.origOptions.byweekday!)) {
                schedule.Weekly.Days = r.origOptions.byweekday!.map((ele : any) => convertRRuleDay(ele))
            } else {
                schedule.Weekly.Days = [convertRRuleDay(r.origOptions.byweekday!)]
            }
            schedule.Weekly.Time = r.origOptions.dtstart!
        } else if (r.origOptions.freq == RRule.MONTHLY) {
            schedule.Frequency = CronFrequency.Monthly
            if (!schedule.Monthly) {
                schedule.Daily = null
                schedule.Weekly = null
                schedule.Monthly = createEmptyMonthlyCron()
            }

            schedule.Monthly.UseDate = (r.origOptions.bymonthday !== null)
            if (schedule.Monthly.UseDate) {
                if (Array.isArray(r.origOptions.bymonthday)) {
                    schedule.Monthly.Dates = r.origOptions.bymonthday
                } else {
                    schedule.Monthly.Dates = [r.origOptions.bymonthday!]
                }
            } else {
                if (Array.isArray(r.origOptions.byweekday!)) {
                    schedule.Monthly.Day = convertRRuleDay(r.origOptions.byweekday![0])
                    schedule.Monthly.Nth = convertRRuleDayToHash(r.origOptions.byweekday![0])
                } else {
                    schedule.Monthly.Day = convertRRuleDay(r.origOptions.byweekday!)
                    schedule.Monthly.Nth = convertRRuleDayToHash(r.origOptions.byweekday!)
                }
            }
            schedule.Monthly.Time = r.origOptions.dtstart!
        }
    }

    if ('rrules' in rule) {
        for (let r of rule.rrules()) {
            parseRrule(r)
        }
    } else {
        parseRrule(rule)
    }

    return schedule
}
