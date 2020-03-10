<template>
    <div>
        <v-list-item v-if="!noHeader">
            <v-list-item-content>
                <v-list-item-title class="title">
                    Audit Trail
                </v-list-item-title>
            </v-list-item-content>
        </v-list-item>
        <advanced-audit-trail-filters
            v-model="filter"
            class="px-4"
        >
        </advanced-audit-trail-filters>
        <v-divider></v-divider>
        <audit-entry-table
            :retrieval-params="auditParams"
        ></audit-entry-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { AuditEventEntry, AuditTrailFilterData, NullAuditTrailFilterData } from '../../ts/auditTrail'
import { PageParamsStore } from '../../ts/pageParams'
import { contactUsUrl } from '../../ts/url'
import { TAllAuditTrailOutput, allAuditTrail} from '../../ts/api/apiAuditTrail'
import AuditEntryTable from './AuditEntryTable.vue'
import AdvancedAuditTrailFilters from './filters/AdvancedAuditTrailFilters.vue'

const Props = Vue.extend({
    props: {
        resourceType: {
            type: Array,
        },
        resourceId: {
            type: Array,
        },
        noHeader: {
            type: Boolean,
            default: false,
        },
    }
})

@Component({
    components: {
        AuditEntryTable,
        AdvancedAuditTrailFilters
    }
})
export default class AuditTrailViewer extends Props {
    filter : AuditTrailFilterData = NullAuditTrailFilterData

    get auditParams() : any {
        let params : any = {
            filter: this.filter,
        }
        params.resourceType = this.resourceType
        params.resourceId = this.resourceId
        return params
    }
}

</script>
