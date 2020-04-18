<template>
    <v-tooltip right color="rgba(0, 0, 0, 0.9)">
        <template v-slot:activator="{ on }">
            <v-chip small v-on="on">
                {{ shortHash }}
            </v-chip>
        </template>
        {{ hash }}
    </v-tooltip>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'

const Props = Vue.extend({
    props: {
        hash: String
    }
})

@Component
export default class HashRenderer extends Props {
    get shortHash() : string {
        return this.hash.slice(0, 10)
    }
}

</script>
