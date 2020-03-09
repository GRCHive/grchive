<template>
    <div>
        <v-list-item>
            <v-list-item-content>
                <v-list-item-title class="title">
                    Audit Trail
                </v-list-item-title>
            </v-list-item-content>
        </v-list-item>
        <v-divider></v-divider>
        <audit-entry-table></audit-entry-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { AuditEventEntry } from '../../ts/auditTrail'
import { PageParamsStore } from '../../ts/pageParams'
import { contactUsUrl } from '../../ts/url'
import { TAllAuditTrailOutput, allAuditTrail} from '../../ts/api/apiAuditTrail'
import AuditEntryTable from './AuditEntryTable.vue'

@Component({
    components: {
        AuditEntryTable
    }
})
export default class AuditTrailViewer extends Vue {
}

</script>
