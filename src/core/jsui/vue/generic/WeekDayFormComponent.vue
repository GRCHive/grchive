<template>
    <div class="day-container">
        <v-btn
            color="primary"
            :outlined="!activeDays[Days.Sunday]"
            @click="toggleActive(Days.Sunday)"
            class="my-1"
            :disabled="readonly"
        >
            Sunday
        </v-btn>

        <v-btn
            color="primary"
            :outlined="!activeDays[Days.Monday]"
            @click="toggleActive(Days.Monday)"
            class="ml-1 my-1"
            :disabled="readonly"
        >
            Monday
        </v-btn>

        <v-btn
            color="primary"
            :outlined="!activeDays[Days.Tuesday]"
            @click="toggleActive(Days.Tuesday)"
            class="ml-1 my-1"
            :disabled="readonly"
        >
            Tuesday
        </v-btn>

        <v-btn
            color="primary"
            :outlined="!activeDays[Days.Wednesday]"
            @click="toggleActive(Days.Wednesday)"
            class="ml-1 my-1"
            :disabled="readonly"
        >
            Wednesday
        </v-btn>

        <v-btn
            color="primary"
            :outlined="!activeDays[Days.Thursday]"
            @click="toggleActive(Days.Thursday)"
            class="ml-1 my-1"
            :disabled="readonly"
        >
            Thursday
        </v-btn>

        <v-btn
            color="primary"
            :outlined="!activeDays[Days.Friday]"
            @click="toggleActive(Days.Friday)"
            class="ml-1 my-1"
            :disabled="readonly"
        >
            Friday
        </v-btn>

        <v-btn
            color="primary"
            :outlined="!activeDays[Days.Saturday]"
            @click="toggleActive(Days.Saturday)"
            class="ml-1 my-1"
            :disabled="readonly"
        >
            Saturday
        </v-btn>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { Days, DaysKey } from '../../ts/time'

const Props = Vue.extend({
    props: {
        value: Array as () => Array<Days>,
        default: () => [],
        readonly: {
            type:  Boolean,
            default: false,
        }
    }
})

@Component
export default class WeekDayFormComponent extends Props {
    Days: any = Days
    activeDays : Record<Days, boolean> = Object()

    toggleActive(d : Days) {
        Vue.set(this.activeDays, d, !this.activeDays[d])
        this.syncToValue()
    }

    syncToValue() {
        this.$emit(
            'input',
            Object.keys(Days)
                .filter((key : any) => !isNaN(Number(Days[key])))
                .filter((key : any) => this.activeDays[Days[<DaysKey>key]])
                .map((key : any) => Days[key])
        )
    }

    @Watch('value')
    syncFromValue() {
        let set : Set<Days> = new Set<Days>(this.value)
        Object.keys(Days).filter((key : any) => !isNaN(Number(Days[key]))).forEach((key : any) => {
            Vue.set(this.activeDays, Days[<DaysKey>key], set.has(Days[<DaysKey>key]))
        })
    }

    mounted() {
        this.syncFromValue()
    }
}

</script>

<style scoped>

.day-container {
    display: flex;
    flex-wrap: wrap;
}

</style>
