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
            <string-filter
                label="Type"
                :value="value.ResourceTypeFilter"
                @input="onChangeTypeFilter"
                class="mb-2"
            >
            </string-filter>

            <string-filter
                label="Action"
                :value="value.ActionFilter"
                @input="onChangeActionFilter"
                class="mb-2"
            >
            </string-filter>

            <string-filter
                label="User"
                :value="value.UserFilter"
                @input="onChangeUserFilter"
                class="mb-2"
            >
            </string-filter>

            <time-range-filter
                label="Time Range"
                :value="value.TimeRangeFilter"
                @input="onChangeTimeRangeFilter"
            >
            </time-range-filter>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { StringFilterData, TimeRangeFilterData } from '../../../ts/filters'
import { AuditTrailFilterData } from '../../../ts/auditTrail'
import StringFilter from './components/StringFilter.vue'
import TimeRangeFilter from './components/TimeRangeFilter.vue'

const Props = Vue.extend({
    props: {
        value: {
            type: Object,
            default: () => Object() as AuditTrailFilterData
        }
    }
})

@Component({
    components: {
        StringFilter,
        TimeRangeFilter
    }
})
export default class AdvancedAuditTrailFilters extends Props {
    showHideFilters: boolean = false

    onChangeTypeFilter(f : StringFilterData) {
        this.value.ResourceTypeFilter = f
        this.$emit('input', this.value)
    }

    onChangeActionFilter(f : StringFilterData) {
        this.value.ActionFilter = f
        this.$emit('input', this.value)
    }

    onChangeUserFilter(f : StringFilterData) {
        this.value.UserFilter = f
        this.$emit('input', this.value)
    }

    onChangeTimeRangeFilter(f : TimeRangeFilterData) {
        this.value.TimeRangeFilter = f
        this.$emit('input', this.value)
    }
}

</script>
