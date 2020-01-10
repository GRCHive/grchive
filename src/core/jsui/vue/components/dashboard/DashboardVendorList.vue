<template>
    <div class="ma-4">
        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    Vendors
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

                    <create-new-vendor-form
                        @do-cancel="showHideNew = false"
                        @do-save="saveVendor">
                    </create-new-vendor-form>

                </v-dialog>
            </v-list-item-action>
        </v-list-item>
        <v-divider></v-divider>

        <vendor-table
            :resources="allVendors"
            :search="filterText"
        >
        </vendor-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'
import { Vendor } from '../../../ts/vendors'
import VendorTable from '../../generic/VendorTable.vue'
import CreateNewVendorForm from './CreateNewVendorForm.vue'
import { allVendors, TAllVendorOutput } from '../../../ts/api/apiVendors'

@Component({
    components: {
        VendorTable,
        CreateNewVendorForm
    }
})
export default class DashboardSystemList extends Vue {
    showHideNew: boolean = false
    filterText : string = ""
    allVendors : Vendor[] = []

    refreshVendors() {
        allVendors({
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllVendorOutput) => {
            this.allVendors = resp.data
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

    saveVendor(vendor : Vendor) {
        this.showHideNew = false
        this.allVendors.unshift(vendor)
    }

    mounted() {
        this.refreshVendors()
    }
}

</script>
