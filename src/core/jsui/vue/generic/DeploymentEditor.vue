<template>
    <div>
        <v-select
            v-model="editableDeployment.DeploymentType"
            filled
            label="Type"
            :items="typeItems"
            hide-details
            :disabled="!canEdit"
        >
        </v-select>

        <v-divider class="my-2"></v-divider>

        <!-- self hosted -->
        <div v-if="deploymentType == 0">
            <v-form
                v-model="formValid"
            >
            </v-form>

            <v-row class="title ml-2">
                Servers

                <v-spacer></v-spacer>
                
                <v-dialog
                    v-model="showHideLinkServer"
                    persistent
                    max-width="40%"
                >
                    <template v-slot:activator="{ on }">
                        <v-btn
                            color="primary"
                            icon
                            :disabled="!canEdit"
                            v-on="on"
                        >
                            <v-icon>mdi-plus</v-icon>
                        </v-btn>
                    </template>

                    <v-card>
                        <v-card-title>
                            Link Servers
                        </v-card-title>
                        <v-divider></v-divider>

                        <server-table :resources="linkableServers"
                                       selectable
                                       multi
                                       v-model="serversToLink"
                        ></server-table>

                        <v-card-actions>
                            <v-btn color="error" @click="showHideLinkServer = false">
                                Cancel
                            </v-btn>
                            <v-spacer></v-spacer>
                            <v-btn color="success" @click="linkServers">
                                Link
                            </v-btn>
                        </v-card-actions>

                    </v-card>
                </v-dialog>
            </v-row>

            <server-table
                :resources="editableDeployment.SelfDeployment.Servers"
                :use-crud-delete="canEdit"
                confirm-delete
                @delete="unlinkServer"
            >
            </server-table>

            <v-divider></v-divider>
        </div>

        <!-- vendor hosted -->
        <div v-if="deploymentType == 1">
            <v-form
                v-model="formValid"
            >
                <vendor-product-search-form-component
                    v-model="editableDeployment.VendorDeployment.Product"
                    :rules="[rules.required]"
                    :disabled="!canEdit"
                    :initial-vendor-id="initialVendorId"
                >
                </vendor-product-search-form-component>
            </v-form>

        </div>

        <v-list-item class="pa-0">
            <v-btn
                color="error"
                @click="cancel"
                v-if="canEdit"
            >
                Cancel
            </v-btn>
            <v-spacer></v-spacer>
            <v-btn
                color="success"
                @click="save"
                :disabled="!canSubmit"
                v-if="canEdit"
            >
                Save
            </v-btn>

            <v-btn
                color="primary"
                @click="canEdit = true"
                v-if="!canEdit"
            >
                Edit
            </v-btn>
        </v-list-item>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { FullDeployment,    
         SelfDeployment,
         VendorDeployment,
         createEmptyVendorDeployment,
         createEmptySelfDeployment,
         deepCopyFullDeployment,
         KSelfHosted,
         KVendorHosted,
         KNoHost
} from '../../ts/deployments'
import {
    TUpdateDeploymentOutput,
    updateDeployment
} from '../../ts/api/apiDeployments'
import {
    TGetAllDocumentRequestOutput,
    getAllDocRequests
} from '../../ts/api/apiDocRequests'
import DocFileTable from './DocFileTable.vue'
import UploadDocumentationForm from '../components/dashboard/UploadDocumentationForm.vue'
import DocSearcherForm from './DocSearcherForm.vue'
import CreateNewRequestForm from '../components/dashboard/CreateNewRequestForm.vue'
import { ControlDocumentationFile } from '../../ts/controls'
import { DocumentRequest } from '../../ts/docRequests'
import { contactUsUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'
import DocRequestTable from './DocRequestTable.vue'
import ServerTable from './ServerTable.vue'
import { Server } from '../../ts/infrastructure'
import { allServers, TAllServerOutput } from '../../ts/api/apiServers'
import VendorProductSearchFormComponent from './VendorProductSearchFormComponent.vue'
import * as rules from '../../ts/formRules'

const VueProps = Vue.extend({
    props: {
        value: Object as () => FullDeployment,
    }
})

@Component({
    components: {
        DocFileTable,
        UploadDocumentationForm,
        DocSearcherForm,
        CreateNewRequestForm,
        DocRequestTable,
        ServerTable,
        VendorProductSearchFormComponent
    }
})
export default class DeploymentEditor extends VueProps {
    canEdit: boolean = false
    formValid: boolean = false
    rules : any = rules
    editableDeployment: FullDeployment = {} as FullDeployment

    showHideLinkServer : boolean = false
    allAvailableServers: Server[] = []
    serversToLink : Server[] = []

    get initialVendorId() : number {
        if (!this.editableDeployment.VendorDeployment) {
            return -1
        }

        if (!this.editableDeployment.VendorDeployment.Product) {
            return -1
        }
        return this.editableDeployment.VendorDeployment.Product.VendorId
    }

    get linkableServers() : Server[] {
        let usedServerIds = new Set<number>()
        if (this.editableDeployment.SelfDeployment) {
            for (let s of this.editableDeployment.SelfDeployment.Servers) {
                usedServerIds.add(s.Id)
            }
        }
        return this.allAvailableServers.filter((ele : Server) => !usedServerIds.has(ele.Id))
    }

    get deploymentType() : number {
        return this.editableDeployment.DeploymentType
    }

    get canSubmit() : boolean {
        if (this.deploymentType == KSelfHosted || this.deploymentType == KVendorHosted) {
            return this.formValid
        }
        return true
    }

    get typeItems() : any[] {
        return [
            {
                text: 'Self Hosted',
                value: KSelfHosted
            },
            {
                text: 'Vendor Hosted',
                value: KVendorHosted
            },
            {
                text: 'N/A',
                value: KNoHost,
            },
        ]
    }

    resetEditCopyFromProps() {
        this.editableDeployment = deepCopyFullDeployment(this.value)
    }

    reloadServers() {
        allServers({
            orgId: PageParamsStore.state.organization!.Id
        }).then((resp : TAllServerOutput) => {
            this.allAvailableServers = resp.data
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
        this.resetEditCopyFromProps()
        this.reloadServers()
    }

    cancel() {
        this.canEdit = false
        this.resetEditCopyFromProps()
    }

    save() {
        updateDeployment({
            deployment: this.editableDeployment
        }).then((resp : TUpdateDeploymentOutput) => {
            this.$emit('input', resp.data)
            Vue.nextTick(this.cancel)
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

    unlinkServer(val : Server) {
        this.editableDeployment.SelfDeployment!.Servers.splice(
            this.editableDeployment.SelfDeployment!.Servers.findIndex(
                (ele : Server) => ele.Id == val.Id
            ),
            1
        )
    }

    linkServers() {
        this.editableDeployment.SelfDeployment!.Servers.push(...this.serversToLink)
        this.showHideLinkServer = false
        this.serversToLink = []
    }
}

</script>
