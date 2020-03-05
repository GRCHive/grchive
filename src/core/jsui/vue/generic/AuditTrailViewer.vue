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

        <div v-if="!!auditTrail">
            <audit-entry-table :resources="auditTrail">
            </audit-entry-table>
        </div>

        <v-row justify="center" align="center" v-else>
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-row>
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
    auditTrail : AuditEventEntry[] | null = null

    refreshData() {
        allAuditTrail({
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllAuditTrailOutput) => {
            this.auditTrail = resp.data
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    mounted() {
        this.refreshData()
    }
}

</script>
