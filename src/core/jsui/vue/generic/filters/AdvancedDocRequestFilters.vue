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
            <doc-request-status-filter
                label="Status"
                :value="value.StatusFilter"
                @input="onChangeStatus"
            >
            </doc-request-status-filter>

            <div class="d-flex my-4">
                <user-filter
                    label="Requester"
                    :value="value.RequesterFilter"
                    @input="onChangeRequester"
                    class="mr-4"
                >
                </user-filter>

                <user-filter
                    label="Assignee"
                    select-noone
                    :value="value.AssigneeFilter"
                    @input="onChangeAssignee"
                >
                </user-filter>
            </div>

            <time-range-filter
                label="Request Time"
                :value="value.RequestTimeFilter"
                @input="onChangeRequestTime"
            >
            </time-range-filter>

            <time-range-filter
                label="Due Date"
                :value="value.DueDateFilter"
                @input="onChangeDueDate"
            >
            </time-range-filter>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Prop } from 'vue-property-decorator'
import { DocRequestFilterData, DocRequestStatusFilterData } from '../../../ts/docRequests'
import { TimeRangeFilterData, UserFilterData } from '../../../ts/filters'
import TimeRangeFilter from './components/TimeRangeFilter.vue'
import DocRequestStatusFilter from './components/DocRequestStatusFilter.vue'
import UserFilter from './components/UserFilter.vue'

@Component({
    components: {
        TimeRangeFilter,
        DocRequestStatusFilter,
        UserFilter,
    }
})
export default class AdvancedDocRequestFilters extends Vue {
    @Prop({required: true})
    value! : DocRequestFilterData

    showHideFilters: boolean = false

    onChangeRequestTime(f : TimeRangeFilterData) {
        this.value.RequestTimeFilter = f
        this.$emit('input', this.value)
    }

    onChangeDueDate(f : TimeRangeFilterData) {
        this.value.DueDateFilter = f
        this.$emit('input', this.value)
    }

    onChangeStatus(f : DocRequestStatusFilterData) {
        this.value.StatusFilter = f
        this.$emit('input', this.value)
    }

    onChangeRequester(f : UserFilterData) {
        this.value.RequesterFilter = f
        this.$emit('input', this.value)
    }

    onChangeAssignee(f : UserFilterData) {
        this.value.AssigneeFilter = f
        this.$emit('input', this.value)
    }
}

</script>
