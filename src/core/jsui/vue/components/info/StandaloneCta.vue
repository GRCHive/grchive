<template>
    <v-row class="max-height" align="center" justify="center" no-gutters>
        <v-col cols="8">
            <v-row justify="center" class="mb-8">
                <p class="display-1 font-weight-bold white--text">
                    Ready to transform your audit process? Request a free demo today.
                </p>
            </v-row>

            <v-row justify="center">
                <v-btn 
                    color="success"
                    :href="getStartedUrl"
                    x-large
                >
                    Request Demo
                </v-btn>
            </v-row>
        </v-col>
    </v-row>
</template>


<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { getStartedUrl } from '../../../ts/url'

@Component
export default class StandaloneCta extends Vue {
    getStartedUrl: string = getStartedUrl
}

</script>
