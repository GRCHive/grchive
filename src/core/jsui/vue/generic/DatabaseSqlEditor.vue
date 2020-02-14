<template>
    <v-container fluid>
        <v-tabs vertical>
            <v-tab>Schemas</v-tab>
            <v-tab-item>
                <database-refresh-manager
                    :db-id="dbId"
                >
                </database-refresh-manager>
            </v-tab-item>

            <v-tab>Queries</v-tab>
            <v-tab-item>
                <database-query-manager
                    :db-id="dbId"
                >
                </database-query-manager>
            </v-tab-item>
        </v-tabs>
    </v-container>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import DatabaseRefreshManager from './DatabaseRefreshManager.vue'
import DatabaseQueryManager from './DatabaseQueryManager.vue'

const Props = Vue.extend({
    props: {
        dbId: Number,
    }
})

@Component({
    components: {
        DatabaseRefreshManager,
        DatabaseQueryManager,
    }
})
export default class DatabaseSqlEditor extends Props {

}

</script>
