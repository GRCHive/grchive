<template>
    <div>
        <v-list-item>
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

@Component({
    components: {
        AuditEntryTable,
        AdvancedAuditTrailFilters
    }
})
export default class AuditTrailViewer extends Vue {
    filter : AuditTrailFilterData = NullAuditTrailFilterData

    get auditParams() : any {
        return {
            filter: this.filter,
        }
    }
}

</script>
