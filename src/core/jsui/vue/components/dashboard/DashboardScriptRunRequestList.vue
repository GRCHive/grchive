<template>
    <div>
        <v-row v-if="isLoading" align="center" justify="center">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-row>

        <div v-else>
            <v-list-item class="pa-0">
                <v-list-item-content class="disable-flex mr-4">
                    <v-list-item-title class="title">
                        Script Run Requests
                    </v-list-item-title>
                </v-list-item-content>
                <v-list-item-action>
                    <v-text-field outlined
                                  v-model="filterText"
                                  prepend-inner-icon="mdi-magnify"
                                  hide-details
                    ></v-text-field>
                </v-list-item-action>
                <v-spacer></v-spacer>

            </v-list-item>
            <v-divider></v-divider>

            <script-run-request-table
                :resources="requests"
                :search="filterText"
            >
            </script-run-request-table>
        </div>
    </div>

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import {
    GenericRequest,
    GenericApproval,
} from '../../../ts/requests'
import { allGenericRequests, TAllGenericRequestsOutput } from '../../../ts/api/apiRequests'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import ScriptRunRequestTable from '../../generic/ScriptRunRequestTable.vue'

@Component({
    components: {
        ScriptRunRequestTable
    }
})
export default class DashboardScriptRunRequestList extends Vue {
    filterText: string = ""
    requests : GenericRequest[] | null = null

    get isLoading() : boolean {
        return !this.requests
    }

    refreshData() {
        allGenericRequests({
            orgId: PageParamsStore.state.organization!.Id,
            scriptsOnly: true,
        }).then((resp : TAllGenericRequestsOutput) => {
            this.requests = resp.data
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops. Something went wrong. Try again.",
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
