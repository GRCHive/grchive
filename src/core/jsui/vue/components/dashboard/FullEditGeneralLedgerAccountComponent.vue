<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-breadcrumbs :items="parentBreadcrumbs" class="pa-0">
            </v-breadcrumbs>
            <v-list-item two-line class="pa-0">
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Account: {{ glAccount.AccountName }}
                        <span class="subtitle-1" v-if="glAccount.AccountName != glAccount.AccountId">
                            ({{ glAccount.AccountId }})
                        </span>
                        <v-btn icon @click="expandDescription = !expandDescription">
                            <v-icon small v-if="!expandDescription" >mdi-chevron-down</v-icon>
                            <v-icon small v-else>mdi-chevron-up</v-icon>
                        </v-btn>
                    </v-list-item-title>

                    <v-list-item-subtitle :class="expandDescription ? `long-text` : `hide-long-text`">
                        {{ glAccount.AccountDescription }}
                    </v-list-item-subtitle>
                </v-list-item-content>
            </v-list-item>
            <v-divider></v-divider>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { GeneralLedgerAccount, GeneralLedgerCategory, GeneralLedger } from '../../../ts/generalLedger'
import { TGetGLAccountInputs, TGetGLAccountOutputs, getGLAccount } from '../../../ts/api/apiGeneralLedger'
import { PageParamsStore } from '../../../ts/pageParams'

@Component({
    components: {
    }
})
export default class FullEditGeneralLedgerAccountComponent extends Vue {
    ready: boolean = false
    ledger : GeneralLedger = new GeneralLedger()

    get parentBreadcrumbs() : any[] {
        let parentCrumbs = []
        let currentParent : GeneralLedgerCategory | null = this.glAccount.ParentCategory
        while (currentParent != null) {
            parentCrumbs.push({
                disabled: true,
                text: currentParent.Name
            })
            currentParent = currentParent.ParentCategory
        }
        return parentCrumbs
    }

    get glAccount() : GeneralLedgerAccount {
        return this.ledger.accounts.values().next().value
    }

    refreshAccountData() {
        let data = window.location.pathname.split('/')
        let accId = Number(data[data.length - 1])

        getGLAccount(<TGetGLAccountInputs>{
            orgId: PageParamsStore.state.organization!.Id,
            accId: accId,
        }).then((resp : TGetGLAccountOutputs) => {
            console.log(resp.data)
            this.ledger.rebuildGL(resp.data.Parents, [resp.data.Account])
            this.ready = true
        }).catch((err : any) => {
            console.log(err)
            //window.location.replace('/404')
        })
    }

    mounted() {
        this.refreshAccountData()
    }
}

</script>
