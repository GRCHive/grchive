<template>
    <div class="ma-4">
        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    Systems
                </v-list-item-title>
            </v-list-item-content>

            <v-spacer></v-spacer>

            <v-list-item-action>
                <v-dialog v-model="showHideNew"
                          persistent
                          max-width="40%"
                >
                    <template v-slot:activator="{ on }">
                        <v-btn color="primary" v-on="on">
                            New
                        </v-btn>
                    </template>

                    <create-new-system-form
                        @do-cancel="showHideNew = false"
                        @do-save="saveSystem">
                    </create-new-system-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>
        <v-divider></v-divider>
        <systems-table
            :resources="systems"
        ></systems-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import CreateNewSystemForm from './CreateNewSystemForm.vue'
import { System } from '../../../ts/systems'
import { TAllSystemsOutputs, getAllSystems } from '../../../ts/api/apiSystems'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'
import SystemsTable from '../../generic/SystemsTable.vue'

@Component({
    components: {
        CreateNewSystemForm,
        SystemsTable
    }
})
export default class DashboardSystemList extends Vue {
    showHideNew: boolean = false
    systems : System[] = []

    saveSystem(newSys : System) {
        this.showHideNew = false
        this.systems.push(newSys)
    }

    refreshSystems() {
        getAllSystems({
            orgId: PageParamsStore.state.organization!.Id
        }).then((resp : TAllSystemsOutputs) => {
            this.systems = resp.data
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops. Something went wrong. Try again.",
                false,
                "",
                contactUsUrl,
                true);
        })
    }

    mounted() {
        this.refreshSystems()
    }
}

</script>
