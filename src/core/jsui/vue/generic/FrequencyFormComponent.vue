<template>
    <section id="freq">
        <v-checkbox
            label="Is Manually Run"
            :value="isManual"
            @change="changeManual">
        </v-checkbox>
        <v-text-field
            label="Interval"
            prefix="Repeat every"
            outlined
            type="number"
            :rules="[rules.numeric]"
            v-if="!isManual"
            :value="freqInterval"
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
        isManual : Boolean,
        freqInterval : Number,
        freqType: Number,
    },
    data: () => ({
        rules
    }),
    computed: {
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
        changeManual(val: boolean) {
            this.$emit('update:isManual', val)
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
