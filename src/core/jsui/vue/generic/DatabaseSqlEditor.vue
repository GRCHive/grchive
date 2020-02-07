<template>
    <v-container fluid>
        <v-row v-if="isLoading" align="center" justify="center">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-row>

        <div v-else>
            <v-tabs vertical>
                <v-tab>Schemas</v-tab>
                <v-tab-item>
                </v-tab-item>

                <v-tab>Queries</v-tab>
                <v-tab-item>
                </v-tab-item>
            </v-tabs>

            <div id="queryOutput">
            </div>
        </div>
    </v-container>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { PageParamsStore } from '../../ts/pageParams'

const Props = Vue.extend({
    props: {
        dbId: Number
    }
})

@Component
export default class DatabaseSqlEditor extends Props {
    isLoading: boolean = true

    refreshData() {
    }

    mounted() {
    }
}

</script>
