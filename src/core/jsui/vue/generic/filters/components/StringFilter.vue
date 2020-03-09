<template>
    <div class="d-flex">
        <span class="font-weight-bold string-label mr-4">{{ label }}</span>
        <v-select
            class="mr-4 flex-grow-0"
            :value="value.Op"
            label="Operator"
            :items="comparisonOperatorItems"
            dense
            outlined
            hide-details
            @input="changeComparisonOperator"
        >
        </v-select>

        <v-text-field
            class="flex-grow-0"
            :value="value.Target"
            @input="changeCompareTo"
            outlined
            v-if="!filterDisabled"
            hide-details
            dense
        >
        </v-text-field>

    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { StringFilterData } from '../../../../ts/filters'
import * as filters from '../../../../ts/filters'

const Props = Vue.extend({
    props: {
        label: {
            type: String,
            default: ""
        },
        value: {
            type: Object,
            default: () => Object() as StringFilterData
        },
    }
})

@Component
export default class StringFilter extends Props {
    comparisonOperatorItems : any[] = filters.stringComparisonOperatorsSelectItems

    get filterDisabled() : boolean {
        return this.value.Op == filters.StringComparisonOperators.Disabled
    }

    changeComparisonOperator(c : filters.StringComparisonOperators) {
        this.value.Op = c
        this.$emit('input', this.value)
    }

    changeCompareTo(val : string) {
        this.value.Target = val
        this.$emit('input', this.value)
    }
}

</script>

<style scoped>

.string-label { 
    margin-top: 10px;
}

</style>
