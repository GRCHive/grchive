<template>
    <div id="viewerContainer" :style="viewerContainerStyle">
        <canvas id="viewer">
        </canvas>
    </div>
</template>

<script lang="ts">

import Vue, {VNode} from 'vue'
import Component from 'vue-class-component'
import { contactUsUrl } from '../../../ts/url'

import pdfjsLib from 'pdfjs-dist/webpack'

const Props = Vue.extend({
    props: {
        // Encoded base64
        pdf: String,
        maxViewerHeight: Number
    },
})

@Component
export default class PdfJsViewer extends Props {

    get viewerContainerStyle() : any {
        return {
            "height": `${this.maxViewerHeight}px`
        }
    }

    get binaryData() : string {
        return atob(this.pdf)
    }

    refreshData() {
        let task = pdfjsLib.getDocument({data : this.binaryData})
        task.promise.then((pdf : any) => {
            pdf.getPage(1).then((page : any) => {
                var scale = 1.5;
                var viewport = page.getViewport({scale: scale});

                // Prepare canvas using PDF page dimensions
                var canvas = <HTMLCanvasElement>document.getElementById("viewer");
                var context = canvas.getContext("2d");
                canvas.height = viewport.height;
                canvas.width = viewport.width;

                // Render PDF page into canvas context
                var renderContext = {
                    canvasContext: context,
                    viewport: viewport
                };
                var renderTask = page.render(renderContext);
                renderTask.promise.then(function () {
                  console.log('Page rendered');
                });
            }, (err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
            })
        })
    }

    mounted() {
        this.refreshData()
    }
}

</script>

<style scoped>

#viewerContainer {
    overflow: scroll;
}

</style>
