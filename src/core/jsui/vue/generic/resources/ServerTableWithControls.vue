<template>
    <div>
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

            <v-list-item-action v-if="!disableNew">
                <v-dialog v-model="showHideNew"
                          persistent
                          max-width="40%"
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
            :resources="filteredServers"
            :search="filterText"
            :use-crud-delete="!disableDelete"
            :confirm-delete="!disableDelete"
            @delete="deleteServer"
            @input="modifySelected"
            :selectable="enableSelect"
            :multi="enableSelect"
        >
        </server-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import ServerTable from '../ServerTable.vue'
import { Server } from '../../../ts/infrastructure'
import { allServers, TAllServerOutput } from '../../../ts/api/apiServers'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'
import CreateNewServerForm from '../../components/dashboard/CreateNewServerForm.vue'
import { deleteServer } from '../../../ts/api/apiServers'

const Props = Vue.extend({
    props: {
        value: {
            type: Array,
            default: () => [],
        },
        exclude: {
            type: Array,
            default: () => [],
        },
        disableNew: {
            type: Boolean,
            default: false,
        },
        disableDelete: {
            type: Boolean,
            default: false,
        },
        enableSelect: {
            type: Boolean,
            default: false,
        }
    }
})

@Component({
    components: {
        ServerTable,
        CreateNewServerForm
    }
})
export default class ServerTableWithControls extends Props {
    showHideNew : boolean = false
    filterText: string = ""
    servers : Server[] = []

    get excludeSet() : Set<number> {
        return new Set<number>(this.exclude.map((ele : any) => ele.Id))
    }

    get filteredServers() : Server[] {
        return this.servers.filter((ele : Server) => !this.excludeSet.has(ele.Id))
    }

    onSave(s : Server) {
        this.showHideNew = false
        this.servers.unshift(s)
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

    modifySelected(vals : Server[]) {
        this.$emit('input', vals)
    }

    deleteServer(server : Server) {
        deleteServer({
            serverId: server.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then(() => {
            this.servers.splice(
                this.servers.findIndex((ele : Server) =>
                    ele.Id == server.Id),
                1)
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
}


</script>
