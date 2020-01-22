<template>
    <div id="viewerContainer" :style="viewerContainerStyle" ref="container" @scroll="scrollContainer">
        <div :style="canvasContainerStyle" v-if="readyToRender">
            <pdf-page-renderer
                v-for="(item, i) in allPages"
                :key="i"
                :page="allPages[i]"
                :viewport="pdfViewports[i]"
                :visible="pageVisibility[i]"
                :style="rendererStyle(i)"
            >
            </pdf-page-renderer>
        </div>

        <v-progress-circular
            indeterminate
            size="64"
            v-else
        ></v-progress-circular>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import PdfPageRenderer from './PdfPageRenderer.vue'
import { contactUsUrl } from '../../../ts/url'

import pdfjsLib from 'pdfjs-dist/webpack'
import { PDFDocumentProxy, PDFPageProxy, PageViewport } from 'pdfjs-dist/build/pdf'
import { TRange } from '../../../ts/range'

const Props = Vue.extend({
    props: {
        // Encoded base64
        pdf: String,
        maxViewerHeight: Number
    },
})

const PageSpacerPx = 8

@Component({
    components: {
        PdfPageRenderer,
    }
})
export default class PdfJsViewer extends Props {
    fullPdf : PDFDocumentProxy | null = null
    allPages : Array<PDFPageProxy> = []
    scale : number = 1.0
    currentScroll : number = 0.0

    $refs! : {
        container : HTMLDivElement
    }

    get readyToRender() : boolean {
        if (!this.fullPdf) {
            return false
        }

        for (let i = 0; i < this.pdfViewports.length; ++i) {
            if (typeof this.pdfViewports[i] === "undefined") {
                return false
            }
        }
        return true
    }

    get viewerContainerStyle() : any {
        return {
            "background-color": "#404040",
            "height": `${this.maxViewerHeight}px`,
        }
    }

    get canvasContainerStyle() : any {
        return {
            "width": `${this.totalWidth}px`,
            "height": `${this.totalHeight}px`,
            "margin-left": "auto",
            "margin-right": "auto",
        }
    }

    rendererStyle(i : number) : any {
        return {
            "margin-top": `${i == 0 ? PageSpacerPx : 0}px`,
            "margin-bottom": `${PageSpacerPx}px`,
        }
    }

    get pdfViewports() : Array<PageViewport> {
        return this.allPages.map((ele : PDFPageProxy) => ele.getViewport({scale : this.scale }))
    }

    get binaryData() : string {
        return atob(this.pdf)
    }

    get numPages() : number {
        if (!this.fullPdf) {
            return 0
        }

        return this.fullPdf.numPages
    }

    get totalWidth() : number {
        return Math.max(...this.pdfViewports.map((ele : PageViewport) => ele.width))
    }

    get totalHeight() : number {
        return this.pdfViewports.map((ele : PageViewport) => ele.height).reduce((a,c) => a + c, (this.numPages + 1) * PageSpacerPx)
    }

    get pageVisibility() : Array<boolean> {
        let vis = new Array<boolean>(this.numPages)
        let viewerRange = new TRange<number>(this.currentScroll, this.currentScroll + this.maxViewerHeight)
        let currentPageY = PageSpacerPx

        for (let i = 0; i < this.numPages; ++i) {
            let pageRange = new TRange<number>(currentPageY, currentPageY + this.pdfViewports[i].height) 
            vis[i] = pageRange.intersects(viewerRange)
            currentPageY = pageRange.max + PageSpacerPx
        }

        return vis
    }

    scrollContainer() {
        this.currentScroll = this.$refs.container.scrollTop
    }

    refreshData() {
        let task = pdfjsLib.getDocument({data : this.binaryData})
        task.promise.then((pdf : PDFDocumentProxy) => {
            this.fullPdf = pdf

            this.allPages = new Array<PDFPageProxy>(this.numPages)
            for (let i = 1; i <= this.numPages; ++i) {
                pdf.getPage(i).then((pg : PDFPageProxy) => {
                    Vue.set(this.allPages, i - 1, pg)
                })
            }
        })
    }

    mounted() {
        this.refreshData()
    }
}

</script>

<style scoped>

#viewerContainer {
    overflow: auto;
}

</style>
