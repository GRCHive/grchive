<template>
    <div :style="parentDivStyle">
        <div v-if="visible">
            <canvas
                :width="viewport.width"
                :height="viewport.height"
                ref="canvas"
            >
            </canvas>
            <div
                class="textLayer"
                :style="textLayerStyle"
                ref="textLayer"
            >
            </div>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { PDFPageProxy, PageViewport } from 'pdfjs-dist/build/pdf'
import { TextLayerBuilder } from 'pdfjs-dist/web/pdf_viewer'

const Props = Vue.extend({
    props: {
        page: PDFPageProxy,
        viewport: PageViewport,
        visible: Boolean
    }
})

@Component
export default class PdfPageRenderer extends Props {
    $refs! : {
        canvas : HTMLCanvasElement
        textLayer: HTMLDivElement
    }

    @Watch("visible")
    @Watch("viewport")
    @Watch("page")
    checkRender() {
        if (!this.visible) {
            return
        }

        Vue.nextTick(() => {
            let context = this.$refs.canvas.getContext("2d");
            // Render PDF page into canvas context
            let renderContext = {
                canvasContext: context,
                viewport: this.viewport,
            };
            this.page.render(renderContext).promise.then(() => {
                this.page.getTextContent().then((textContent) => {
                    let textBuilder = new TextLayerBuilder({
                        textLayerDiv: this.$refs.textLayer, 
                        pageIndex: this.page.pageIndex,
                        viewport: this.viewport,
                    })

                    textBuilder.setTextContent(textContent)
                    textBuilder.render()
                })
            })
        })
    }
    
    get parentDivStyle() : any {
        return {
            "background-color": "white",
            "width": `${this.viewport.width}px`,
            "height": `${this.viewport.height}px`,
            "position": "relative",
        }
    }

    get textLayerStyle() : any {
        return {
            "width": `${this.viewport.width}px`,
            "height": `${this.viewport.height}px`,
        }
    }

    mounted() {
        this.checkRender()
    }
}

</script>

<style scoped>

.textLayer {
    position: absolute;
    left: 0;
    top: 0;
    right: 0;
    bottom: 0;
    overflow: hidden;
    opacity: 0.2;
    line-height: 1.0;
}

.textLayer ::selection {
    background: rgb(0, 0, 255);
}

>>>.textLayer > span {
    color: transparent;
    position: absolute;
    white-space: pre;
    cursor: text;
    transform-origin: 0% 0%;
}

</style>
