<template>
    <v-autocomplete
        :label="label"
        filled
        cache-items
        :items="displayItems"
        :loading="loading"
        hide-no-data
        hide-selected
        clearable
        :value="user"
        @change="changeUser"
        :disabled="disabled"
        :value-comparator="compare"
    ></v-autocomplete>
</template>

<script lang="ts">

import Vue from 'vue'
import { createUserString } from '../../ts/users'
import Metadata from '../../ts/metadata'

export default Vue.extend({
    props: {
        label: String,
        user: {
            type: Object as () => User
        },
        disabled: {
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
            return displayText
        },
    },
    methods: {
        changeUser(val : User) {
            this.$emit('update:user', val)
        },
        compare(a : User, b : User) : boolean {
            return a.Id == b.Id
        }
    },
})
</script>
