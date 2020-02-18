<template>
    <section id="freq">
        <v-select
            label="Frequency Type"
            :value="freqType"
            :items="typeItems"
            @input="changeType(arguments[0], false)"
            :value-comparator="compareType"
            hide-details
            filled
            class="mb-4"
        >
        </v-select>
        <v-text-field
            label="Interval"
            prefix="Repeat every"
            outlined
            type="number"
            :rules="[rules.numeric]"
            v-if="freqType >= 0"
            :value="freqInterval"
            :disabled="disabled"
            :readonly="readonly"
            @change="changeInterval"
            :key="`interval-${key}`"
        >
            <template v-slot:append-outer v-bind:freqType="freqType">
                <div id="choices">
                    <v-select :items="frequencyChoices"
                              outlined
                              :value="freqType"
                              @change="changeType(arguments[0], true)">
                    </v-select>
                </div>
            </template>
        </v-text-field>

        <v-text-field
            label="Other Frequency Description"
            filled
            v-if="freqType == -2"
            :disabled="disabled"
            :readonly="readonly"
            :value="freqOther"
            @change="changeOther"
        >
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
        freqOther: String,
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
        rules,
        key: 0
    }),
    computed: {
        typeItems(): any[] {
            return [
                {
                    text: "Ad-Hoc",
                    value: -1,
                },
                {
                    text: "Other",
                    value: -2,
                },
                {
                    text: "Interval",
                    value: 0,
                },
            ]
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
        compareType(a : number, b: number) : boolean {
            if (a == b) {
                return true
            }

            // All intervals should look the same to the v-select.
            if (a >= 0 && b >= 0) {
                return true
            }

            return false
        },
        changeInterval(val : string) {
            this.$emit('update:freqInterval', parseInt(val, 10))
        },
        changeOther(val : string) {
            this.$emit('update:freqOther', val)
        },
        changeType(val : number, strict : boolean = true) {
            if (!strict && val == 0 && this.freqType >= 0) {
                return
            }
            this.key += 1
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
    min-width: 150px;
}

</style>
