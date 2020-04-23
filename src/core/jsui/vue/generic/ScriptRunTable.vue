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
    createSingleRunLogUrl,
} from '../../ts/url'
import { createUserString } from '../../ts/users'
import { ClientScript } from '../../ts/clientScripts'
import { ManagedCode } from '../../ts/code'
import { VIcon, VProgressCircular } from 'vuetify/lib'
import { standardFormatTime } from '../../ts/time'
import { getClientScriptCodeFromLink, TGetClientScriptCodeFromLinkOutput } from '../../ts/api/apiScripts'
import ScriptBuildRunStatus from './code/ScriptBuildRunStatus.vue'

@Component({
    components: {
        BaseResourceTable,
        HashRenderer,
        ScriptBuildRunStatus
    }
})
export default class ScriptRunTable extends ResourceTableProps {
    linkIdToScript : Record<number, ClientScript> = Object()
    linkIdToCode : Record<number, ManagedCode> = Object()
    linkIdProcessed : Set<number> = new Set<number>()

    get tableHeaders() : any[] {
        return [
            {
                text: 'Id',
                value: 'id',
            },
            {
                text: 'Script',
                value: 'script',
            },
            {
                text: 'User',
                value: 'user',
            },
            {
                text: 'Request Time',
                value: 'start',
            },
            {
                text: 'Build Status',
                value: 'build',
            },
            {
                text: 'Run Status',
                value: 'run',
            },
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    transformInputResourceToTableItem(inp : any) : any {
        return {
            value: inp,
            start: standardFormatTime(inp.StartTime),
            user: createUserString(MetadataStore.getters.getUser(inp.UserId)),
        }
    }

    goToRun(item : any) {
        window.location.assign(createSingleRunLogUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            item.value.Id
        ))
    }

    retrieveInfoFromLinkId(linkId : number) {
        if (this.linkIdProcessed.has(linkId)) {
            return
        }

        this.linkIdProcessed.add(linkId)
        getClientScriptCodeFromLink({
            orgId: PageParamsStore.state.organization!.Id,
            linkId: linkId,
        }).then((resp : TGetClientScriptCodeFromLinkOutput) => {
            Vue.set(this.linkIdToScript, linkId, resp.data.Script)
            Vue.set(this.linkIdToCode, linkId, resp.data.Code)
        }).catch((err: any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    renderId(props : any) : VNode {
        return this.$createElement(
            'span',
            {
                class: ['font-weight-bold']
            },
            `#${props.item.value.Id}`
        )
    }

    renderScriptName(props : any) : VNode {
        this.retrieveInfoFromLinkId(props.item.value.LinkId)

        let linkId = props.item.value.LinkId
        if (!(linkId in this.linkIdToCode) || !(linkId in this.linkIdToScript)) {
            return this.$createElement(
                'span',
                'Loading...'
            )
        }

        let script = this.linkIdToScript[linkId]
        let code = this.linkIdToCode[linkId]

        return this.$createElement(
            'a',
            {
                style: {
                    display: 'flex',
                    'align-items': 'center',
                    'justify-items': 'center',
                },
                attrs: {
                    href: createSingleScriptUrl(PageParamsStore.state.organization!.OktaGroupName, script.Id),
                }
            },
            [
                this.$createElement(
                    'p',
                    {
                        class: ['ma-0'],
                    },
                    `${script.Name}`
                ),
                this.$createElement(
                    HashRenderer,
                    {
                        props: {
                            hash: code.GitHash,
                        },
                        class: ['ml-2'],
                    },
                ),
            ],
        )
    }

    renderStatus(required : boolean, success: boolean, start : Date | null, end : Date | null, forceFail : boolean) : VNode {
        return this.$createElement(
            ScriptBuildRunStatus,
            {
                props: {
                    success: success,
                    start: start,
                    end: end,
                    showTimeStamp: true,
                    notRequired: !required,
                    forceFail: forceFail
                }
            }
        )
    }

    renderBuildStatus(props : any) : VNode {
        return this.renderStatus(
            props.item.value.RequiresBuild,
            props.item.value.BuildSuccess,
            props.item.value.BuildStartTime,
            props.item.value.BuildFinishTime,
            false,
        )
    }

    renderRunStatus(props : any) : VNode {
        return this.renderStatus(
            true,
            props.item.value.RunSuccess,
            props.item.value.RunStartTime,
            props.item.value.RunFinishTime,
            !props.item.value.BuildSuccess && props.item.value.RequiresBuild,
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
                    resourceName: "script runs",
                    showExpand: false
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any) => this.$emit('delete', item.value),
                    'click:row': this.goToRun
                },
                scopedSlots: {
                    'item.id': this.renderId,
                    'item.script': this.renderScriptName,
                    'item.build': this.renderBuildStatus,
                    'item.run': this.renderRunStatus,
                }
            }
        )
    }

}

</script>
