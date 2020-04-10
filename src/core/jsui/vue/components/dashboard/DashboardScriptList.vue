<template>
    <div class="ma-4">
        <client-script-table-with-controls></client-script-table-with-controls>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import ClientScriptTableWithControls from '../../generic/resources/ClientScriptTableWithControls.vue'

@Component({
    components: {
        ClientScriptTableWithControls
    }
})
export default class DashboardScriptList extends Vue {
}
</script>

