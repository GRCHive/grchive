<template>
    <div class="max-height">
        <dashboard-app-bar ref="dashboardAppBar">
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar></dashboard-home-page-nav-bar>

        <v-overlay :value="isLoading">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <v-content v-if="!isLoading">
            <v-list-item two-line>
                <v-list-item-content>
                    <v-list-item-title>
                        <a :href="url">
                            <span class="font-weight-bold">
                                Script: 
                            </span>
                            {{ script.Name }}
                        </a>
                    </v-list-item-title>

                    <v-list-item-subtitle>
                        {{ code.GitHash }}
                    </v-list-item-subtitle>
                </v-list-item-content>

                <v-list-item-content>
                    <v-list-item-title>
                        <span class="font-weight-bold">
                            Run By: 
                        </span>
                        {{  runByUser }}
                    </v-list-item-title>

                    <v-list-item-subtitle>
                        <span class="font-weight-bold">
                            Run At: 
                        </span>
                        {{ runTimeStr }}
                    </v-list-item-subtitle>
                </v-list-item-content>

                <v-list-item-content>
                    <span class="font-weight-bold overline item-content-no-flex">
                        Build Status: 
                    </span>
                    <script-build-run-status
                        class="item-content-no-flex"
                        :success="run.BuildSuccess"
                        :start="run.BuildStartTime"
                        :end="run.BuildFinishTime"
                        show-time-stamp
                    >
                    </script-build-run-status>
                </v-list-item-content>

                <v-list-item-content>
                    <span class="font-weight-bold overline item-content-no-flex">
                        Run Status: 
                    </span>
                    <script-build-run-status
                        class="item-content-no-flex"
                        :success="run.RunSuccess"
                        :start="run.RunStartTime"
                        :end="run.RunFinishTime"
                        show-time-stamp
                    >
                    </script-build-run-status>
                </v-list-item-content>
            </v-list-item>
            <v-divider></v-divider>
            <v-row>
                <v-col cols="6" class="pb-0">
                    <p class="overline text-center">Run Log</p>
                    <log-viewer :run-id="runId" run-log></log-viewer>
                </v-col>

                <v-col cols="6" class="pb-0">
                    <p class="overline text-center">Build Log</p>
                    <log-viewer :run-id="runId"></log-viewer>
                </v-col>
            </v-row>
        </v-content>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import LogViewer from '../../generic/logs/LogViewer.vue'
import { PageParamsStore } from '../../../ts/pageParams'
import { ClientScript } from '../../../ts/clientScripts'
import { ManagedCode, ScriptRun } from '../../../ts/code'
import { getClientScriptCodeFromLink, TGetClientScriptCodeFromLinkOutput } from '../../../ts/api/apiScripts'
import { getCodeRun, TGetCodeRunOutput } from '../../../ts/api/apiCode'
import {
    contactUsUrl,
    createSingleScriptUrl,
} from '../../../ts/url'
import { standardFormatTime } from '../../../ts/time'
import { createUserString } from '../../../ts/users'
import MetadataStore from '../../../ts/metadata'
import ScriptBuildRunStatus from '../../generic/code/ScriptBuildRunStatus.vue'

@Component({
    components: {
        DashboardAppBar,
        DashboardHomePageNavBar,
        LogViewer,
        ScriptBuildRunStatus,
    }
})
export default class DashboardOrgSingleRunLog extends Vue {
    run : ScriptRun | null = null
    script : ClientScript | null = null
    code : ManagedCode | null = null

    get url() : string {
        return createSingleScriptUrl(PageParamsStore.state.organization!.OktaGroupName, this.script!.Id)
    }

    get isLoading() : boolean {
        return !this.run || !this.script || !this.code
    }

    get runId() : number {
        return parseInt(PageParamsStore.state.resource!.Id, 10)
    }

    get runTimeStr() : string {
        return standardFormatTime(this.run!.StartTime)
    }

    get runByUser() : string {
        return createUserString(MetadataStore.getters.getUser(this.run!.UserId))
    }

    refreshData() {
        getCodeRun({
            orgId: PageParamsStore.state.organization!.Id,
            runId: this.runId
        }).then((resp:  TGetCodeRunOutput) => {
            this.run = resp.data
            getClientScriptCodeFromLink({
                orgId: PageParamsStore.state.organization!.Id,
                linkId: this.run.LinkId,
            }).then((resp : TGetClientScriptCodeFromLinkOutput) => {
                this.script = resp.data.Script
                this.code = resp.data.Code
            }).catch((err: any) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
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
}

</script>

<style scoped>

.item-content-no-flex {
    flex: auto !important;
}

</style>
