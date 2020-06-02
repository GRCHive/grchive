<template>
    <div>
        <v-list-item>
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    Shell Script Runs
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
        </v-list-item>
        <v-divider></v-divider>

        <shell-run-table
            :resources="shellRuns"
            :search="filterText"
        >
        </shell-run-table>

    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import ShellRunTable from '../ShellRunTable.vue'
import {
    allShellRuns, TAllShellRunOutput, TAllShellRunInput
} from '../../../ts/api/apiShellRun'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'
import { ShellScriptRun } from '../../../ts/shell'

const Props = Vue.extend({
    props: {
        shellId: {
            type: Number,
            default: -1,
        },
        serverId: {
            type: Number,
            default: -1,
        },
    }
})

@Component({
    components: {
        ShellRunTable,
    }
})
export default class ShellRunTableWithControls extends Props {
    filterText: string = ""
    shellRuns: ShellScriptRun[] = []

    refreshData() {
        let input = <TAllShellRunInput>{
            orgId: PageParamsStore.state.organization!.Id,
        }

        if (this.shellId != -1) {
            input.shellId = this.shellId
        } else if (this.serverId != -1) {
            input.serverId = this.serverId
        }

        allShellRuns(input).then((resp : TAllShellRunOutput) => {
            this.shellRuns = resp.data
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
}

</script>
