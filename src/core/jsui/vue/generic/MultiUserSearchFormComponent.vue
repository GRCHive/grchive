<template>
    <v-autocomplete
        :label="label"
        filled
        cache-items
        :items="displayItems"
        :loading="loading"
        hide-no-data
        hide-selected
        :clearable="!readonly && !disabled"
        :value="value"
        @change="changeUser"
        :disabled="disabled"
        :readonly="readonly"
        :value-comparator="compare"
        :hide-details="hideDetails"
        multiple
        chips
        deletable-chips
        :dense="dense"
    ></v-autocomplete>
</template>

<script lang="ts">

import Vue from 'vue'
import { createUserString } from '../../ts/users'
import Metadata from '../../ts/metadata'

export default Vue.extend({
    props: {
        label: String,
        value: {
            type: Array,
        },
        disabled: {
            type: Boolean,
            default: false
        },
        readonly: {
            type: Boolean,
            default: false
        },
        hideDetails: {
            type: Boolean,
            default: false
        },
        selectNoone: {
            type: Boolean,
            default: false
        },
        dense: {
            type: Boolean,
            default: false
        }
    },
    data : () => ({
        loading: false,
    }),
    computed: {
        displayItems() : Object[] {
            let displayText = [] as Object[]
            for (let user of Metadata.state.availableUsers) {
                displayText.push({
                    text: createUserString(user),
                    value: user
                })
            }

            if (this.selectNoone) {
                displayText.push({
                    text: 'No User',
                    value: null,
                })
            }

            return displayText
        },
    },
    methods: {
        changeUser(val : User[]) {
            this.$emit('input', val)
        },
        compare(a : User | null, b : User | null) : boolean {
            if (a == null && b == null) {
                return true
            }

            if (a == null || b == null) {
                return false
            }
            return a.Id == b.Id
        }
    },
})
</script>
