import axios from 'axios'
import { deleteFormJson, postFormJson, putFormJson } from '../../http'
import { getAPIRequestConfig } from '../apiUtility'
import { PbcNotificationCadenceSettings } from '../../settings/pbcNotifications'
import {
    apiv2OrgPbcNotificationsSettings,
    apiv2OrgSinglePbcNotificationSetting,
} from '../../url'

export interface TGetPbcNotificationSettingsInput {
    orgId: number
}

export interface TGetPbcNotificationSettingsOutput {
    data: PbcNotificationCadenceSettings[]
}

export function getOrgPbcNotificationSettings(inp : TGetPbcNotificationSettingsInput) : Promise<TGetPbcNotificationSettingsOutput> {
    return axios.get(
        apiv2OrgPbcNotificationsSettings(inp.orgId),
        getAPIRequestConfig()
    )
}

export interface TSinglePbcNotificationSettingInput {
    orgId: number
    notificationId: number
}

export function deletePbcNotificationSetting(inp : TSinglePbcNotificationSettingInput) : Promise<void> {
    return deleteFormJson(
        apiv2OrgSinglePbcNotificationSetting(inp.orgId, inp.notificationId),
        {},
        getAPIRequestConfig()
    )
}

export interface TNewPbcNotificationSettingInput {
    orgId: number
    daysBefore: number
}

export interface TNewPbcNotificationSettingOutput {
    data: PbcNotificationCadenceSettings
}

export function newPbcNotificationSetting(inp : TNewPbcNotificationSettingInput) : Promise<TNewPbcNotificationSettingOutput> {
    return postFormJson(
        apiv2OrgPbcNotificationsSettings(inp.orgId),
        inp,
        getAPIRequestConfig()
    )
}

export interface TEditPbcNotificationSettingInput {
    orgId: number
    setting: PbcNotificationCadenceSettings
    applyAll: boolean
}

export function editPbcNotificationSetting(inp : TEditPbcNotificationSettingInput) : Promise<void> {
    return putFormJson(
        apiv2OrgSinglePbcNotificationSetting(inp.orgId, inp.setting.Id),
        inp,
        getAPIRequestConfig(),
    )
}
