<template>
    <div>
        <v-list-item>
            <v-list-item-action>
                <v-btn @click="resetToToday">
                    Today
                </v-btn>
            </v-list-item-action>

            <v-list-item-action>
                <v-btn icon @click="prev">
                    <v-icon>mdi-chevron-left</v-icon>
                </v-btn>
            </v-list-item-action>

            <v-list-item-action>
                <v-btn icon @click="next">
                    <v-icon>mdi-chevron-right</v-icon>
                </v-btn>
            </v-list-item-action>

            <v-list-item-content>
                {{ monthYear }}
            </v-list-item-content>
        </v-list-item>
        <div style="max-height: calc(100% - 60px); height: calc(100% - 60px);">
            <v-calendar
                ref="calendar"
                v-model="currentTimeStr"
                type="month"
                :events="scheduledVueitfyEvents"
                @change="refreshSchedule"
            >
            </v-calendar>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { VCalendar } from 'vuetify/lib'
import { TAllScheduledTasksOutput, allScheduledTasks } from '../../../ts/api/apiTasks'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import { ScheduledTaskMetadata } from '../../../ts/tasks'
import {
    TimeRange,
    vuetifyCalendarTimeFormat,
    standardFormatDate,
    parseStandardFormatDate
} from '../../../ts/time'

@Component
export default class ScheduleViewer extends Vue {
    currentTimeStr : string = standardFormatDate(new Date())

    allTasks : ScheduledTaskMetadata[] = []
    taskTimes : Record<number, TimeRange[]> = Object()

    resetToToday() {
        this.currentTimeStr = standardFormatDate(new Date())
    }

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

    get currentTime() : Date {
        return parseStandardFormatDate(this.currentTimeStr)
    }

    get monthYear() : string {
        let opts = new Intl.DateTimeFormat('en-US', { month : 'long' })
        return `${opts.format(this.currentTime)} ${this.currentTime.getFullYear()}`
    }

    prev() {
        //@ts-ignore
        this.$refs.calendar.prev()
    }

    next() {
        //@ts-ignore
        this.$refs.calendar.next()
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
