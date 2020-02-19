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
                label="Number of Related Controls"
            >
            </numeric-filter>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import NumericFilter from './components/NumericFilter.vue'

@Component({
    components: {
        NumericFilter
    }
})
export default class AdvancedRiskFilters extends Vue {
    showHideFilters: boolean = false
}

</script>
