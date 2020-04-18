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
                                {{ !!script ? "Script" : "Data" }}:
                            </span>
                            {{ name }}
                        </a>
                    </v-list-item-title>

                    <v-list-item-subtitle>
                        {{ code.GitHash }}
                    </v-list-item-subtitle>
                </v-list-item-content>

                <v-list-item-content>
                    <v-list-item-title>
                        <span class="font-weight-bold">
                            Modified By: 
                        </span>
                        {{  modByUser }}
                    </v-list-item-title>

                    <v-list-item-subtitle>
                        <span class="font-weight-bold">
                            Modified At: 
                        </span>
                        {{ modTimeStr }}
                    </v-list-item-subtitle>
                </v-list-item-content>

                <v-list-item-content>
                    <span class="font-weight-bold overline item-content-no-flex">
                        Build Status: 
                    </span>
                    <code-build-status
                        class="item-content-no-flex"
                        :commit="commit"
                        show-time-stamp
                    >
                    </code-build-status>
                </v-list-item-content>
            </v-list-item>

            <log-viewer :commit="commit"></log-viewer>
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
import { ClientData } from '../../../ts/clientData'
import { ManagedCode } from '../../../ts/code'
import {
    getCodeLink, TGetCodeLinkOutput,
    getCode, TGetCodeOutput
} from '../../../ts/api/apiCode'
import {
    contactUsUrl,
    createSingleScriptUrl,
    createSingleClientDataUrl,
} from '../../../ts/url'
import { standardFormatTime } from '../../../ts/time'
import { createUserString } from '../../../ts/users'
import MetadataStore from '../../../ts/metadata'
import CodeBuildStatus from '../../generic/code/CodeBuildStatus.vue'

@Component({
    components: {
        DashboardAppBar,
        DashboardHomePageNavBar,
        LogViewer,
        CodeBuildStatus
    }
})
export default class DashboardOrgSingleBuildLog extends Vue {
    code : ManagedCode | null = null
    script : ClientScript | null = null
    data : ClientData | null = null

    get name() : string {
        return !!this.script ? this.script.Name : this.data!.Name
    }

    get url() : string {
        return !!this.script ?
            createSingleScriptUrl(PageParamsStore.state.organization!.OktaGroupName, this.script!.Id) :
            createSingleClientDataUrl(PageParamsStore.state.organization!.OktaGroupName, this.data!.Id)
    }

    get modByUser() : string {
        return createUserString(MetadataStore.getters.getUser(this.code!.UserId))
    }

    get modTimeStr() : string {
        return standardFormatTime(this.code!.ActionTime)
    }

    get isLoading() : boolean {
        return !this.code || (!this.script && !this.data)
    }

    get commit() : string {
        return PageParamsStore.state.resource!.Id
    }

    refreshData() {
        getCode({
            orgId: PageParamsStore.state.organization!.Id,
            codeCommit: this.commit,
        }).then((resp : TGetCodeOutput) => {
            this.code = resp.data.Full
            getCodeLink({
                orgId: PageParamsStore.state.organization!.Id,
                codeId: this.code.Id,
            }).then((resp : TGetCodeLinkOutput) => {
                this.script = resp.data.Script
                this.data = resp.data.Data
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

<style scoped>

.item-content-no-flex {
    flex: auto !important;
}

</style>
