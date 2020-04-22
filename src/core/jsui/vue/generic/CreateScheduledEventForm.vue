<template>
    <div v-if="!!value">
        <v-text-field
            v-model="value.Name"
            label="Name"
            filled
            :rules="[rules.required, rules.createMaxLength(256)]"
        >
        </v-text-field> 

        <v-textarea
            v-model="value.Description"
            label="Description"
            filled
            hide-details
        >
        </v-textarea> 

        <v-checkbox
            :input-value="value.Repeat"
            label="Repeat?"
            @click="changeRepeat"
        >
        </v-checkbox>

        <div v-if="value.Repeat">
            <v-select
                label="Frequency"
                :items="frequencyItems"
                filled
                hide-details
                :value="value.Frequency"
                @input="changeFrequency"
            >
            </v-select>

            <div v-if="value.Frequency == CronFrequency.Daily && !!value.Daily">
                <v-row align="center">
                    <v-col cols="3">
                        <span>Run every day at</span>
                    </v-col>

                    <v-col cols="9">
                        <div class="time-container" v-for="(t, index) in value.Daily.Times" :key="`container${index}`">
                            <date-time-picker-form-component
                                class="time-time"
                                v-model="value.Daily.Times[index]"
                                disable-date
                            >
                            </date-time-picker-form-component>

                            <v-btn
                                class="time-delete"
                                color="error"
                                icon
                            >
                                <v-icon>mdi-delete</v-icon>
                            </v-btn>
                        </div>

                        <v-btn
                            color="primary"
                            block
                            outlined
                            @click="value.Daily.Times.push(new Date())"
                        >
                            Add
                        </v-btn>
                    </v-col>
                </v-row>
            </div>

            <div v-else-if="value.Frequency == CronFrequency.Weekly && !!value.Weekly">
                <v-row align="center">
                    <v-col cols="3">
                        <span>Run every week on:</span>
                    </v-col>

                    <v-col cols="9">
                        <week-day-form-component
                            v-model="value.Weekly.Days"
                        >
                        </week-day-form-component>
                    </v-col>
                </v-row>

                <v-row align="center">
                    <v-col cols="3">
                        <span>At:</span>
                    </v-col>

                    <v-col cols="9">
                        <date-time-picker-form-component
                            v-model="value.Weekly.Time"
                            disable-date
                        >
                        </date-time-picker-form-component>


                    </v-col>
                </v-row>
            </div>

            <div v-else-if="value.Frequency == CronFrequency.Monthly && !!value.Monthly">
                <v-row align="center">
                    <v-col cols="3">
                        <span>Run every month on:</span>
                    </v-col>

                    <v-col cols="9">
                        <v-switch
                            v-model="value.Monthly.UseDate"
                            label="Date Mode"
                        >
                        </v-switch>
                        
                        <v-select
                            v-model="value.Monthly.Dates"
                            chips
                            deletable-chips
                            multiple
                            outlined
                            label="Dates"
                            :items="dateItems"
                            hide-details
                            v-if="value.Monthly.UseDate"
                        >
                            <template v-slot:append>
                                <span style="margin-top: 10px;">days of the month</span>
                            </template>
                        </v-select>

                        <div class="cron-container" v-else>
                            <span>Every</span>
                            <v-select
                                class="ml-2"
                                v-model="value.Monthly.Nth"
                                outlined
                                :items="CronWeekdayHashItems"
                                hide-details
                            >
                            </v-select>
                            <v-select
                                class="ml-2"
                                v-model="value.Monthly.Day"
                                outlined
                                :items="DaysSelectItems"
                                hide-details
                            >
                            </v-select>
                        </div>
                    </v-col>
                </v-row>

                <v-row align="center">
                    <v-col cols="3">
                        <span>At:</span>
                    </v-col>

                    <v-col cols="9">
                        <date-time-picker-form-component
                            v-model="value.Monthly.Time"
                            disable-date
                        >
                        </date-time-picker-form-component>
                    </v-col>
                </v-row>
            </div>
        </div>
        
        <div v-else>
            <date-time-picker-form-component
                :value="value.OneTimeDate"
                @input="changeOneTimeDate"
            >
            </date-time-picker-form-component>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import {
    ScheduledEvent,
    CronFrequency,
    createEmptyScheduledEvent,
    cleanScheduledEventFromJson,
    createEmptyDailyCron,
    createEmptyWeeklyCron,
    createEmptyMonthlyCron,
    CronWeekdayHashItems
} from '../../ts/event'
import {
    DaysSelectItems
} from '../../ts/time'
import DateTimePickerFormComponent from './DateTimePickerFormComponent.vue'
import WeekDayFormComponent from './WeekDayFormComponent.vue'
import Ordinal from 'ordinal'
import * as rules from '../../ts/formRules'

const Props = Vue.extend({
    props: {
        value: {
            type: Object,
            default: () => null as ScheduledEvent | null
        }
    }
})

@Component({
    components: {
        DateTimePickerFormComponent,
        WeekDayFormComponent,
    }
})
export default class CreateScheduledEventForm extends Props {
    CronFrequency : any = CronFrequency
    CronWeekdayHashItems: any = CronWeekdayHashItems
    DaysSelectItems: any = DaysSelectItems
    rules : any = rules

    doChange(fn : (e : ScheduledEvent) => void) {
        let e : ScheduledEvent
        if (!this.value) {
            e = createEmptyScheduledEvent()
        } else {
            e = JSON.parse(JSON.stringify(this.value))
            cleanScheduledEventFromJson(e)
        }

        fn(e)

        this.$emit('input', e)
    }

    changeRepeat(e : MouseEvent) {
        // Clicking on the checkbox throws this event twice for
        // whatever reason and the @input event doesn't work on
        // v-checkbox...for whatever reason. v-model works fine though.
        e.stopPropagation()

        this.doChange((e : ScheduledEvent) => {
            e.Repeat = !e.Repeat
        })
    }

    changeOneTimeDate(d : Date) {
        this.doChange((e : ScheduledEvent) => {
            e.OneTimeDate = d
        })
    }

    changeFrequency(f : CronFrequency) {
        this.doChange((e : ScheduledEvent) => {
            e.Frequency = f
            if (f == CronFrequency.Daily) {
                e.Daily = createEmptyDailyCron()
                e.Weekly = null
                e.Monthly = null
            } else if (f == CronFrequency.Weekly) {
                e.Daily = null
                e.Weekly = createEmptyWeeklyCron()
                e.Monthly = null
            } else if (f == CronFrequency.Monthly) {
                e.Daily = null
                e.Weekly = null
                e.Monthly = createEmptyMonthlyCron()
            }
        })
    }

    get frequencyItems() : any[] {
        return Object.keys(CronFrequency)
            .filter((key : any) => !isNaN(Number(CronFrequency[key])))
            .map((key : any) => ({
                text: key,
                value: CronFrequency[key]
            }))
    }

    get dateItems() : any[] {
        return Array.from(Array(31).keys()).map((d : number) => {
            // 0 - 30 to 1-31
            let cd = d + 1

            return {
                text: Ordinal(cd),
                value: cd,
            }
        })
    }

    mounted() {
        if (!this.value) {
            this.$emit('input', createEmptyScheduledEvent())
        }
    }
}

</script>

<style scoped>

.cron-container {
    display: flex;
    align-items: center;
}

.time-container {
    display: flex;
    align-items: center;
}

.time-time {
    flex-grow: 8;
}

.time-delete {
    flex-grow: 0;
}

</style>
