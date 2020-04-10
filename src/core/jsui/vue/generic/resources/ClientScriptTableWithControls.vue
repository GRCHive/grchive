<template>
    <div>
        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    Scripts
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

                    <create-new-script-form
                        @do-cancel="showHideNew = false"
                        @do-save="onNewScript"
                    >
                    </create-new-script-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>
        <v-divider></v-divider>

        <client-script-table
            :value="value"
            :resources="filteredScripts"
            :search="filterText"
            :use-crud-delete="!disableDelete"
            :confirm-delete="!disableDelete"
            @delete="deleteScript"
            @input="modifySelected"
            :selectable="enableSelect"
            :multi="enableSelect"
        ></client-script-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import ClientScriptTable from '../ClientScriptTable.vue'
import CreateNewScriptForm from '../../components/dashboard/CreateNewScriptForm.vue'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import { ClientScript } from '../../../ts/clientScripts'
import {
    allClientScripts, TAllClientScriptsOutput,
    deleteClientScript
} from '../../../ts/api/apiScripts'

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
        ClientScriptTable,
        CreateNewScriptForm,
    }
})
export default class ClientScriptTableWithControls extends Props {
    showHideNew: boolean = false
    filterText : string = ""
    data : ClientScript[] = []

    get excludeSet() : Set<number> {
        return new Set<number>(this.exclude.map((ele : any) => ele.Id))
    }

    get filteredScripts() : ClientScript[] {
        return this.data.filter((ele : ClientScript) => !this.excludeSet.has(ele.Id))
    }

    refreshData() {
        allClientScripts({
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllClientScriptsOutput) => {
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

    onNewScript(data : ClientScript) {
        this.showHideNew = false
        this.data.unshift(data)
    }

    mounted() {
        this.refreshData()
    }

    deleteScript(script : ClientScript) {
        deleteClientScript({
            orgId: PageParamsStore.state.organization!.Id,
            scriptId: script.Id,
        }).then(() => {
            let idx = this.data.findIndex((ele : ClientScript) => {
                return ele.Id == script.Id
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

    modifySelected(vals : ClientScript[]) {
        this.$emit('input', vals)
    }
}

</script>
