<template>
    <v-row>
        <v-date-picker v-if="!disableDate"
            class="mr-4"
            :value="date"
            @input="onDateChange"
        >
        </v-date-picker>
        <v-time-picker
            :value="time"
            @input="onTimeChange"
        >
        </v-time-picker>
    </v-row>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'

const Props = Vue.extend({
    props: {
        value: {
            type: Date,
            default: () => new Date(),
        },
        disableDate: {
            type: Boolean,
            default: false,
        },
    }
})

@Component
export default class DateTimePicker extends Props {
    get realDate() : Date {
        return new Date(this.value)
    }

    get date() : string {
        return `${this.realDate.getFullYear()}-${(this.realDate.getMonth()+1).toString().padStart(2, "0")}-${this.realDate.getDate().toString().padStart(2, "0")}`
    }

    get time(): string {
        return `${this.realDate.getHours().toString().padStart(2, "0")}:${this.realDate.getMinutes().toString().padStart(2, "0")}`
    }

    onDateChange(d : string) {
        let n : Date = new Date(this.value)

        let data : string[] = d.split('-')
        n.setFullYear(parseInt(data[0]))
        n.setMonth(parseInt(data[1])-1)
        n.setDate(parseInt(data[2]))

        this.$emit('input', n)
    }

    onTimeChange(t : string) {
        let n : Date = new Date(this.value)

        let data : string[] = t.split(':')
        n.setHours(parseInt(data[0]))
        n.setMinutes(parseInt(data[1]))
        this.$emit('input', n)
    }
}

</script>
