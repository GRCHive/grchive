<template>
    <div>
        <v-row justify="center" class="mb-2">
            <slot name="image"></slot>
        </v-row>

        <v-row justify="center">
            <p class="headline font-weight-bold text-center">{{ headline }}</p>
        </v-row>

        <v-row justify="center">
            <p class="subtitle-1">
                {{ description }}
            </p>
        </v-row>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'

const Props = Vue.extend({
    props: {
        headline: String,
        description: String,
    }
})

@Component
export default class FeatureSummary extends Props {
}

</script>
