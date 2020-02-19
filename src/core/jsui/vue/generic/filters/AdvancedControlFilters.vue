<template>
    <div>
        <p class="ma-0">
            {{ showHideFilters ? "Hide" : "Show" }} Advanced Filters
            <v-btn icon @click="showHideFilters = !showHideFilters">
                <v-icon small v-if="!showHideFilters" >mdi-chevron-down</v-icon>
                <v-icon small v-else>mdi-chevron-up</v-icon>
            </v-btn>
        </p>

        <div class="mb-4" v-if="showHideFilters">
            <numeric-filter
                label="Number of Linked Risks"
                :value="value.NumRisks"
                @input="onChangeNumRisksFilter"
            >
            </numeric-filter>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import NumericFilter from './components/NumericFilter.vue'
import { NumericFilterData } from '../../../ts/filters'
import { ControlFilterData } from '../../../ts/controls'

const Props = Vue.extend({
    props: {
        value: {
            type: Object,
            default: () => Object() as ControlFilterData
        }
    }
})

@Component({
    components: {
        NumericFilter
    }
})
export default class AdvancedControlFilters extends Props {
    showHideFilters: boolean = false

    onChangeNumRisksFilter(f : NumericFilterData) {
        this.value.NumRisks = f
        this.$emit('input', this.value)
    }
}

</script>
