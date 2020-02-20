<template>
    <systems-table
        :resources="allSystems"
        :value="value"
        selectable
        @input="onInput"
    >
    </systems-table>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import SystemsTable from './SystemsTable.vue'
import { System } from '../../ts/systems'
import { getAllSystems, TAllSystemsOutputs } from '../../ts/api/apiSystems'
import { PageParamsStore } from '../../ts/pageParams'
import { contactUsUrl } from '../../ts/url'

const Props = Vue.extend({
    props: {
        value: Array
    }
})

@Component({
    components: {
        SystemsTable
    }
})
export default class SystemSearchFormComponent extends Props {
    allSystems: System[] = []

    loadSystems() {
        getAllSystems({
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp: TAllSystemsOutputs) => {
            this.allSystems = resp.data
        }).catch((err: any) => {
            //@ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong, please reload the page and try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    mounted() {
        this.loadSystems()
    }

    onInput(v : System[]) {
        this.$emit('input', v)
    }
}

</script>
