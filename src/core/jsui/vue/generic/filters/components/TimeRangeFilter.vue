<template>
    <div class="d-flex">
        <span class="font-weight-bold string-label mr-4">{{ label }}</span>
        <v-menu
            offset-y
            v-model="showStartDate"
            :close-on-content-click="false"
            content-class="dt-menu"
        >
            <template v-slot:activator="{ on }">
                <v-text-field
                    :value="startStr"
                    label="Start Time"
                    prepend-icon="mdi-calendar"
                    readonly
                    v-on="on">
                </v-text-field>
            </template>

            <date-time-picker
                v-model="value.Start"
                class="dt-bg"
            >
            </date-time-picker>
        </v-menu>

        <v-menu
            offset-y
            v-model="showEndDate"
            :close-on-content-click="false"
            content-class="dt-menu"
        >
            <template v-slot:activator="{ on }">
                <v-text-field
                    :value="endStr"
                    label="End Time"
                    prepend-icon="mdi-calendar"
                    readonly
                    v-on="on">
                </v-text-field>
            </template>

            <date-time-picker
                v-model="value.End"
                class="dt-bg"
            >
            </date-time-picker>
        </v-menu>

    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { TimeRangeFilterData } from '../../../../ts/filters'
import { standardFormatTime } from '../../../../ts/time'
import DateTimePicker from '../../DateTimePicker.vue'

const Props = Vue.extend({
    props: {
        label: {
            type: String,
            default: ""
        },
        value: {
            type: Object,
            default: () => Object() as TimeRangeFilterData
        },
    }
})

@Component({
    components: {
        DateTimePicker
    }
})
export default class TimeRangeFilter extends Props {
    showStartDate : boolean = false
    showEndDate: boolean = false
    
    get startStr() : string {
        return standardFormatTime(this.value.Start)
    }

    get endStr() : string {
        return standardFormatTime(this.value.End)
    }
}

</script>

<style scoped>

.string-label { 
    margin-top: 10px;
}

.dt-bg {
    background-color: transparent;
}

.dt-menu {
    box-shadow: none !important;
}

</style>
