<template>
    <div class="ma-4">
        <v-row>
            <v-col cols="6" v-if="previewReady">
            </v-col>

            <v-col cols="6" v-if="!previewReady" align="center" justify="center">
                <v-progress-circular
                    indeterminate
                    size="64"
                ></v-progress-circular>
            </v-col>

            <v-col cols="6" v-if="metadataReady">
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
            <v-divider></v-divider>

            </v-col>

            <v-col cols="6" v-if="!metadataReady" align="center" justify="center">
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
import { contactUsUrl, createSingleDocCatUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import { getSingleControlDocument, TGetSingleControlDocumentOutput } from '../../../ts/api/apiControlDocumentation'
import { ControlDocumentationFile, ControlDocumentationCategory } from '../../../ts/controls'

@Component
export default class FullEditDocumentationComponent extends Vue {
    ready: boolean = false
    parentCat : ControlDocumentationCategory | null = null
    metadata : ControlDocumentationFile | null = null

    get previewReady() : boolean {
        return false
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

    refreshData() {
        let data = window.location.pathname.split('/')
        let resourceId = Number(data[data.length - 1])

        getSingleControlDocument({
            fileId: resourceId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetSingleControlDocumentOutput) => {
            this.parentCat = resp.data.Category
            this.metadata = resp.data.File
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
