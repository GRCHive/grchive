<template>
    <div class="max-height">
        <dashboard-app-bar ref="dashboardAppBar">
        </dashboard-app-bar>

        <dashboard-home-page-nav-bar></dashboard-home-page-nav-bar>

        <v-content>
            <v-overlay :value="loading">
                <v-progress-circular indeterminate size="64"></v-progress-circular>
            </v-overlay>

            <div v-if="!loading">
                <v-list-item two-line>
                    <v-list-item-content>
                        <v-list-item-title class="title">
                            {{shellTypeStr}} Script: 
                            <a :href="shellScriptUrl">{{ script.Name }} v{{ versionNum }}</a>
                        </v-list-item-title>

                        <v-list-item-subtitle class="subtitle-1">
                            <span class="font-weight-bold">Run By:</span> {{ runUser }}
                        </v-list-item-subtitle>
                    </v-list-item-content>

                    <v-list-item-content>
                        <div>
                            <div>
                                <span>
                                    <span class="font-weight-bold">Start Time:</span>
                                    {{ startTimeStr }}
                                </span>
                            </div>

                            <div>
                                <span>
                                    <span class="font-weight-bold">End Time:</span>
                                    {{ endTimeStr }}
                                </span>
                            </div>
                        </div>
                    </v-list-item-content>

                    <v-list-item-content>
                        <span class="font-weight-bold" style="flex: 0 1;">Status: </span>
                        <div style="flex: 1 1;">
                            <shell-run-status
                                :server-runs="serverRuns"
                            >
                            </shell-run-status>
                        </div>
                    </v-list-item-content>
                </v-list-item>

                <v-divider></v-divider>

                <v-tabs>
                    <v-tab>Logs</v-tab>
                    <v-tab-item>
                        <v-tabs vertical>
                            <template v-for="(item, idx) in serverRuns">
                                <v-tab :key="`tab-${idx}`">
                                    <span v-if="hasServerInfo(item.ServerId)">{{ singleServerInfo(item.ServerId).Name }}</span>
                                    <span v-else>Loading...</span>
                                </v-tab>

                                <v-tab-item :key="`tab-item-${idx}`">
                                    <generic-log-viewer
                                        full-height
                                        :raw-log="item.EncryptedLog"
                                    >
                                    </generic-log-viewer>
                                </v-tab-item>
                            </template>
                        </v-tabs>
                    </v-tab-item>

                    <v-tab>Script</v-tab>
                    <v-tab-item>
                        <generic-code-editor
                            v-if="!!scriptText"
                            v-model="scriptText"
                            :lang="shellCMLanguage"
                            readonly
                            full-height
                        >
                        </generic-code-editor>
                    </v-tab-item>
                </v-tabs>

            </div>
        </v-content>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import DashboardAppBar from '../../components/dashboard/DashboardAppBar.vue'
import DashboardHomePageNavBar from '../../components/dashboard/DashboardHomePageNavBar.vue'
import ShellRunStatus from '../../generic/shell/ShellRunStatus.vue'
import MetadataStore from '../../../ts/metadata'
import GenericCodeEditor from '../../generic/code/GenericCodeEditor.vue'
import GenericLogViewer from '../../generic/logs/GenericLogViewer.vue'
import { createUserString } from '../../../ts/users'
import {
    ShellScriptRun,
    ShellScript,
    ShellScriptVersion,
    ShellScriptRunPerServer,
    ShellTypes,
    ShellTypeToCodeMirror,
} from '../../../ts/shell'
import {
    getShellRunInformation, TGetShellRunOutput
} from '../../../ts/api/apiShellRun'
import {
    getServer, TGetServerOutput
} from '../../../ts/api/apiServers'
import {
    getShellScriptVersion, TGetShellScriptVersionOutput,
} from '../../../ts/api/apiShell'
import {
    contactUsUrl,
    createSingleShellUrl,
} from '../../../ts/url'
import {
    standardFormatTime
} from '../../../ts/time'
import {
    Server
} from '../../../ts/infrastructure'
import { PageParamsStore } from '../../../ts/pageParams'

@Component({
    components : {
        DashboardAppBar,
        DashboardHomePageNavBar,
        ShellRunStatus,
        GenericCodeEditor,
        GenericLogViewer,
    },
})
export default class DashboardOrgSingleShellRun extends Vue {
    run : ShellScriptRun | null = null
    script : ShellScript | null = null
    version: ShellScriptVersion | null = null
    versionNum: number = -1
    serverRuns : ShellScriptRunPerServer[] | null = null
    serverInfo: Record<number, Server> = Object()
    scriptText : string | null = null

    get hasServerInfo() : (idx : number) => boolean {
        return (idx : number) => {
            return idx in this.serverInfo
        }
    }

    get singleServerInfo() : (idx : number) => Server {
        return (idx : number) => {
            return this.serverInfo[idx]
        }
    }

    get shellTypeStr() : string {
        return ShellTypes[this.script!.TypeId]
    }

    get shellCMLanguage() : string {
        return ShellTypeToCodeMirror.get(<ShellTypes>this.script!.TypeId)!
    }

    get shellScriptUrl() : string {
        return createSingleShellUrl(PageParamsStore.state.organization!.OktaGroupName, this.script!.Id) + `?version=${this.versionNum}`
    }

    get loading() : boolean {
        return (!this.run || !this.script || !this.version || !this.serverRuns)
    }

    get runUser() : string {
        return createUserString(MetadataStore.getters.getUser(this.run!.RunUserId))
    }

    get startTimeStr() : string {
        if (!this.run!.RunTime) {
            return "Pending..."
        }
        return standardFormatTime(this.run!.RunTime!)
    }

    get endTimeStr() : string {
        if (!this.run!.EndTime) {
            return "Pending..."
        }
        return standardFormatTime(this.run!.EndTime!)
    }

    loadServerInformation() {
        if (!this.serverRuns) {
            return
        }

        // Load information one by one in a chain to prevent a mass
        // of HTTP requests going out at the same time in the case where
        // a lot of servers are being requested.
        let idx = 0 
        let loadServerInformationHelper = () => {
            if (idx >= this.serverRuns!.length) {
                return
            }

            getServer({
                orgId: PageParamsStore.state.organization!.Id,
                serverId: this.serverRuns![idx].ServerId,
            }).then((resp : TGetServerOutput) => {
                Vue.set(this.serverInfo, resp.data.Server.Id, resp.data.Server)

                idx += 1
                loadServerInformationHelper()
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
        loadServerInformationHelper()
    }

    refreshData() {
        getShellRunInformation({
            orgId: PageParamsStore.state.organization!.Id,
            runId: parseInt(PageParamsStore.state.resource!.Id, 10),
            includeLogs: true,
        }).then((resp : TGetShellRunOutput) => {
            this.run = resp.data.Run
            this.script = resp.data.Script
            this.version = resp.data.Version
            this.versionNum = resp.data.VersionNum
            this.serverRuns = resp.data.ServerRuns

            this.loadServerInformation()

            getShellScriptVersion({
                orgId: PageParamsStore.state.organization!.Id,
                shellId: this.script!.Id,
                version: this.version!.Id,
            }).then((resp : TGetShellScriptVersionOutput) => {
                this.scriptText = resp.data
            }).catch((err : any) => {
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
