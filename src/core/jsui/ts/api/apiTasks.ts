import axios from 'axios'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import {
    allScheduleUrl
} from '../url'
import {
    ScheduledTaskMetadata,
    cleanScheduledTaskMetadataFromJson
} from '../tasks'
import {
    TimeRange,
    cleanTimeRangeFromJson,
} from '../time'
import moment from 'moment-timezone'

export interface TAllScheduledTasksInput {
    orgId: number
    range: TimeRange
}

export interface TAllScheduledTasksOutput {
    data: {
        Tasks: ScheduledTaskMetadata[],
        Times: Record<number, TimeRange[]>
    }
}

export function allScheduledTasks(inp : TAllScheduledTasksInput) : Promise<TAllScheduledTasksOutput> {
    let params = {
        ...inp,
        range: JSON.stringify(inp.range),
        timezone: moment.tz.guess(),
    }

    return axios.get(allScheduleUrl + '?' + qs.stringify(params), getAPIRequestConfig()).then((resp : TAllScheduledTasksOutput) => {
        for (let t of resp.data.Tasks) {
            resp.data.Times[t.Id].forEach(cleanTimeRangeFromJson)
        }

        resp.data.Tasks.forEach(cleanScheduledTaskMetadataFromJson)
        return resp
    })
}
