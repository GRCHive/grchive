<template>
    <div class="ma-4">
        <db-table-with-controls></db-table-with-controls>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import DbTableWithControls from '../../generic/resources/DbTableWithControls.vue'

@Component({
    components: {
        DbTableWithControls
    }
})
export default class DashboardDatabaseList extends Vue {
}

</script>
