<template>
    <div class="d-flex">
        <span class="font-weight-bold label-helper mr-4">{{ label }}</span>
        <v-autocomplete
            :items="items"
            dense
            multiple
            hide-details
            hide-selected
            :value="value.ValidStatuses"
            @input="onChange"
            chips
            clearable
            filled
        >
            <template v-slot:selection="{item}">
                <doc-request-status-display
                    :status="item.value"
                    can-close
                    @click:close="onClose(item)"
                >
                </doc-request-status-display>
            </template>

            <template v-slot:item="{item}">
                <doc-request-status-display
                    :status="item.value"
                >
                </doc-request-status-display>
            </template>
        </v-autocomplete>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Prop } from 'vue-property-decorator'
import {
    DocRequestStatus,
    DocRequestStatusFilterData,
} from '../../../../ts/docRequests'
import DocRequestStatusDisplay from '../../requests/DocRequestStatusDisplay.vue'

@Component({
    components: {
        DocRequestStatusDisplay,
    }
})
export default class DocRequestStatusFilter extends Vue {
    @Prop()
    label! : string

    @Prop()
    value!: DocRequestStatusFilterData

    onClose(item : any) {
        const idx = this.value.ValidStatuses.findIndex((ele : DocRequestStatus) => ele == item.value)
        if (idx == -1) {
            return
        }
        this.value.ValidStatuses.splice(idx, 1)
        this.$emit('input', this.value)
    }

    get items() : any[] {
        return [
            {
                text: 'Open',
                value: DocRequestStatus.Open,
            },
            {
                text: 'In Progress',
                value: DocRequestStatus.InProgress,
            },
            {
                text: 'Feedback',
                value: DocRequestStatus.Feedback,
            },
            {
                text: 'Complete',
                value: DocRequestStatus.Complete,
            },
            {
                text: 'Overdue',
                value: DocRequestStatus.Overdue,
            },
        ]
    }

    onChange(v : DocRequestStatus[]) {
        this.value.ValidStatuses = v
        this.$emit('input', this.value)
    }
}

</script>

<style scoped>

.label-helper { 
    margin-top: 10px;
}

</style>
