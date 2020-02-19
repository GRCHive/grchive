<template>
    <div class="d-flex">
        <span class="font-weight-bold numeric-label mr-4">{{ label }}</span>
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
            type="number"
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
import * as filters from '../../../../ts/filters'
import { NumericFilterData } from '../../../../ts/filters'

const Props = Vue.extend({
    props: {
        label: {
            type: String,
            default: ""
        },
        value: {
            type: Object,
            default: () => Object() as NumericFilterData
        },
    }
})

@Component
export default class NumericFilter extends Props {
    comparisonOperatorItems : any[] = filters.comparisonOperatorsSelectItems

    get filterDisabled() : boolean {
        return this.value.Op == filters.ComparisonOperators.Disabled
    }

    changeComparisonOperator(c : filters.ComparisonOperators) {
        this.value.Op = c
        this.$emit('input', this.value)
    }

    changeCompareTo(val : string) {
        this.value.Target = parseInt(val, 10)
        this.$emit('input', this.value)
    }
}

</script>

<style scoped>

.numeric-label { 
    margin-top: 10px;
}

</style>
