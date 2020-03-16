<template>
    <div>
        <v-list dense class="pa-0">
            <v-list-item class="pa-0">
                <v-list-item-action class="ma-0">
                    <v-btn icon @click="minimizeSystem = !minimizeSystem">
                        <v-icon small>
                            {{ !minimizeSystem ? "mdi-chevron-up" : "mdi-chevron-down" }}
                        </v-icon>
                    </v-btn>
                </v-list-item-action>

                <v-subheader class="flex-grow-1 pr-0">
                    LINKED SYSTEMS
                </v-subheader>

                <v-list-item-action class="ma-0">
                    <v-dialog persistent max-width="40%" v-model="showLinkSystem">
                        <template v-slot:activator="{ on }">
                            <v-btn
                                icon
                                v-on="on"
                            >
                                <v-icon small>
                                    mdi-plus
                                </v-icon>
                            </v-btn>
                        </template>

                        <v-card>
                            <v-card-title>
                                Link System
                            </v-card-title>
                            <v-divider></v-divider>

                            <system-table-with-controls
                                class="ma-4"
                                v-model="systemsToLink"
                                :exclude="linkedSystems"
                                disable-new
                                disable-delete
                                enable-select
                            >
                            </system-table-with-controls>

                            <v-card-actions>
                                <v-btn
                                    color="error"
                                    @click="cancelSystemLink"
                                >
                                    Cancel
                                </v-btn>
                                <v-spacer></v-spacer>
                                <v-btn
                                    color="success"
                                    @click="saveSystemLink"
                                    :disabled="systemsToLink.length == 0"
                                >
                                    Link
                                </v-btn>
                            </v-card-actions>
                        </v-card>
                    </v-dialog>
                </v-list-item-action>
            </v-list-item>
        </v-list>
        <systems-table
            :resources="linkedSystems"
            use-crud-delete
            @delete="deleteLinkedSystem"
            v-if="!minimizeSystem"
        >
        </systems-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import VueSetup from '../../../../ts/vueSetup' 
import SystemsTable from '../../../generic/SystemsTable.vue'
import SystemTableWithControls from '../../../generic/resources/SystemTableWithControls.vue'
import { PageParamsStore } from '../../../../ts/pageParams'
import { System } from '../../../../ts/systems'
import { contactUsUrl } from '../../../../ts/url'
import { 
    newNodeSystemLink,
    deleteNodeSystemLink
} from '../../../../ts/api/apiNodeSystemLinks'

@Component({
    components: {
        SystemsTable,
        SystemTableWithControls,
    }
})
export default class NodeLinkedSystemsEditor extends Vue {
    systemsToLink : System[] = []
    showLinkSystem : boolean = false
    minimizeSystem: boolean = false

    cancelSystemLink() {
        this.systemsToLink = []
        this.showLinkSystem = false
    }

    saveSystemLink() {
        if (this.systemsToLink.length == 0) {
            return
        }

        let nodeId : number = this.currentNode.Id
        let system : System = this.systemsToLink[0]
        newNodeSystemLink({
            nodeId: nodeId,
            systemId: system.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then(() => {
            VueSetup.store.commit('addNodeSystemLink', {
                nodeId: nodeId,
                system: system,
            })
            this.systemsToLink = []
            this.showLinkSystem = false
        }).catch((err : any) => {
            //@ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong, please reload the page and try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    deleteLinkedSystem(system : System) {
        let id : number = this.currentNode.Id
        deleteNodeSystemLink({
            nodeId: id,
            systemId: system.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then(() => {
            VueSetup.store.commit('deleteNodeSystemLink', {
                nodeId: id,
                systemId: system.Id,
            })
        }).catch((err : any) => {
            //@ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong, please reload the page and try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    get linkedSystems() : System[] | null {
        return VueSetup.store.getters.systemsLinkedToNode(this.currentNode.Id)
    }

    get currentNode() : ProcessFlowNode {
        return VueSetup.store.getters.currentNodeInfo
    }
}

</script>
