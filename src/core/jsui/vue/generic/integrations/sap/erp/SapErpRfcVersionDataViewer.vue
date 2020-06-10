<template>
    <div>
        <div v-if="!!fullVersionData">
            <div v-if="isFinishedRunning">
                <div v-if="fullVersionData.Success">
                    <generic-json-viewer
                        :data="fullVersionData.Data"
                    >
                    </generic-json-viewer>
                </div>

                <v-row v-else align="center" justify="center">
                    <div class="max-width">
                        <p class="subtitle-1 text-center">Oops! We were unable to run this SAP ERP RFC.</p>
                        <v-divider class="mb-3"></v-divider>
                        <pre>{{ fullVersionData.Logs }}</pre>
                    </div>
                </v-row>
            </div>

            <div v-else>
                <v-row justify="center">
                    <span class="title">In progress...Check back soon!</span>
                </v-row>
                <v-row justify="center">
                    <v-progress-linear indeterminate></v-progress-linear>
                </v-row>
            </div>
        </div>

        <v-row justify="center" v-else>
            <v-progress-circular size="64" indeterminate></v-progress-circular>
        </v-row>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Prop, Watch } from 'vue-property-decorator'
import { SapErpRfcVersion } from '../../../../../ts/integrations/sap'
import {
    getSapErpRfcVersion, TGetSapErpRfcVersionOutput,
} from '../../../../../ts/api/integrations/apiSapErp'
import { PageParamsStore } from '../../../../../ts/pageParams'
import { contactUsUrl } from '../../../../../ts/url'
import GenericJsonViewer from '../../../json/GenericJsonViewer.vue'

const refreshIntervalSeconds : number = 5

@Component({
    components: {
        GenericJsonViewer,
    }
})
export default class SapErpRfcVersionDataViewer extends Vue {
    @Prop({required: true})
    readonly version! : SapErpRfcVersion | null

    @Prop({required: true})
    readonly integrationId! : number

    fullVersionData : SapErpRfcVersion | null = null

    get isFinishedRunning() : boolean {
        if (!this.fullVersionData) {
            return false
        }

        return !!this.fullVersionData.FinishedTime
    }

    @Watch('version')
    refreshVersion() {
        this.fullVersionData = null
        if (!this.version) {
            return
        }

        getSapErpRfcVersion({
            orgId : PageParamsStore.state.organization!.Id,
            integrationId: this.integrationId,
            rfcId: this.version!.RfcId,
            versionId: this.version!.Id,
        }).then((resp : TGetSapErpRfcVersionOutput) => {
            this.fullVersionData = resp.data
            if (!this.isFinishedRunning) {
                setTimeout(this.refreshVersion, refreshIntervalSeconds*1000)
            }
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true)
        })
    }

    mounted() {
        this.refreshVersion()
    }
}

</script>
