<template>
    <div>
        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    Databases
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

            <v-list-item-action v-if="!disableNew">
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
            :resources="filteredDbs"
            :search="filterText"
            :use-crud-delete="!disableDelete"
            :confirm-delete="!disableDelete"
            @delete="deleteDb"
            @input="modifySelected"
            :selectable="enableSelect"
            :multi="enableSelect"
        ></db-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import DbTable from '../DbTable.vue'
import { Database, NullDatabaseFilterData } from '../../../ts/databases'
import { TAllDatabaseOutputs, allDatabase } from '../../../ts/api/apiDatabases'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'
import CreateNewDatabaseForm from '../../components/dashboard/CreateNewDatabaseForm.vue'
import { deleteDatabase } from '../../../ts/api/apiDatabases'
import { KNoHost } from '../../../ts/deployments'

const Props = Vue.extend({
    props: {
        value: {
            type: Array,
            default: () => [],
        },
        exclude: {
            type: Array,
            default: () => [],
        },
        disableNew: {
            type: Boolean,
            default: false,
        },
        disableDelete: {
            type: Boolean,
            default: false,
        },
        enableSelect: {
            type: Boolean,
            default: false,
        },
        deploymentType: {
            type: Number,
            default: KNoHost
        },
    }
})

@Component({
    components: {
        DbTable,
        CreateNewDatabaseForm
    }
})
export default class DbTableWithControls extends Props {
    allDbs : Database[] = []
    showHideNew: boolean = false
    filterText : string = ""

    get excludeSet() : Set<number> {
        return new Set<number>(this.exclude.map((ele : any) => ele.Id))
    }

    get filteredDbs() : Database[] {
        return this.allDbs.filter((ele : Database) => !this.excludeSet.has(ele.Id))
    }

    refreshDatabases() {
        let params : any = {
            orgId: PageParamsStore.state.organization!.Id,
            filter: NullDatabaseFilterData,
        }

        if (this.deploymentType != KNoHost) {
            params.deploymentType = this.deploymentType
        }

        allDatabase(params).then((resp : TAllDatabaseOutputs) => {
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
        this.allDbs.unshift(db)
    }

    modifySelected(vals : Database[]) {
        this.$emit('input', vals)
    }

    deleteDb(db : Database) {
        deleteDatabase({
            dbId: db.Id,
            orgId: PageParamsStore.state.organization!.Id
        }).then(() => {
            this.allDbs.splice(
                this.allDbs.findIndex((ele : Database) =>
                    ele.Id == db.Id),
                1)
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
}

</script>
