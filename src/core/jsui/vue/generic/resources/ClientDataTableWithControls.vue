<template>
    <div>
        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    Data Objects
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

                    <create-new-client-data-form
                        @do-cancel="showHideNew = false"
                        @do-save="onNewClientData"
                    >
                    </create-new-client-data-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>
        <v-divider></v-divider>

        <client-data-table
            :value="value"
            :resources="filteredData"
            :search="filterText"
            :use-crud-delete="!disableDelete"
            :confirm-delete="!disableDelete"
            @delete="deleteData"
            @input="modifySelected"
            :selectable="enableSelect"
            :multi="enableSelect"
        ></client-data-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import ClientDataTable from '../ClientDataTable.vue'
import { FullClientDataWithLink } from '../../../ts/clientData'
import { 
    allClientData, TAllClientDataOutput,
    newClientData, TNewClientDataOutput,
    deleteClientData
} from '../../../ts/api/apiClientData'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import CreateNewClientDataForm from '../../components/dashboard/CreateNewClientDataForm.vue'

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
        }
    }
})

@Component({
    components: {
        ClientDataTable,
        CreateNewClientDataForm
    }
})
export default class ClientDataTableWithControls extends Props {
    showHideNew: boolean = false
    filterText : string = ""
    data : FullClientDataWithLink[] = []

    get excludeSet() : Set<number> {
        return new Set<number>(this.exclude.map((ele : any) => ele.Data.Id))
    }

    get filteredData() : FullClientDataWithLink[] {
        return this.data.filter((ele : FullClientDataWithLink) => !this.excludeSet.has(ele.Data.Id))
    }

    refreshData() {
        allClientData({
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllClientDataOutput) => {
            this.data = resp.data
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

    onNewClientData(data : FullClientDataWithLink) {
        this.showHideNew = false
        this.data.unshift(data)
    }

    mounted() {
        this.refreshData()
    }

    deleteData(data : FullClientDataWithLink) {
        deleteClientData({
            orgId: PageParamsStore.state.organization!.Id,
            dataId: data.Data.Id,
        }).then(() => {
            let idx = this.data.findIndex((ele : FullClientDataWithLink) => {
                return ele.Data.Id == data.Data.Id
            })
            if (idx == -1) {
                return
            }
            this.data.splice(idx, 1)
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

    modifySelected(vals : FullClientDataWithLink[]) {
        this.$emit('input', vals)
    }
}

</script>
