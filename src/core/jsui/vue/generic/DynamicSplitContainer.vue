<template>
    <div class="dynParent">
        <div class="dynChildA" :style="colAStyle">
            <slot
                name="first-col"
            >
            </slot>
        </div>

        <template v-if="enableColB">
            <div class="dynDivider"
                 @mousedown="onClickDivider"
            ></div>

            <div class="dynChildB" :style="colBStyle">
                <slot name="second-col"></slot>
            </div>
        </template>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'

const Props = Vue.extend({
    props: {
        enableColB: {
            type: Boolean,
            default: false,
        }
    }
})

@Component
export default class DynamicSplitContainer extends Props {
    dividerActive: boolean = false

    colAOffset: number = 0

    get colAStyle() : any {
        if (this.enableColB) {
            return {
                'max-width': `calc(50% + ${this.colAOffset}px - 4px)`,
                'min-width': '20px',
            }
        } else {
            return {}
        }
    }


    get colBStyle() : any {
        if (this.enableColB) {
            return {
                'max-width': `calc(50% - ${this.colAOffset}px - 4px)`,
                'min-width': '20px',
            }
        } else {
            return {}
        }
    }


    onClickDivider() {
        this.dividerActive = true
    }

    onReleaseDivider() {
        this.dividerActive = false
    }

    onMoveDivider(e : MouseEvent) {
        if (!this.dividerActive) {
            return
        }

        this.colAOffset += e.movementX
    }

    mounted() {
        document.addEventListener('mousemove', this.onMoveDivider)
        document.addEventListener('mouseup', this.onReleaseDivider)
    }
}

</script>

<style scoped>

.dynParent {
    display: flex;
}

.dynChildA {
    box-sizing: border-box;
    flex: 5 1 auto;
}

.dynChildB {
    box-sizing: border-box;
    flex: 1 1 auto;
}

.dynDivider {
    width: 8px;
    cursor: ew-resize;
    flex: 0 0 auto;
    background-color: gray;
}

</style>
