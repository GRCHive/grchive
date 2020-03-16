<template>
    <div class="ma-4">
        <server-table-with-controls>
        </server-table-with-controls>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import ServerTableWithControls from '../../generic/resources/ServerTableWithControls.vue'

@Component({
    components: {
        ServerTableWithControls,
    }
})
export default class DashboardServerList extends Vue {
}

</script>
