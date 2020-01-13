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
                    <create-new-vendor-form
                        edit-mode
                        :reference-vendor="currentVendor"
                        ref="editForm"
                        @do-save="onEdit"
                    >
                    </create-new-vendor-form>
                </v-tab-item>

                <v-tab-item>
                    <v-container fluid>
                        <v-row>
                            <v-col cols="2">
                                <v-list>
                                    <v-subheader>PRODUCTS</v-subheader>

                                    <v-dialog
                                        v-model="showHideNewProduct"
                                        persistent
                                        max-width="40%"
                                    >
                                        <template v-slot:activator="{ on }">
                                            <v-list-item v-on="on">
                                                <v-list-item-icon>
                                                    <v-icon>mdi-plus</v-icon>
                                                </v-list-item-icon>

                                                <v-list-item-content>
                                                    New Product
                                                </v-list-item-content>
                                            </v-list-item>
                                        </template>

                                        <create-new-vendor-product-form
                                            :parent-vendor="currentVendor"
                                            @do-cancel="showHideNewProduct = false"
                                            @do-save="onSaveProduct"
                                        >
                                        </create-new-vendor-product-form>
                                    </v-dialog>

                                    <v-list-item-group
                                        v-model="selectedProductIndex"
                                        mandatory
                                    >
                                        <v-list-item
                                            v-for="(product, i) in allProducts"
                                            :key="i"
                                        >
                                            <v-list-item-content>
                                                <v-list-item-title>
                                                    {{ product.Name }}
                                                </v-list-item-title>
                                            </v-list-item-content>
                                        </v-list-item>

                                    </v-list-item-group>
                                </v-list>
                            </v-col>

                            <v-col cols="10" v-if="!!selectedProduct">
                                <div v-if="!productFullyLoaded" class="max-height">
                                    <v-row justify="center" align="center" class="max-height">
                                        <v-progress-circular indeterminate size="64"></v-progress-circular>
                                    </v-row>
                                </div>

                                <div v-else>
                                    <v-row>
                                        <v-col cols="4">
                                            <create-new-vendor-product-form
                                                edit-mode
                                                :reference-product="selectedProduct"
                                                :parent-vendor="currentVendor"
                                                @do-save="onEditProduct"
                                            >
                                            </create-new-vendor-product-form>
                                        </v-col>

                                        <v-col cols="8">
                                            <v-list-item class="title">
                                                SOC Reports
                                            </v-list-item>

                                            <doc-file-manager
                                                v-model="productSocFiles"
                                                :cat-id="currentVendor.DocCatId"
                                                @new-doc="onAddNewSocFile"
                                            >
                                                <template v-slot:multiActions="{ hasSelected, selectedFiles }">
                                                    <v-btn 
                                                        color="warning"
                                                        :disabled="!hasSelected"
                                                        @click="unlinkSocFiles(selectedFiles)"
                                                    >
                                                        Unlink
                                                    </v-btn>
                                                </template>

                                                <template v-slot:singleActions>
                                                    <v-dialog
                                                        v-model="showAddSoc"
                                                        persistent
                                                        max-width="40%"
                                                    >
                                                        <template v-slot:activator="{ on }">
                                                            <v-btn color="info" v-on="on">
                                                                Add Existing
                                                            </v-btn>
                                                        </template>

                                                        <doc-searcher-form
                                                            :force-cat-id="currentVendor.DocCatId"
                                                            :exclude-files="productSocFiles"
                                                            @do-select="onSelectSOC"
                                                            @do-cancel="showAddSoc = false">
                                                        </doc-searcher-form>
                                                    </v-dialog>

                                                </template>
                                            </doc-file-manager>

                                            <v-divider class="my-3"></v-divider>

                                            <v-list-item class="title">
                                                SOC Report Requests

                                                <v-spacer></v-spacer>

                                                <v-dialog v-model="showRequestSoc" persistent max-width="40%">
                                                    <template v-slot:activator="{ on }">
                                                        <v-btn color="warning" icon v-on="on">
                                                            <v-icon>mdi-plus</v-icon>
                                                        </v-btn>
                                                    </template>

                                                    <create-new-request-form
                                                        :cat-id="currentVendor.DocCatId"
                                                        :vendor-product-id="selectedProduct.Id"
                                                        @do-cancel="showRequestSoc = false"
                                                        @do-save="onRequestSOC">
                                                    </create-new-request-form>
                                                </v-dialog>
                                            </v-list-item>

                                            <doc-request-table :resources="productSocRequests">
                                            </doc-request-table>

                                        </v-col>
                                    </v-row>
                                </div>
                            </v-col>
                        </v-row>
                    </v-container>
                </v-tab-item>

                <v-tab-item>
                    <full-edit-documentation-category-component
                        content-only
                        :resource-id="currentVendor.DocCatId"
                        :key="refreshDocCat"
                    >
                    </full-edit-documentation-category-component>
                </v-tab-item>
            </v-tabs-items>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { Vendor, VendorProduct } from '../../../ts/vendors'
import { contactUsUrl, createOrgVendorsUrl } from '../../../ts/url'
import { deleteVendor, getVendor, TGetVendorOutput } from '../../../ts/api/apiVendors'
import { 
    allVendorProducts, TAllVendorProductOutput,
    getVendorProduct, TGetVendorProductOutput,
    linkVendorProductSocFiles, unlinkVendorProductSocFiles
} from '../../../ts/api/apiVendorProduct'
import { PageParamsStore } from '../../../ts/pageParams'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import CreateNewVendorForm from './CreateNewVendorForm.vue'
import FullEditDocumentationCategoryComponent from './FullEditDocumentationCategoryComponent.vue'
import CreateNewVendorProductForm from './CreateNewVendorProductForm.vue'
import { ControlDocumentationFile, extractControlDocumentationFileHandle } from '../../../ts/controls'
import { DocumentRequest } from '../../../ts/docRequests'
import {
    TGetAllDocumentRequestOutput,
    getAllDocRequests
} from '../../../ts/api/apiDocRequests'
import DocFileManager from '../../generic/DocFileManager.vue'
import DocRequestTable from '../../generic/DocRequestTable.vue'
import DocSearcherForm from '../../generic/DocSearcherForm.vue'
import CreateNewRequestForm from './CreateNewRequestForm.vue'

@Component({
    components: {
        GenericDeleteConfirmationForm,
        CreateNewVendorForm,
        FullEditDocumentationCategoryComponent,
        CreateNewVendorProductForm,
        DocFileManager,
        DocRequestTable,
        DocSearcherForm,
        CreateNewRequestForm
    }
})
export default class FullEditVendorComponent extends Vue {
    currentVendor: Vendor | null = null
    expandDescription : boolean = false
    showHideDelete : boolean = false
    tab : number | null = 0

    selectedProductIndex : number | null = null
    allProducts : VendorProduct[] = []
    loadedProducts: boolean = false
    showHideNewProduct: boolean = false

    productFullyLoaded : boolean = false
    productSocFiles : ControlDocumentationFile[] = []
    productSocRequests : DocumentRequest[] = []

    showAddSoc: boolean = false
    showRequestSoc: boolean = false

    // Ideally there'd be some shared ControlDocumentationFile[] local repository
    // but there isn't so we need to force the FullEditDocumentationCategoryComponent
    // component to regrab all available files in case we do an upload in the product
    // tab.
    refreshDocCat : number = 0

    $refs! : {
        editForm : CreateNewVendorForm
    }

    get selectedProduct() : VendorProduct | null {
        if (this.selectedProductIndex == null) {
            return null
        }

        if (this.selectedProductIndex < 0 || this.selectedProductIndex >= this.allProducts.length) {
            return null
        }

        return this.allProducts[this.selectedProductIndex]
    }

    get ready() : boolean {
        return !!this.currentVendor && this.loadedProducts
    }

    refreshData() {
        let data = window.location.pathname.split('/')
        let resourceId = Number(data[data.length - 1])

        getVendor({
            vendorId: resourceId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetVendorOutput) => {
            this.currentVendor = resp.data
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })

        allVendorProducts({
            orgId: PageParamsStore.state.organization!.Id,
            vendorId: resourceId,
        }).then((resp : TAllVendorProductOutput) => {
            this.allProducts = resp.data
            this.loadedProducts = true
        }).catch((err: any) => {
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

    onSaveProduct(p : VendorProduct) {
        this.allProducts.unshift(p)
        this.selectedProductIndex = 0
        this.showHideNewProduct = false
    }

    onEditProduct(p : VendorProduct) {
        this.allProducts[this.selectedProductIndex!].Name = p.Name
        this.allProducts[this.selectedProductIndex!].Description = p.Description
        this.allProducts[this.selectedProductIndex!].Url = p.Url
    }

    @Watch('selectedProduct')
    onChangeProduct() {
        this.productFullyLoaded = false
        if (!this.selectedProduct) {
            return
        }

        getVendorProduct({
            productId: this.selectedProduct!.Id,
            vendorId: this.currentVendor!.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetVendorProductOutput) => {
            this.productSocFiles = resp.data.SocFiles
            this.productFullyLoaded = true
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })

        getAllDocRequests({
            orgId: PageParamsStore.state.organization!.Id,
            catId: this.currentVendor!.DocCatId,
            vendorProductId: this.selectedProduct!.Id,
        }).then((resp : TGetAllDocumentRequestOutput) => {
            this.productSocRequests = resp.data
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

    onAddNewSocFile(f : ControlDocumentationFile) {
        this.onSelectSOC([f], false)
        this.refreshDocCat += 1
    }

    onSelectSOC(f : ControlDocumentationFile[], needAddToModel : boolean = true) {
        linkVendorProductSocFiles({
            productId: this.selectedProduct!.Id,
            vendorId: this.currentVendor!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            socFiles: f.map(extractControlDocumentationFileHandle)
        }).then(() => {
            if (needAddToModel) {
                this.productSocFiles.unshift(...f)
            }
            this.showAddSoc = false
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

    unlinkSocFiles(files : ControlDocumentationFile[]) {
        unlinkVendorProductSocFiles({
            productId: this.selectedProduct!.Id,
            vendorId: this.currentVendor!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            socFiles: files.map(extractControlDocumentationFileHandle)
        }).then(() => {
            let idSet = new Set<number>(files.map((ele : ControlDocumentationFile) => ele.Id))
            this.productSocFiles = this.productSocFiles.filter((ele : ControlDocumentationFile) => !idSet.has(ele.Id))
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

    onRequestSOC(req : DocumentRequest) {
        this.productSocRequests.push(req)
        this.showRequestSoc = false
    }
}

</script>
