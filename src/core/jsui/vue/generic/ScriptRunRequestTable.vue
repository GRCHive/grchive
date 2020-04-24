<script lang="ts">

import Vue, { VNode } from 'vue'
import GenericRequestTable from './GenericRequestTable.vue'
import Component from 'vue-class-component'
import { ClientScript } from '../../ts/clientScripts'
import { ManagedCode } from '../../ts/code'
import { getGenericRequestScript, TGetGenericRequestScriptOutput } from '../../ts/api/apiRequests'
import { contactUsUrl, createSingleScriptUrl, createSingleScriptRequestUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'
import BaseResourceTable from './BaseResourceTable.vue'
import HashRenderer from './code/HashRenderer.vue'

@Component({
    components: {
        BaseResourceTable,
        HashRenderer
    }
})
export default class ScriptRunRequestTable extends GenericRequestTable {
    idToScript : Record<number, ClientScript> = Object()
    idToCode: Record<number, ManagedCode> = Object()
    pendingScriptRequests : Set<number> = new Set<number>()

    get tableHeaders() : any[] {
        //@ts-ignore
        let headers = GenericRequestTable.options.computed.tableHeaders.get()
        headers.unshift(...[
            {
                text: 'Script',
                value: 'script',
            },
            {
                text: 'Version',
                value: 'version',
            },
        ])
        return headers
    }

    goToScriptRequest(item : any) {
        window.location.assign(createSingleScriptRequestUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            item.value.Id))
    }

    loadAuxData(reqId : number) {
        if (this.pendingScriptRequests.has(reqId)) {
            return
        }

        this.pendingScriptRequests.add(reqId)
        getGenericRequestScript({
            requestId: reqId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetGenericRequestScriptOutput) => {
            Vue.set(this.idToScript, reqId, resp.data.Script)
            Vue.set(this.idToCode, reqId, resp.data.Code)
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
        ret.script = this.idToScript[inp.Id]
        ret.version = this.idToCode[inp.Id]
        return ret
    }

    renderScript(props : any) : VNode {
        if (!props.item.script) {
            return this.$createElement('span', 'Loading...')
        }
        return this.$createElement(
            'a',
            {
                attrs: {
                    href: createSingleScriptUrl(PageParamsStore.state.organization!.OktaGroupName, props.item.script.Id)
                }
            },
            props.item.script.Name
        )
    }

    renderVersion(props : any) : VNode {
        if (!props.item.version) {
            return this.$createElement('span', 'Loading...')
        }
        return this.$createElement(
            HashRenderer,
            {
                props: {
                    hash: props.item.version.GitHash,
                }
            })
    }

    render() : VNode {
        return this.$createElement(
            BaseResourceTable,
            {
                props: {
                    ...this.$props,
                    tableHeaders: this.tableHeaders,
                    tableItems: this.tableItems,
                    resourceName: "script request",
                    showExpand: true
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any) => this.$emit('delete', item.value),
                    'click:row': this.goToScriptRequest
                },
                scopedSlots: {
                    'item.approval': this.renderApproval,
                    'item.requestTime': this.renderRequestTime,
                    'item.dueDate': this.renderDueDate,
                    'item.script': this.renderScript,
                    'item.version': this.renderVersion,
                    'expanded-item': this.renderExpansion,
                }
            }
        )
    }
}

</script>
