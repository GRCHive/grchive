<template>
    <section id="freq">
        <v-checkbox
            label="Run Ad Hoc"
            :value="isAdHoc"
            :disabled="disabled"
            :readonly="readonly"
            @change="changeAdHoc">
            hide-details
        </v-checkbox>
        <v-text-field
            label="Interval"
            prefix="Repeat every"
            outlined
            type="number"
            :rules="[rules.numeric]"
            v-if="!isAdHoc"
            :value="freqInterval"
            :disabled="disabled"
            :readonly="readonly"
            @change="changeInterval"
        >
            <template v-slot:append-outer v-bind:freqType="freqType">
                <section id="choices">
                    <v-select :items="frequencyChoices"
                              outlined
                              :value="freqType"
                              @change="changeType">
                    </v-select>
                </section>
            </template>
        </v-text-field>
    </section>
</template>

<script lang="ts">

import Vue from 'vue'
import { frequencyTypes } from '../../ts/frequency'
import * as rules from "../../ts/formRules"

export default Vue.extend({
    props: {
        freqInterval : Number,
        freqType: Number,
        disabled: {
            type: Boolean,
            default: false
        },
        readonly: {
            type: Boolean,
            default: false
        }
    },
    data: () => ({
        rules
    }),
    computed: {
        isAdHoc() : boolean {
            return (this.freqType == -1)
        },
        frequencyChoices() : Object[] {
            let items = [] as Object[]
            let counter = 0
            for (let freq of frequencyTypes) {
                items.push({
                    text: freq,
                    value: counter
                })

                counter += 1
            }
            return items
        }
    },
    methods: {
        changeInterval(val : string) {
            this.$emit('update:freqInterval', parseInt(val, 10))
        },
        changeAdHoc(val: boolean) {
            if (val) {
                this.changeType(-1)
            } else {
                this.changeType(0)
            }
        },
        changeType(val : number) {
            this.$emit('update:freqType', val)
        }
    }
})
</script>

<style scoped>

/* We need to account for the 16px margin that gets tacked onto the v-slot:append-outer
   that I can't figure out how to remove...
 */
#freq {
    margin-bottom: -16px;
}

#choices {
    transform: translateY(-16px);
}

</style>
