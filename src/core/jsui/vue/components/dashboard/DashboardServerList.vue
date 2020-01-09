<template>
    <div class="ma-4">
        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    Servers
                </v-list-item-title>
            </v-list-item-content>

            <v-list-item-action>
                <v-text-field outlined
                              v-model="filterText"
                              prepend-inner-icon="mdi-magnify"
                              hide-details
                ></v-text-field>
            </v-list-item-action>

            <v-spacer></v-spacer>

            <v-list-item-action>
                <v-dialog v-model="showHideNew"
                          persistent
                          max-width="40%"
                          ref="form"
                >
                    <template v-slot:activator="{ on }">
                        <v-btn color="primary" v-on="on">
                            New
                        </v-btn>
                    </template>

                    <create-new-server-form
                        @do-save="onSave"
                        @do-cancel="showHideNew = false">
                    </create-new-server-form>
                </v-dialog>
            </v-list-item-action>

        </v-list-item>
        <v-divider></v-divider>

        <server-table
            :resources="servers"
            :search="filterText"
        >
        </server-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import ServerTable from '../../generic/ServerTable.vue'
import { Server } from '../../../ts/infrastructure'
import { allServers, TAllServerOutput } from '../../../ts/api/apiServers'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'
import CreateNewServerForm from './CreateNewServerForm.vue'

@Component({
    components: {
        ServerTable,
        CreateNewServerForm
    }
})
export default class DashboardServerList extends Vue {
    showHideNew : boolean = false
    filterText: string = ""
    servers : Server[] = []

    $refs!: {
        form: CreateNewServerForm
    }

    onSave(s : Server) {
        this.showHideNew = false
        this.servers.push(s)
        this.$refs.form.clearForm()
    }

    mounted() {
        allServers({
            orgId: PageParamsStore.state.organization!.Id
        }).then((resp : TAllServerOutput) => {
            this.servers = resp.data
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
}

</script>
