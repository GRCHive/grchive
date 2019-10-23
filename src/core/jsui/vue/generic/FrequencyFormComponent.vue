<template>
    <section id="freq">
        <v-checkbox
            v-model="isManual"
            label="Is Manually Run">
        </v-checkbox>
        <v-text-field
            label="Interval"
            prefix="Repeat every"
            outlined
            type="number"
            v-model="freqInterval"
            :rules="[rules.numeric]"
            v-if="!isManual"
        >
            <template v-slot:append-outer v-bind:freqType="freqType">
                <section id="choices">
                    <v-select :items="frequencyChoices" outlined v-model="freqType">
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
    data : () => ({
        isManual : false,
        freqInterval : 0,
        freqType: 0,
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
