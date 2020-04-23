<template>
    <div>
        <v-calendar
            v-model="currentTimeStr"
            type="month"
            :events="scheduledVueitfyEvents"
            @change="refreshSchedule"
        >
        </v-calendar>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { TAllScheduledTasksOutput, allScheduledTasks } from '../../../ts/api/apiTasks'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import { ScheduledTaskMetadata } from '../../../ts/tasks'
import { TimeRange, vuetifyCalendarTimeFormat } from '../../../ts/time'

@Component
export default class ScheduleViewer extends Vue {
    currentTimeStr : string = new Date().toString()

    allTasks : ScheduledTaskMetadata[] = []
    taskTimes : Record<number, TimeRange[]> = Object()

    get scheduledVueitfyEvents() : any[] {
        return this.allTasks.map((ele : ScheduledTaskMetadata) => {
            return this.taskTimes[ele.Id].map((rn : TimeRange) => {
                return {
                    name : ele.Name,
                    start: vuetifyCalendarTimeFormat(rn.Start),
                    end: vuetifyCalendarTimeFormat(rn.End),
                    color: 'blue',
                }
            })
        }).flat()
    }

    refreshSchedule(props : any) {
        function createDate(d : any, hour : number, minute : number, seconds: number) : Date {
            return new Date(
                d.year,
                d.month-1,
                d.day,
                hour,
                minute,
                seconds
            )
        }

        allScheduledTasks({
            orgId: PageParamsStore.state.organization!.Id,
            range: {
                Start: createDate(props.start, 0, 0, 0),
                End: createDate(props.end, 23, 59, 59),
            },
        }).then((resp : TAllScheduledTasksOutput) => {
            this.allTasks = resp.data.Tasks
            this.taskTimes = resp.data.Times
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }
}

</script>
