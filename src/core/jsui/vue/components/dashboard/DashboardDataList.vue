<template>
    <div class="ma-4">
        <client-data-table-with-controls></client-data-table-with-controls>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import ClientDataTableWithControls from '../../generic/resources/ClientDataTableWithControls.vue'

@Component({
    components: {
        ClientDataTableWithControls
    }
})
export default class DashboardDataList extends Vue {
}
</script>
