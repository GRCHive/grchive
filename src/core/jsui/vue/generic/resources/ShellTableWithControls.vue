<template>
    <div>
        <v-list-item>
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    Shell Scripts
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
                          max-width="80%"
                >
                    <template v-slot:activator="{ on }">
                        <v-btn color="primary" v-on="on">
                            New
                        </v-btn>
                    </template>

                    <create-new-shell-script-form
                        :shell-type="shellType"
                        @do-cancel="showHideNew = false"
                        @do-save="onNewShellScript"
                    >
                    </create-new-shell-script-form>
                </v-dialog>
            </v-list-item-action>

        </v-list-item>
        <v-divider></v-divider>

        <shell-table
            :resources="scripts"
            :search="filterText"
            use-crud-delete
            confirm-delete
            @delete="deleteShellScript"
        >
        </shell-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import CreateNewShellScriptForm from '../../components/dashboard/CreateNewShellScriptForm.vue'
import ShellTable from '../ShellTable.vue'
import { ShellTypes, ShellScript } from '../../../ts/shell'
import {
    allShellScripts, TAllShellScriptsOutput,
    deleteShellScript
} from '../../../ts/api/apiShell'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'

const Props = Vue.extend({
    props: {
        shellType: Number
    }
})

@Component({
    components: {
        CreateNewShellScriptForm,
        ShellTable,
    }
})
export default class ShellTableWithControls extends Props {
    showHideNew : boolean = false
    filterText: string = ""
    scripts : ShellScript[] = []

    refreshData() {
        allShellScripts({
            orgId: PageParamsStore.state.organization!.Id,
            shellType: this.shellType,
        }).then((resp : TAllShellScriptsOutput) => {
            this.scripts = resp.data
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

    onNewShellScript(s : ShellScript) {
        this.scripts.unshift(s)
        this.showHideNew = false
    }

    deleteShellScript(s : ShellScript) {
        deleteShellScript({
            orgId: PageParamsStore.state.organization!.Id,
            shellId: s.Id,
        }).then(() => {
            let idx = this.scripts.findIndex((ele : ShellScript) => ele.Id == s.Id)
            if (idx == -1) {
                return
            }
            this.scripts.splice(idx, 1)
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
