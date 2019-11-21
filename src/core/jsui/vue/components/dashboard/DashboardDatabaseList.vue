<template>
    <div class="ma-4">
        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    Databases
                </v-list-item-title>
            </v-list-item-content>

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
                    <create-new-database-form
                        @do-cancel="showHideNew = false"
                        @do-save="onSaveDatabase">
                    </create-new-database-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>
        <v-divider></v-divider>
        <db-table
            :resources="allDbs"
        ></db-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import DbTable from '../../generic/DbTable.vue'
import { Database } from '../../../ts/databases'
import { TAllDatabaseOutputs, allDatabase } from '../../../ts/api/apiDatabases'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'
import CreateNewDatabaseForm from './CreateNewDatabaseForm.vue'

@Component({
    components: {
        DbTable,
        CreateNewDatabaseForm
    }
})
export default class DashboardDatabaseList extends Vue {
    allDbs : Database[] = []
    showHideNew: boolean = false

    refreshDatabases() {
        allDatabase({
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllDatabaseOutputs) => {
            this.allDbs = resp.data
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops. Something went wrong. Try again.",
                false,
                "",
                contactUsUrl,
                true);
        })
    }

    mounted() {
        this.refreshDatabases()
    }

    onSaveDatabase(db : Database) {
        this.showHideNew = false
        this.allDbs.push(db)
    }
}

</script>
