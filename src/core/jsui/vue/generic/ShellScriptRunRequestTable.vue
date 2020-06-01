<script lang="ts">

import Vue, { VNode } from 'vue'
import GenericRequestTable from './GenericRequestTable.vue'
import Component from 'vue-class-component'
import { contactUsUrl, createSingleShellRequestUrl, createSingleShellUrl} from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'
import {
    getGenericRequestShell, TGetGenericRequestShellOutput,
} from '../../ts/api/apiRequests'
import {
    ShellScript, ShellScriptVersion
} from '../../ts/shell'
import BaseResourceTable from './BaseResourceTable.vue'
import { VChip } from 'vuetify/lib'

@Component({
    components: {
        BaseResourceTable,
    }
})
export default class ScriptRunRequestTable extends GenericRequestTable {
    pendingScriptRequests : Set<number> = new Set<number>()
    idToShell : Record<number, ShellScript> = Object()
    idToVersion : Record<number, number> = Object()

    get tableHeaders() : any[] {
        //@ts-ignore
        let headers = GenericRequestTable.options.computed.tableHeaders.get()
        headers.unshift(...[
            {
                text: 'Shell Script',
                value: 'shell',
            },
            {
                text: 'Version',
                value: 'version',
            },
        ])
        return headers
    }

    goToShellScriptRequest(item : any) {
        window.location.assign(createSingleShellRequestUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            item.value.Id))
    }

    loadAuxData(reqId : number) {
        if (this.pendingScriptRequests.has(reqId)) {
            return
        }

        this.pendingScriptRequests.add(reqId)
        getGenericRequestShell({
            requestId: reqId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetGenericRequestShellOutput) => {
            Vue.set(this.idToShell, reqId, resp.data.Shell)
            Vue.set(this.idToVersion, reqId, resp.data.VersionNum)
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

    transformInputResourceToTableItem(inp : any) : any {
        this.loadAuxData(inp.Id)
        //@ts-ignore
        let ret = GenericRequestTable.options.methods.transformInputResourceToTableItem.call(this, inp)
        ret.shell = this.idToShell[inp.Id]
        ret.version = this.idToVersion[inp.Id]
        return ret
    }

    renderShell(props : any) : VNode {
        if (!props.item.shell) {
            return this.$createElement('span', 'Loading...')
        }
        return this.$createElement(
            'a',
            {
                attrs: {
                    href: createSingleShellUrl(
                        PageParamsStore.state.organization!.OktaGroupName,
                        props.item.shell.Id
                    ) + `?version=${props.item.version}`
                }
            },
            props.item.shell.Name
        )
    }

    renderVersion(props : any) : VNode {
        if (!props.item.version) {
            return this.$createElement('span', 'Loading...')
        }
        return this.$createElement(
            VChip,
            {
                props: {
                    pill: true
                },
            },
            [
                `v${props.item.version}`
            ]
        )
    }

    render() : VNode {
        return this.$createElement(
            BaseResourceTable,
            {
                props: {
                    ...this.$props,
                    tableHeaders: this.tableHeaders,
                    tableItems: this.tableItems,
                    resourceName: "shell script run request",
                    showExpand: true
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any) => this.$emit('delete', item.value),
                    'click:row': this.goToShellScriptRequest
                },
                scopedSlots: {
                    'item.approval': this.renderApproval,
                    'item.requestTime': this.renderRequestTime,
                    'item.dueDate': this.renderDueDate,
                    'item.shell': this.renderShell,
                    'item.version': this.renderVersion,
                    'expanded-item': this.renderExpansion,
                }
            }
        )
    }
}

</script>
