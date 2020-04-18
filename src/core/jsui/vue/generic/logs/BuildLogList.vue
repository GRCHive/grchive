<template>
    <build-log-table
        :resources="allCode"
    >
    </build-log-table>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import BuildLogTable from '../BuildLogTable.vue'
import { ManagedCode } from '../../../ts/code'
import { TAllCodeOutput, allCode } from '../../../ts/api/apiCode'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'

const Props = Vue.extend({
    props: {
    }
})

@Component({
    components: {
        BuildLogTable,
    }
})
export default class BuildLogList extends Props {
    allCode : ManagedCode[] = []

    refreshData(){ 
        allCode({
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllCodeOutput) => {
            this.allCode = resp.data
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    mounted() {
        this.refreshData()
    }
}

</script>
