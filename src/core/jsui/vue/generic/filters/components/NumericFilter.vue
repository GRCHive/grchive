<template>
    <v-select
        :value="comparison"
        label="Operator"
        :items="comparisonOperatorItems"
        dense
        filled hide-details
        height="44px"
        @input="changeComparisonOperator"
    >
        <template v-slot:prepend>
            <span class="font-weight-bold numeric-filter-label">{{ label }}</span>
        </template>

        <template v-slot:append-outer>
            <v-text-field
                class="numeric-filter-val"
                :value="compareTo"
                @change="changeCompareTo"
                type="number"
                filled
                v-if="!filterDisabled"
                hide-details
                dense
            >
            </v-text-field>
        </template>
    </v-select>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import * as filters from '../../../../ts/filters'

const Props = Vue.extend({
    props: {
        label: {
            type: String,
            default: ""
        },
    }
})

@Component
export default class NumericFilter extends Props {
    comparison : filters.ComparisonOperators = filters.ComparisonOperators.Disabled
    compareTo : number = 0

    comparisonOperatorItems : any[] = filters.comparisonOperatorsSelectItems

    get filterDisabled() : boolean {
        return this.comparison == filters.ComparisonOperators.Disabled
    }

    changeComparisonOperator(c : filters.ComparisonOperators) {
        this.comparison = c
    }

    changeCompareTo(val : string) {
        this.compareTo = parseInt(val, 10)
    }
}

</script>

<style scoped>

.numeric-filter-label {
    transform: translateY(-6px);
}

.numeric-filter-val {
    transform: translateY(-14px);
    margin-bottom: -6px;
}

>>>.numeric-filter-val input {
    padding-top: 0!important;
    padding-bottom: 0!important;
    margin-top: 11px !important;
    margin-bottom: 11px !important;
}

</style>
