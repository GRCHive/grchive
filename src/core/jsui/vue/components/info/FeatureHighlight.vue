<template>
    <v-row justify="center" class="py-12">
        <v-col align-self="center" cols="4" :offset="textOffset" :order="textOrder">
            <p class="display-1 font-weight-bold">{{ headline }}</p>
            <p class="subtitle-1">
                {{ description }}
            </p>
            <v-btn
                class="subtitle-1 pa-0"
                :href="ctaUrl"
                text
                color="accent"
                v-if="!!ctaText && !!ctaUrl"
            >
                {{ ctaText }}
                <v-icon>
                    mdi-chevron-right
                </v-icon>
            </v-btn>
        </v-col>

        <v-col :cols="largeImage ? 6 : 4" :offset="imageOffset" :order="imageOrder" align-self="center">
            <v-row justify="center">
                <slot name="image"></slot>
            </v-row>
        </v-col>
    </v-row>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'

const Props = Vue.extend({
    props: {
        headline: String,
        description: String,
        ctaUrl: String,
        ctaText: String,
        flip: {
            type: Boolean,
            default: false,
        },
        largeImage: {
            type: Boolean,
            default: false,
        }
    }
})

@Component
export default class FeatureHighlight extends Props {
    get imageOrder() : string {
        return this.flip ? "first" : "last"
    }

    get imageOffset() : number {
        return this.flip ? 0 : 1
    }

    get textOrder() : string {
        return this.flip ? "last" : "first"
    }

    get textOffset() : number {
        return this.flip ? 1 : 0
    }
}

</script>
