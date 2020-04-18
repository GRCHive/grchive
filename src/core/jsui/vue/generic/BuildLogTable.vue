<script lang="ts">

import Vue, { VNode } from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import MetadataStore from '../../ts/metadata'
import HashRenderer from './code/HashRenderer.vue'
import { PageParamsStore } from '../../ts/pageParams'
import {
    contactUsUrl,
    createSingleScriptUrl,
    createSingleClientDataUrl,
} from '../../ts/url'
import { createUserString } from '../../ts/users'
import { ClientScript } from '../../ts/clientScripts'
import { ClientData } from '../../ts/clientData'
import { DroneCiStatus } from '../../ts/code'
import { VIcon, VProgressCircular } from 'vuetify/lib'
import { standardFormatTime } from '../../ts/time'
import { getCodeLink, TGetCodeLinkOutput } from '../../ts/api/apiCode'
import CodeBuildStatus from './code/CodeBuildStatus.vue'

@Component({
    components: {
        BaseResourceTable,
        HashRenderer,
        CodeBuildStatus,
    }
})
export default class ScriptRunTable extends ResourceTableProps {
    codeIdToObject: Record<number, ClientScript | ClientData> = Object()
    codeIdIsScript : Record<number, boolean> = Object()
    codeIdProcessed : Set<number> = new Set<number>()

    get tableHeaders() : any[] {
        return [
            {
                text: 'Id',
                value: 'id',
            },
            {
                text: 'Data/Script',
                value: 'object',
            },
            {
                text: 'Modified Time',
                value: 'time',
            },
            {
                text: 'Build Status',
                value: 'build',
            },
            {
                text: 'User',
                value: 'user',
            },
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    transformInputResourceToTableItem(inp : any) : any {
        return {
            value: inp,
            time: standardFormatTime(inp.ActionTime),
            user: createUserString(MetadataStore.getters.getUser(inp.UserId)),
        }
    }

    goToLogs(item : any) {
    }

    renderId(props : any) : VNode {
        return this.$createElement(
            HashRenderer,
            {
                props: {
                    hash: props.item.value.GitHash,
                }
            },
        )
    }

    retrieveAuxInfo(codeId: number) {
        if (this.codeIdProcessed.has(codeId)) {
            return
        }

        this.codeIdProcessed.add(codeId)
        getCodeLink({
            orgId: PageParamsStore.state.organization!.Id,
            codeId: codeId,
        }).then((resp : TGetCodeLinkOutput) => {
            Vue.set(this.codeIdToObject, codeId, resp.data.Data || resp.data.Script!)
            Vue.set(this.codeIdIsScript, codeId, resp.data.Script !== null)
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

    renderData(s : ClientData) : VNode {
        return this.$createElement(
            'a',
            {
                attrs: {
                    href: createSingleClientDataUrl(PageParamsStore.state.organization!.OktaGroupName, s.Id),
                }
            },
            [
                this.$createElement(
                    'span',
                    {
                        class: ['font-weight-bold'],
                    },
                    'Data: ',
                ),
                s.Name,
            ]
        )
    }

    renderScript(s : ClientScript) : VNode {
        return this.$createElement(
            'a',
            {
                attrs: {
                    href: createSingleScriptUrl(PageParamsStore.state.organization!.OktaGroupName, s.Id),
                }
            },
            [
                this.$createElement(
                    'span',
                    {
                        class: ['font-weight-bold'],
                    },
                    'Script: ',
                ),
                s.Name
            ]
        )
    }

    renderDataScript(props : any) : VNode {
        this.retrieveAuxInfo(props.item.value.Id)

        let codeId = props.item.value.Id
        if (!(codeId in this.codeIdToObject)) {
            return this.$createElement(
                'span',
                'Loading...'
            )
        }

        let obj = this.codeIdToObject[codeId]
        let isScript = this.codeIdIsScript[codeId]

        if (isScript) {
            return this.renderScript(obj as ClientScript)
        } else {
            return this.renderData(obj as ClientData)
        }
    }

    renderBuildStatus(props : any) : VNode {
        this.retrieveAuxInfo(props.item.value.Id)
        return this.$createElement(
            CodeBuildStatus,
            {
                props: {
                    commit: props.item.value.GitHash,
                    showTimeStamp: true,
                }
            }
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
                    resourceName: "build logs",
                    showExpand: false
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any) => this.$emit('delete', item.value),
                    'click:row': this.goToLogs
                },
                scopedSlots: {
                    'item.id': this.renderId,
                    'item.object': this.renderDataScript,
                    'item.build': this.renderBuildStatus,
                }
            }
        )
    }

}

</script>
