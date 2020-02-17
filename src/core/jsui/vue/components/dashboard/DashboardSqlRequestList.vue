<template>
    <div>
        <v-row v-if="!requestsToDisplay" align="center" justify="center">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-row>

        <div v-else>
            <v-list-item class="pa-0">
                <v-list-item-content class="disable-flex mr-4">
                    <v-list-item-title class="title">
                        SQL Query Requests
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

                <v-list-item-action>
                    <v-dialog v-model="showHideNew"
                              persistent
                              max-width="40%"
                    >
                        <template v-slot:activator="{ on }">
                            <v-btn color="primary" v-on="on">
                                New
                            </v-btn>
                        </template>

                    </v-dialog>
                </v-list-item-action>
            </v-list-item>
            <v-divider></v-divider>

            <sql-request-table
                :resources="requestsToDisplay"
                :search="filterText"
            >
            </sql-request-table>
        </div>
    </div>
</template>


<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { DbSqlQueryRequest } from '../../../ts/sql'
import { DatabaseStore, getStoreForDatabase } from '../../../ts/vuex/databaseStore'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'
import {
    allSqlRequest, TAllSqlRequestOutput,
} from '../../../ts/api/apiSqlRequests'
import SqlRequestTable from '../../generic/SqlRequestTable.vue'

const Props = Vue.extend({
    props: {
        dbId: {
            type: Number,
            default: -1,
        }
    }
})

@Component({
    components: {
        SqlRequestTable
    }
})
export default class DashboardSqlRequestList extends Props {
    showHideNew : boolean = false
    filterText : string = ""
    store : DatabaseStore | null = this.dbId != -1 ? getStoreForDatabase(this.dbId) : null
    orgRequests : DbSqlQueryRequest[] = []

    get requestsToDisplay() :  DbSqlQueryRequest[] | null {
        if (!!this.store) {
            return this.store.state.allRequests
        }
        return this.orgRequests
    }

    refreshData() {
        if (!!this.store) {
            this.store.dispatch('pullRequests')
        } else {
            allSqlRequest({
                orgId: PageParamsStore.state.organization!.Id,
            }).then((resp : TAllSqlRequestOutput) => {
                this.orgRequests = resp.data
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
    }

    mounted() {
        this.refreshData()
    }
}

</script>
