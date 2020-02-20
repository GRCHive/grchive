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
                    <v-col cols="6">
                        <create-new-system-form
                            ref="editForm"
                            :edit-mode="true"
                            :reference-system="currentSystem"
                            @do-save="onEdit">
                        </create-new-system-form>
                    </v-col>

                    <v-col cols="6">
                        <v-card class="mb-4">
                            <v-card-title>
                                Related Resources
                            </v-card-title>
                            <v-divider></v-divider>

                            <v-tabs>
                                <v-tab>Databases</v-tab>
                                <v-tab-item>
                                    <v-list-item>
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
                                    </v-list-item>

                                    <db-table
                                        :resources="relatedDbs"
                                        use-crud-delete
                                        @delete="onDeleteDbLink"
                                    ></db-table>
                                </v-tab-item>

                            <v-tab>Process Flows</v-tab>
                            <v-tab-item>
                                <process-flow-table
                                    :resources="relatedFlows"
                                >
                                </process-flow-table>
                            </v-tab-item>

                            <v-tab>Risks</v-tab>
                            <v-tab-item>
                                <risk-table
                                    :resources="relatedRisks"
                                >
                                </risk-table>
                            </v-tab-item>

                            <v-tab>Controls</v-tab>
                            <v-tab-item>
                                <control-table
                                    :resources="relatedControls"
                                >
                                </control-table>
                            </v-tab-item>
                            </v-tabs>
                        </v-card>

                        <v-card>
                            <v-card-title>
                                Deployment
                            </v-card-title>
                            <v-divider></v-divider>

                            <div v-if="!deployment">
                                <v-row align="center" justify="center">
                                    <v-btn color="primary"
                                           fab
                                           outlined
                                           x-large
                                           class="my-6"
                                           @click="addDeployment">
                                        <v-icon>mdi-plus</v-icon>
                                    </v-btn>
                                </v-row>
                            </div>

                            <div v-else>
                                <deployment-editor
                                    class="pa-4"
                                    v-model="deployment"
                                    :system-id="currentSystem.Id"
                                >
                                </deployment-editor>
                            </div>

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
import { linkDatabasesToSystem, deleteDbSysLink } from '../../../ts/api/apiSystems'
import { PageParamsStore } from '../../../ts/pageParams'
import { System } from '../../../ts/systems'
import CreateNewSystemForm from './CreateNewSystemForm.vue'
import { contactUsUrl, createOrgSystemUrl } from '../../../ts/url'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import DbTable from '../../generic/DbTable.vue'
import { Database } from '../../../ts/databases'
import { FullDeployment } from '../../../ts/deployments'
import DeploymentEditor from '../../generic/DeploymentEditor.vue'
import { newDeployment, TNewDeploymentOutput } from '../../../ts/api/apiDeployments'
import { TAllRiskSystemLinkOutput, allRiskSystemLink } from '../../../ts/api/apiRiskSystemLinks'
import { TAllControlSystemLinkOutput, allControlSystemLink } from '../../../ts/api/apiControlSystemLinks'
import { TAllNodeSystemLinkOutput, allNodeSystemLink } from '../../../ts/api/apiNodeSystemLinks'
import RiskTable from '../../generic/RiskTable.vue'
import ControlTable from '../../generic/ControlTable.vue'
import ProcessFlowTable from '../../generic/ProcessFlowTable.vue'

@Component({
    components: {
        CreateNewSystemForm,
        GenericDeleteConfirmationForm,
        DbTable,
        DeploymentEditor,
        RiskTable,
        ControlTable,
        ProcessFlowTable
    }
})
export default class FullEditSystemComponent extends Vue {
    currentSystem: System = {} as System
    relatedDbs: Database[] = []
    allDb: Database[] = []

    dbToLink: Database[] = []
    deployment: FullDeployment | null = null

    relatedControls : ProcessFlowControl[] = []
    relatedRisks : ProcessFlowRisk[] = []
    relatedFlows : ProcessFlowBasicData[] = []

    ready : boolean = false
    expandDescription: boolean = false
    showHideDelete: boolean = false
    showHideLinkDb: boolean = false

    $refs!: {
        editForm: CreateNewSystemForm
    }

    addDeployment() {
        newDeployment({
            orgId: PageParamsStore.state.organization!.Id,
            systemId: this.currentSystem.Id,
            dbId: null,
        }).then((resp : TNewDeploymentOutput) => {
            this.deployment = resp.data
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

    refreshRelatedFlows() {
        allNodeSystemLink({
            systemId: this.currentSystem.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllNodeSystemLinkOutput) => {
            this.relatedFlows = <ProcessFlowBasicData[]>resp.data
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

    refreshRelatedRisks() {
        allRiskSystemLink({
            systemId: this.currentSystem.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllRiskSystemLinkOutput) => {
            this.relatedRisks = resp.data.Risks!
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

    refreshRelatedControls() {
        allControlSystemLink({
            systemId: this.currentSystem.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllControlSystemLinkOutput) => {
            this.relatedControls = resp.data.Controls!
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
            this.deployment = resp.data.Deployment

            this.refreshRelatedFlows()
            this.refreshRelatedRisks()
            this.refreshRelatedControls()

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

    onDeleteDbLink(db : Database) {
        deleteDbSysLink({
            sysId: this.currentSystem.Id,
            orgId: PageParamsStore.state.organization!.Id,
            dbId: db.Id,
        }).then(() => {
            this.relatedDbs = this.relatedDbs.filter((ele : Database) => ele.Id != db.Id)
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
