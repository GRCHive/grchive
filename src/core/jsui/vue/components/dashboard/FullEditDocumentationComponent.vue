<template>
    <div class="max-height">
        <v-row class="max-height">
            <v-col cols="6" v-if="previewReady" class="max-height py-0">
                <pdf-js-viewer :pdf="previewBase64" ref="pdfViewer" :max-viewer-height="viewerMaxHeight">
                </pdf-js-viewer>
            </v-col>

            <v-col 
                cols="6"
                v-else
                align="center"
                justify="center"
            >
                <v-progress-circular
                    indeterminate
                    size="64"
                    v-if="hasPreview"
                ></v-progress-circular>

                <span v-else>No Preview Available</span>
            </v-col>

            <v-col cols="6" v-if="metadataReady">
                <div class="mt-4">
                    <v-list-item two-line class="pa-0">
                        <v-list-item-content>
                            <v-list-item-title class="title">
                                File: {{ metadata.AltName }}
                            </v-list-item-title>

                            <v-list-item-subtitle>
                                Parent Category: <a :href="parentCatUrl">{{ parentCat.Name }}</a>
                            </v-list-item-subtitle>
                        </v-list-item-content>
                    </v-list-item>
                </div>
            <v-divider></v-divider>

            </v-col>

            <v-col cols="6" v-else align="center" justify="center">
                <v-progress-circular
                    indeterminate
                    size="64"
                ></v-progress-circular>
            </v-col>
        </v-row>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { contactUsUrl, createSingleDocCatUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import { 
    getSingleControlDocument,
    TGetSingleControlDocumentOutput,
    downloadSingleControlDocument,
    TDownloadSingleControlDocumentOutput
} from '../../../ts/api/apiControlDocumentation'
import { ControlDocumentationFile, ControlDocumentationCategory } from '../../../ts/controls'
import PdfJsViewer from '../../generic/pdf/PdfJsViewer.vue'

@Component({
    components: {
        PdfJsViewer
    }
})
export default class FullEditDocumentationComponent extends Vue {
    ready: boolean = false
    parentCat : ControlDocumentationCategory | null = null
    metadata : ControlDocumentationFile | null = null
    viewerMaxHeight: number = 100

    previewMetadata : ControlDocumentationFile | null = null
    previewData : Blob | null = null
    previewBase64 : string | null = null

    $refs!: {
        pdfViewer: PdfJsViewer
    }

    get hasPreview() : boolean {
        return !!this.previewMetadata
    }

    get previewReady() : boolean {
        return !!this.previewBase64
    }

    @Watch('previewData')
    encodePreview() {
        if (!this.previewData) {
            this.previewBase64 = null
            return
        }
        
        let reader = new FileReader()
        reader.onload = (e) => {
            if (!e.target) {
                return
            }

            let blobStr : string = <string>e.target.result
            this.previewBase64 = blobStr.replace(/^data:.*\/.*;base64,/g, '')
        }
        reader.readAsDataURL(this.previewData!)
    }

    get metadataReady() : boolean {
        return !!this.metadata
    }

    get parentCatUrl() : string {
        if (!this.parentCat) {
            return "#"
        }

        return createSingleDocCatUrl(PageParamsStore.state.organization!.OktaGroupName, this.parentCat!.Id)
    }

    reloadPreview() {
        if (!this.previewMetadata) {
            return
        }

        downloadSingleControlDocument({
            fileId: this.previewMetadata!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            catId: this.previewMetadata!.CategoryId,
        }).then((resp : TDownloadSingleControlDocumentOutput) => {
            this.previewData = resp.data
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

    refreshData() {
        let data = window.location.pathname.split('/')
        let resourceId = Number(data[data.length - 1])

        getSingleControlDocument({
            fileId: resourceId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetSingleControlDocumentOutput) => {
            this.parentCat = resp.data.Category
            this.previewMetadata = resp.data.PreviewFile
            this.metadata = resp.data.File
            this.reloadPreview()
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
        window.addEventListener("resize", this.updateViewerRect)
    }

    @Watch('previewReady')
    updateViewerRect() {
        Vue.nextTick(() => {
            if (!this.previewReady || !this.$refs.pdfViewer) {
                return
            }

            let viewerRect = this.$refs.pdfViewer.$el.getBoundingClientRect()
            let windowHeight = window.innerHeight
            this.viewerMaxHeight = windowHeight - viewerRect.y
        })
    }
}

</script>
