<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-list-item two-line class="pa-0">
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Database: {{ currentDb.Name }}
                    </v-list-item-title>

                    <v-list-item-subtitle>
                        {{ fullTypeString }}
                    </v-list-item-subtitle>
                </v-list-item-content>

                <v-dialog v-model="showHideDelete"
                          persistent
                          max-width="40%"
                >
                    <template v-slot:activator="{ on }">
                        <v-btn color="error" v-on="on">
                            Delete
                        </v-btn>
                    </template>

                    <generic-delete-confirmation-form
                        item-name="databases"
                        :items-to-delete="[currentDb.Name]"
                        :use-global-deletion="false"
                        @do-cancel="showHideDelete = false"
                        @do-delete="onDelete">
                    </generic-delete-confirmation-form>
                </v-dialog>

            </v-list-item>
            <v-divider></v-divider>

            <v-container fluid>
                <create-new-database-form
                    ref="editForm"
                    :edit-mode="true"
                    :reference-db="currentDb"
                    @do-save="onEdit">
                </create-new-database-form>
            </v-container>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { getDatabase, TGetDatabaseOutputs } from '../../../ts/api/apiDatabases'
import { deleteDatabase, TDeleteDatabaseOutputs } from '../../../ts/api/apiDatabases'
import { PageParamsStore } from '../../../ts/pageParams'
import { Database, getDbTypeAsString } from '../../../ts/databases'
import CreateNewDatabaseForm from './CreateNewDatabaseForm.vue'
import { contactUsUrl, createOrgDatabaseUrl } from '../../../ts/url'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'

@Component({
    components: {
        CreateNewDatabaseForm,
        GenericDeleteConfirmationForm
    }
})
export default class FullEditDatabaseComponent extends Vue {
    currentDb: Database = {} as Database
    ready : boolean = false
    showHideDelete: boolean = false

    $refs!: {
        editForm: CreateNewDatabaseForm
    }

    get fullTypeString() : string {
        return `${getDbTypeAsString(this.currentDb)} ${this.currentDb.Version}`
    }

    refreshDbData() {
        let data = window.location.pathname.split('/')
        let dbId = Number(data[data.length - 1])

        getDatabase({
            dbId: dbId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetDatabaseOutputs) => {
            this.currentDb = resp.data.Database
            this.ready = true

            Vue.nextTick(() => {
                this.$refs.editForm.clearForm()
            })
        }).catch((err : any) => {
            window.location.replace('/404')
        })
    }

    mounted() {
        this.refreshDbData()
    }

    onDelete() {
        deleteDatabase({
            dbId: this.currentDb.Id,
            orgId: PageParamsStore.state.organization!.Id
        }).then((resp : TDeleteDatabaseOutputs) => {
            window.location.replace(createOrgDatabaseUrl(PageParamsStore.state.organization!.OktaGroupName))
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

    onEdit(db : Database) {
        this.currentDb.Name = db.Name
        this.currentDb.TypeId = db.TypeId
        this.currentDb.OtherType = db.OtherType
        this.currentDb.Version = db.Version
    }
}

</script>
