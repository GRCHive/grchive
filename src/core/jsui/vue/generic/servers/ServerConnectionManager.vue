<template>
    <div>
        <v-dialog
            persistent
            max-width="40%"
            v-model="showHideNewSSHPassword"
        >
            <create-server-ssh-password-connection
                :server-id="serverId"
                @do-save="onNewSshPasswordConnection"
                @do-cancel="showHideNewSSHPassword = false"
            >
            </create-server-ssh-password-connection>
        </v-dialog>

        <v-dialog
            persistent
            max-width="40%"
            v-model="showHideNewSSHKey"
        >
        </v-dialog>

        <v-card>
            <v-card-title>
                Connection

                <v-spacer></v-spacer>

                <v-menu offset-y>
                    <template v-slot:activator="{ on }">
                        <v-btn
                            color="primary"
                            icon
                            v-on="on"
                            :loading="loading"
                        >
                            <v-icon>mdi-plus</v-icon>
                        </v-btn>
                    </template>

                    <v-list dense>
                        <v-list-item
                            @click="showHideNewSSHPassword = true"
                            :disabled="!!sshPasswordConn"
                        >
                            <v-list-item-title>SSH (Password)</v-list-item-title>
                        </v-list-item>
                        <v-list-item
                            @click="showHideNewSSHKey = true"
                            :disabled="!!sshKeyConn"
                        >
                            <v-list-item-title>SSH (Key)</v-list-item-title>
                        </v-list-item>
                    </v-list>
                </v-menu>
            </v-card-title>
            <v-divider></v-divider>

            <v-tabs>
                <template v-if="!!sshPasswordConn">
                    <v-tab>
                        SSH (Password)
                        <v-dialog v-model="showHideDeleteSshPassword"
                                  persistent
                                  max-width="40%"
                        >
                            <template v-slot:activator="{ on }">
                                <v-btn
                                    icon
                                    color="error"
                                    @click.stop
                                    @mousedown.stop
                                    @mouseup.stop
                                    v-on="on"
                                >
                                    <v-icon>
                                        mdi-close-circle-outline
                                    </v-icon>
                                </v-btn>
                            </template>

                            <generic-delete-confirmation-form
                                item-name="server connection information"
                                :items-to-delete="[`SSH (Password)`]"
                                :use-global-deletion="false"
                                @do-cancel="showHideDeleteSshPassword = false"
                                @do-delete="onDeleteSshPasswordConnection"
                                :delete-in-progress="deleteInProgress"
                            >
                            </generic-delete-confirmation-form>

                        </v-dialog>
                    </v-tab>
                    <v-tab-item>
                    </v-tab-item>
                </template>

                <template v-if="!!sshKeyConn">
                    <v-tab>SSH (Key)</v-tab>
                    <v-tab-item>
                    </v-tab-item>
                </template>
            </v-tabs>
        </v-card>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import CreateServerSshPasswordConnection from './CreateServerSshPasswordConnection.vue'
import GenericDeleteConfirmationForm from '../../components/dashboard/GenericDeleteConfirmationForm.vue'
import {
    ServerSSHConnectionGeneric,
} from '../../../ts/infrastructure'
import {
    getAllServerConnections, TAllServerConnectionOutput,
    deleteServerSSHPasswordConnection,
} from '../../../ts/api/apiServerConnection'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'

const Props = Vue.extend({
    props: {
        serverId: Number
    }
})

@Component({
    components: {
        CreateServerSshPasswordConnection,
        GenericDeleteConfirmationForm,
    }
})
export default class ServerConnectionManager extends Props {
    loading : boolean = false
    sshPasswordConn : ServerSSHConnectionGeneric | null = null
    sshKeyConn : ServerSSHConnectionGeneric | null = null

    showHideNewSSHPassword : boolean = false
    showHideNewSSHKey : boolean = false

    showHideDeleteSshPassword : boolean = false
    deleteInProgress : boolean = false

    onNewSshPasswordConnection(s : ServerSSHConnectionGeneric) {
        this.sshPasswordConn = s
        this.showHideNewSSHPassword = false
    }

    onDeleteSshPasswordConnection() {
        this.deleteInProgress = true
        deleteServerSSHPasswordConnection({
            orgId: PageParamsStore.state.organization!.Id,
            serverId: this.serverId,
            connectionId: this.sshPasswordConn!.Id,
        }).then(() => {
            this.showHideDeleteSshPassword = false
            this.sshPasswordConn = null
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops. Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        }).finally(() => {
            this.deleteInProgress = false
        })
    }

    refreshData() {
        this.loading = true
        getAllServerConnections({
            orgId: PageParamsStore.state.organization!.Id,
            serverId: this.serverId,
        }).then((resp : TAllServerConnectionOutput) => {
            this.sshPasswordConn = resp.data.SshPassword
            this.sshKeyConn = resp.data.SshKey
            this.loading = false
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops. Something went wrong. Try again.",
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
