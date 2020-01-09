<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-list-item two-line class="pa-0">
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Server: {{ currentServer.Name }}
                        <v-btn icon @click="expandDescription = !expandDescription">
                            <v-icon small v-if="!expandDescription" >mdi-chevron-down</v-icon>
                            <v-icon small v-else>mdi-chevron-up</v-icon>
                        </v-btn>
                    </v-list-item-title>

                    <v-list-item-subtitle :class="expandDescription ? `long-text` : `hide-long-text`">
                        {{ currentServer.Description }}
                    </v-list-item-subtitle>
                </v-list-item-content>

                <v-dialog v-model="showHideDelete"
                          persistent
                          max-width="40%"
                >
                    <template v-slot:activator="{ on }">
                        <v-btn color="error" v-on="on">
                            Delete
                        </v-btn>
                    </template>

                    <generic-delete-confirmation-form
                        item-name="servers"
                        :items-to-delete="[currentServer.Name]"
                        :use-global-deletion="false"
                        @do-cancel="showHideDelete = false"
                        @do-delete="onDelete">
                    </generic-delete-confirmation-form>
                </v-dialog>

            </v-list-item>
            <v-divider></v-divider>

            <v-container fluid>
                <v-row>
                    <v-col cols="8">
                        <create-new-server-form
                            edit-mode
                            :reference-server="currentServer"
                            ref="editForm"
                            @do-save="onEdit"
                        >
                        </create-new-server-form>
                    </v-col>

                    <v-col cols="4">
                    </v-col>
                </v-row>
            </v-container>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import CreateNewServerForm from './CreateNewServerForm.vue'
import { Server } from '../../../ts/infrastructure'
import { getServer, TGetServerOutput } from '../../../ts/api/apiServers'
import { deleteServer } from '../../../ts/api/apiServers'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl, createOrgServersUrl } from '../../../ts/url'

@Component({
    components: {
        GenericDeleteConfirmationForm,
        CreateNewServerForm
    }
})
export default class FullEditServerComponent extends Vue {
    currentServer: Server = {} as Server
    ready : boolean = false
    expandDescription : boolean = false
    showHideDelete: boolean = false

    $refs! : {
        editForm : CreateNewServerForm
    }

    refreshData() {
        let data = window.location.pathname.split('/')
        let resourceId = Number(data[data.length - 1])

        getServer({
            serverId: resourceId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetServerOutput) => {
            this.currentServer = resp.data.Server
            this.ready = true

            Vue.nextTick(() => {
                this.$refs.editForm.clearForm()
            })
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

    onDelete() {
        deleteServer({
            serverId: this.currentServer.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then(() => {
            window.location.replace(createOrgServersUrl(PageParamsStore.state.organization!.OktaGroupName))
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

    onEdit(s : Server) {
        this.currentServer = s
    }
}

</script>

