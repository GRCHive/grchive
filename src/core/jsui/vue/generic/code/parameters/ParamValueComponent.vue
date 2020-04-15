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
        param: {
            type: Object,
            default: () => null as CodeParamType | null
        }
    }
})

@Component
export default class ParamValueComponent extends Props {
}

</script>
