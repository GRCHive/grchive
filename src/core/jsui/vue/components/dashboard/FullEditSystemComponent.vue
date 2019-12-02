<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-list-item two-line class="pa-0">
                <v-list-item-content>
                    <v-list-item-title class="title">
                        System: {{ currentSystem.Name }}
                        <v-btn icon @click="expandDescription = !expandDescription">
                            <v-icon small v-if="!expandDescription" >mdi-chevron-down</v-icon>
                            <v-icon small v-else>mdi-chevron-up</v-icon>
                        </v-btn>
                    </v-list-item-title>

                    <v-list-item-subtitle :class="expandDescription ? `long-text` : `hide-long-text`">
                        {{ currentSystem.Description }}
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
                        item-name="systems"
                        :items-to-delete="[currentSystem.Name]"
                        :use-global-deletion="false"
                        @do-cancel="showHideDelete = false"
                        @do-delete="onDelete">
                    </generic-delete-confirmation-form>
                </v-dialog>

            </v-list-item>
            <v-divider></v-divider>

            <v-container fluid>
                <v-row>
                    <v-col cols="8">
                        <create-new-system-form
                            ref="editForm"
                            :edit-mode="true"
                            :reference-system="currentSystem"
                            @do-save="onEdit">
                        </create-new-system-form>
                    </v-col>

                    <v-col cols="4">
                        <v-card class="mb-4">
                            <v-card-title>
                                Relevant Databases
                                <v-spacer></v-spacer>
                                
                                <v-dialog persistent
                                          max-width="40%"
                                          v-model="showHideLinkDb">
                                    <template v-slot:activator="{ on }">
                                        <v-btn color="primary" icon v-on="on">
                                            <v-icon>mdi-plus</v-icon>
                                        </v-btn>
                                    </template>

                                    <v-card>
                                        <v-card-title>
                                            Link Databases
                                        </v-card-title>
                                        <v-divider></v-divider>

                                        <db-table :resources="allDb"
                                                  v-model="dbToLink"
                                                  selectable
                                                  multi
                                        ></db-table>

                                        <v-card-actions>
                                            <v-btn color="error" @click="showHideLinkDb = false">
                                                Cancel
                                            </v-btn>
                                            <v-spacer></v-spacer>
                                            <v-btn color="success" @click="linkDbs">
                                                Link
                                            </v-btn>
                                        </v-card-actions>
                                    </v-card>
                                </v-dialog>
                            </v-card-title>
                            <v-divider></v-divider>
                            <db-table :resources="relatedDbs"></db-table>
                        </v-card>

                        <v-card>
                            <v-card-title>
                                Infrastructure
                            </v-card-title>
                            <v-divider></v-divider>

                            <v-row align="center" justify="center">
                                <v-btn color="primary"
                                       fab
                                       outlined
                                       x-large
                                       class="my-6">
                                    <v-icon>mdi-plus</v-icon>
                                </v-btn>
                            </v-row>
                        </v-card>
                    </v-col>
                </v-row>
            </v-container>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { getSystem, TGetSystemOutputs } from '../../../ts/api/apiSystems'
import { deleteSystem, TDeleteSystemOutputs } from '../../../ts/api/apiSystems'
import { linkDatabasesToSystem } from '../../../ts/api/apiSystems'
import { PageParamsStore } from '../../../ts/pageParams'
import { System } from '../../../ts/systems'
import CreateNewSystemForm from './CreateNewSystemForm.vue'
import { contactUsUrl, createOrgSystemUrl } from '../../../ts/url'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import DbTable from '../../generic/DbTable.vue'
import { Database } from '../../../ts/databases'

@Component({
    components: {
        CreateNewSystemForm,
        GenericDeleteConfirmationForm,
        DbTable
    }
})
export default class FullEditSystemComponent extends Vue {
    currentSystem: System = {} as System
    relatedDbs: Database[] = []
    allDb: Database[] = []

    dbToLink: Database[] = []

    ready : boolean = false
    expandDescription: boolean = false
    showHideDelete: boolean = false
    showHideLinkDb: boolean = false

    $refs!: {
        editForm: CreateNewSystemForm
    }

    refreshSystemData() {
        let data = window.location.pathname.split('/')
        let systemId = Number(data[data.length - 1])

        getSystem({
            sysId: systemId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetSystemOutputs) => {
            this.currentSystem = resp.data.System
            this.ready = true
            this.allDb = resp.data.AllDatabases

            let idSet = new Set(resp.data.RelevantDatabaseIds)
            this.relatedDbs = resp.data.AllDatabases.filter((ele : Database) => idSet.has(ele.Id))

            Vue.nextTick(() => {
                this.$refs.editForm.clearForm()
            })
        }).catch((err : any) => {
            window.location.replace('/404')
        })
    }

    mounted() {
        this.refreshSystemData()
    }

    onDelete() {
        deleteSystem({
            sysId: this.currentSystem.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TDeleteSystemOutputs) => {
            window.location.replace(createOrgSystemUrl(PageParamsStore.state.organization!.OktaGroupName))
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

    onEdit(sys : System) {
        this.currentSystem.Name = sys.Name
        this.currentSystem.Purpose = sys.Purpose
        this.currentSystem.Description = sys.Description
    }

    linkDbs() {
        linkDatabasesToSystem({
            sysId: this.currentSystem.Id,
            orgId: PageParamsStore.state.organization!.Id,
            dbIds: this.dbToLink.map((ele : Database) => ele.Id)
        }).then(() => {
            let idSet = new Set([...this.dbToLink, ...this.relatedDbs].map((ele : Database) => ele.Id))
            this.relatedDbs = this.allDb.filter((ele : Database) => idSet.has(ele.Id))
            this.dbToLink = []
            this.showHideLinkDb = false
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
