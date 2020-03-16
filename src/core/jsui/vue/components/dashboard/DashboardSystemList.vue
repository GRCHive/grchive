<template>
    <div class="ma-4">
        <system-table-with-controls></system-table-with-controls>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import SystemTableWithControls from '../../generic/resources/SystemTableWithControls.vue'

@Component({
    components: {
        SystemTableWithControls
    }
})
export default class DashboardSystemList extends Vue {
}

</script>
