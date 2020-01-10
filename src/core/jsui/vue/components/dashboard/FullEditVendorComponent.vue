<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-list-item two-line class="pa-0">
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Vendor: {{ currentVendor.Name }}
                        <v-btn icon @click="expandDescription = !expandDescription">
                            <v-icon small v-if="!expandDescription" >mdi-chevron-down</v-icon>
                            <v-icon small v-else>mdi-chevron-up</v-icon>
                        </v-btn>
                    </v-list-item-title>

                    <v-list-item-subtitle :class="expandDescription ? `long-text` : `hide-long-text`">
                        {{ currentVendor.Description }}
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
                        item-name="vendors"
                        :items-to-delete="[currentVendor.Name]"
                        :use-global-deletion="false"
                        @do-cancel="showHideDelete = false"
                        @do-delete="onDelete">
                    </generic-delete-confirmation-form>
                </v-dialog>
            </v-list-item>
            <v-divider></v-divider>

            <v-tabs v-model="tab">
                <v-tab>Overview</v-tab>
                <v-tab>Products</v-tab>
                <v-tab>Documentation</v-tab>
            </v-tabs>

            <v-tabs-items v-model="tab">
                <v-tab-item>
                    <v-container fluid>
                        <create-new-vendor-form
                            edit-mode
                            :reference-vendor="currentVendor"
                            ref="editForm"
                            @do-save="onEdit"
                        >
                        </create-new-vendor-form>
                    </v-container>
                </v-tab-item>

                <v-tab-item>
                    <v-container fluid>
                    </v-container>
                </v-tab-item>

                <v-tab-item>
                    <v-container fluid>
                    </v-container>
                </v-tab-item>
            </v-tabs-items>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Vendor } from '../../../ts/vendors'
import { contactUsUrl, createOrgVendorsUrl } from '../../../ts/url'
import { deleteVendor, getVendor, TGetVendorOutput } from '../../../ts/api/apiVendors'
import { PageParamsStore } from '../../../ts/pageParams'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import CreateNewVendorForm from './CreateNewVendorForm.vue'

@Component({
    components: {
        GenericDeleteConfirmationForm,
        CreateNewVendorForm
    }
})
export default class FullEditVendorComponent extends Vue {
    currentVendor: Vendor | null = null
    ready : boolean = false
    expandDescription : boolean = false
    showHideDelete : boolean = false
    tab : number | null = 0

    $refs! : {
        editForm : CreateNewVendorForm
    }

    refreshData() {
        let data = window.location.pathname.split('/')
        let resourceId = Number(data[data.length - 1])

        getVendor({
            vendorId: resourceId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetVendorOutput) => {
            this.currentVendor = resp.data
            this.ready = true
            Vue.nextTick(() => {
                this.$refs.editForm.clearForm()
            })
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

    onDelete() {
        deleteVendor({
            vendorId: this.currentVendor!.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then(() => {
            window.location.replace(createOrgVendorsUrl(PageParamsStore.state.organization!.OktaGroupName))
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

    onEdit(v : Vendor) {
        this.currentVendor = v
        Vue.nextTick(() => {
            this.$refs.editForm.clearForm()
        })
    }
}

</script>
