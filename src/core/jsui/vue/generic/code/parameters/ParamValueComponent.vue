<template>
    <v-list-item
        class="pa-0"
        v-if="!!param"
    >
        <v-list-item-content class="mr-1">
            <v-text-field
                :value="param.Name"
                filled
                label="Name"
                hide-details
                dense
                readonly
            >
            </v-text-field>
        </v-list-item-content>

        <v-list-item-content class="ml-1">
            <v-text-field
                :value="value"
                @input="input"
                filled
                label="Value"
                hide-details
                dense
                :readonly="readonly"
            >
            </v-text-field>
        </v-list-item-content>
    </v-list-item>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import {
    CodeParamType,
} from '../../../../ts/code'

const Props = Vue.extend({
    props: {
        value: {
            type: [Object, Number, String, Date],
            default: () => null as any | null
        },
        param: {
            type: Object,
            default: () => null as CodeParamType | null
        },
        readonly: {
            type: Boolean,
            default: false,
        }
    }
})

@Component
export default class ParamValueComponent extends Props {
    input(v : any) {
        this.$emit('input', v)
    }
}

</script>
