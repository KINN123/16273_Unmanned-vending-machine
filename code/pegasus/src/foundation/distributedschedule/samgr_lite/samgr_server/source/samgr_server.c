/*
 * Copyright (c) 2020 Huawei Device Co., Ltd.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#include "samgr_server.h"
#include <fcntl.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <ohos_init.h>
#include <ohos_types.h>
#include <ohos_errno.h>
#include <liteipc.h>
#include <liteipc_adapter.h>
#include <log.h>
#include "cJSON.h"
#include "policy_define.h"
#include "samgr_lite.h"
#include "memory_adapter.h"
#include "thread_adapter.h"
#include "default_client.h"

#undef LOG_TAG
#undef LOG_DOMAIN
#define LOG_TAG "Samgr"
#define LOG_DOMAIN 0xD001800

typedef int(*ProcFunc)(SamgrServer *server, int32 option, void *origin, IpcIo *req, IpcIo *reply);
#define MAX_SA_SIZE 0x100
#define RETRY_TIMES 3
#define RETRY_INTERVAL 1
#define UID_HAP 10000
#define MAX_SYSCAP_NUM_PER_REPLY 118

static const char *GetName(Service *service);
static BOOL Initialize(Service *service, Identity identity);
static TaskConfig GetTaskConfig(Service *service);
static BOOL MessageHandle(Service *service, Request *request);
static int32 Invoke(IServerProxy *iProxy, int funcId, void *origin, IpcIo *req, IpcIo *reply);
static int OnEndpointExit(const IpcContext *context, void* ipcMsg, IpcIo* data, void* argv);
static int ProcEndpoint(SamgrServer *server, int32 option, void *origin, IpcIo *req, IpcIo *reply);
static int32 ProcPutFeature(SamgrServer *server, const void *origin, IpcIo *req, IpcIo *reply, SvcIdentity *identity);
static int32 ProcGetFeature(SamgrServer *server, const void *origin, IpcIo *req, IpcIo *reply, SvcIdentity *identity);
static int ProcFeature(SamgrServer *server, int32 option, void *origin, IpcIo *req, IpcIo *reply);
static int RegisterSamgrEndpoint(const IpcContext* context, SvcIdentity* identity);
static void TransmitPolicy(int ret, const SvcIdentity* identity, IpcIo *reply,
                           const PolicyTrans *policy, uint32 policyNum);
static void TransmitFixedPolicy(IpcIo *reply, PolicyTrans policy);
static IpcAuthInterface *GetIpcAuthInterface(void);
static int ProcSysCap(SamgrServer *server, int32 option, void *origin, IpcIo *req, IpcIo *reply);
static void ParseSysCap(void);

static SamgrServer g_server = {
    .GetName = GetName,
    .Initialize = Initialize,
    .GetTaskConfig = GetTaskConfig,
    .MessageHandle = MessageHandle,
    SERVER_IPROXY_IMPL_BEGIN,
    .Invoke = Invoke,
    IPROXY_END,
};

static ProcFunc g_functions[] = {
    [RES_ENDPOINT] = ProcEndpoint,
    [RES_FEATURE] = ProcFeature,
    [RES_SYSCAP] = ProcSysCap,
};

static const char *GetSysCapName(const SysCapImpl *serviceImpl)
{
    if (serviceImpl == NULL) {
        return NULL;
    }
    return serviceImpl->name;
}

static void InitializeRegistry(void)
{
    HILOG_INFO(HILOG_MODULE_SAMGR, "Initialize Registry!");
    g_server.mtx = MUTEX_InitValue();
    SASTORA_Init(&g_server.store);
    g_server.samgr = SAMGR_CreateEndpoint("samgr", RegisterSamgrEndpoint);
    SAMGR_GetInstance()->RegisterService((Service *)&g_server);
    g_server.sysCapMtx = MUTEX_InitValue();
    g_server.sysCapabilitys = VECTOR_Make((VECTOR_Key)GetSysCapName, (VECTOR_Compare)strcmp);
    ParseSysCap();
    HILOG_INFO(HILOG_MODULE_SAMGR, "InitializeRegistry ParseSysCap size: %d", VECTOR_Size(&(g_server.sysCapabilitys)));
}
SYS_SERVICE_INIT(InitializeRegistry);

static BOOL CanRequest(const void *origin)
{
    pid_t uid = GetCallingUid(origin);
    return uid < UID_HAP;
}

static const char *GetName(Service *service)
{
    (void)service;
    return SAMGR_SERVICE;
}

static BOOL Initialize(Service *service, Identity identity)
{
    SamgrServer *server = (SamgrServer *)service;
    server->identity = identity;
    SaName saName = {SAMGR_SERVICE, NULL};
    SAMGR_AddRouter(server->samgr, &saName, &server->identity, GET_IUNKNOWN(*server));
    return TRUE;
}

static BOOL MessageHandle(Service *service, Request *request)
{
    SamgrServer *server = (SamgrServer *)service;
    switch (request->msgId) {
        case MSG_CLEAN:
            MUTEX_Lock(server->mtx);
            SASTORA_ClearByPid(&server->store, request->msgValue);
            MUTEX_Unlock(server->mtx);
            break;
        default:
            break;
    }
    return TRUE;
}

static TaskConfig GetTaskConfig(Service *service)
{
    (void)service;
    TaskConfig config = {LEVEL_HIGH, PRI_BUTT - 1, 0x400, 20, SINGLE_TASK}; // Cannot use PRI_BUTT directly, so minus 1
    return config;
}

static int32 Invoke(IServerProxy *iProxy, int funcId, void *origin, IpcIo *req, IpcIo *reply)
{
    SamgrServer *server = GET_OBJECT(iProxy, SamgrServer, iUnknown);
    int32 resource = IpcIoPopUint32(req);
    int32 option = IpcIoPopUint32(req);
    if (server == NULL || resource >= RES_BUTT || resource < 0 || g_functions[resource] == NULL) {
        HILOG_ERROR(HILOG_MODULE_SAMGR, "Invalid Msg<%d, %d, %d>", resource, option, funcId);
        return EC_INVALID;
    }
    return g_functions[resource](server, option, origin, req, reply);
}

static int ProcEndpoint(SamgrServer *server, int32 option, void *origin, IpcIo *req, IpcIo *reply)
{
    if (option != OP_POST) {
        IpcIoPushUint32(reply, (uint32)INVALID_INDEX);
        return EC_FAILURE;
    }

    pid_t pid = GetCallingPid(origin);
    PidHandle handle;
    MUTEX_Lock(server->mtx);
    int index = SASTORA_FindHandleByPid(&g_server.store, pid, &handle);
    if (index == INVALID_INDEX) {
        SvcIdentity identity = {(uint32)INVALID_INDEX, (uint32)INVALID_INDEX, (uint32)INVALID_INDEX};
        (void)GenServiceHandle(&identity, GetCallingTid(origin));
#ifdef __LINUX__
        IpcMsg* data = (IpcMsg*)origin;
        if (data == NULL) {
            HILOG_ERROR(HILOG_MODULE_SAMGR, "Register Endpoint origin null pointer!");
            return EC_FAILURE;
        }
        identity.handle = data->target.handle;
        BinderAcquire(g_server.samgr->context, identity.handle);
#endif

        handle.pid = pid;
        handle.uid = GetCallingUid(origin);
        handle.handle = identity.handle;
        handle.deadId = INVALID_INDEX;
        (void)SASTORA_SaveHandleByPid(&server->store, handle);
        (void)UnregisterDeathCallback(identity, handle.deadId);
        (void)RegisterDeathCallback(server->samgr->context, identity, OnEndpointExit, (void*)((uintptr_t)pid),
                                    &handle.deadId);
    }
    MUTEX_Unlock(server->mtx);
    IpcIoPushUint32(reply, handle.handle);
    HILOG_INFO(HILOG_MODULE_SAMGR, "Register Endpoint<%d, %u, %u>", handle.pid, handle.handle, handle.deadId);
    return EC_SUCCESS;
}

static int32 ProcPutFeature(SamgrServer *server, const void *origin, IpcIo *req, IpcIo *reply, SvcIdentity *identity)
{
    size_t len = 0;
    char *service = (char *)IpcIoPopString(req, &len);
    if (service == NULL || len == 0) {
        IpcIoPushInt32(reply, EC_INVALID);
        return EC_INVALID;
    }
    pid_t pid = GetCallingPid(origin);
    uid_t uid = GetCallingUid(origin);
    char *feature = IpcIoPopBool(req) ? NULL : (char *)IpcIoPopString(req, &len);
    MUTEX_Lock(server->mtx);
    PidHandle handle;
    int index = SASTORA_FindHandleByUidPid(&server->store, uid, pid, &handle);
    if (index == INVALID_INDEX) {
        MUTEX_Unlock(server->mtx);
        HILOG_ERROR(HILOG_MODULE_SAMGR, "Endpoint[%d] is not register", pid);
        IpcIoPushInt32(reply, EC_NOSERVICE);
        return EC_NOSERVICE;
    }
    *identity = SASTORA_Find(&server->store, service, feature);
    if (identity->handle != INVALID_INDEX && identity->handle != handle.handle) {
        MUTEX_Unlock(server->mtx);
        IpcIoPushInt32(reply, EC_INVALID);
        return EC_INVALID;
    }
    identity->token = IpcIoPopUint32(req);
    identity->handle = handle.handle;

    PolicyTrans *policy = NULL;
    RegParams regParams = {service, feature, handle.uid, handle.pid};
    uint32 policyNum = 0;
    int ret = g_server.ipcAuth->GetCommunicationStrategy(regParams, &policy, &policyNum);
    if (ret != EC_SUCCESS || policy == NULL) {
        MUTEX_Unlock(server->mtx);
        SAMGR_Free(policy);
        HILOG_DEBUG(HILOG_MODULE_SAMGR, "Remote Get Communication Strategy<%s, %s> No Permission<%d>!",
                    service, feature, ret);
        IpcIoPushInt32(reply, EC_PERMISSION);
        return EC_PERMISSION;
    }

    ret = SASTORA_Save(&server->store, service, feature, identity);
    MUTEX_Unlock(server->mtx);
    HILOG_DEBUG(HILOG_MODULE_SAMGR, "Register Feature<%s, %s> pid<%d>, id<%u, %u> ret:%d",
                service, feature, pid, identity->handle, identity->token, ret);
    TransmitPolicy(ret, identity, reply, policy, policyNum);
    SAMGR_Free(policy);
    return ret;
}

static void TransmitPolicy(int ret, const SvcIdentity* identity, IpcIo *reply,
                           const PolicyTrans *policy, uint32 policyNum)
{
    if (identity == NULL || reply == NULL || policy == NULL) {
        IpcIoPushInt32(reply, EC_INVALID);
        return;
    }
    if (ret != EC_SUCCESS) {
        IpcIoPushInt32(reply, ret);
        return;
    }
    IpcIoPushInt32(reply, ret);
    IpcIoPushSvc(reply, identity);
    IpcIoPushUint32(reply, policyNum);
    uint32 i;
    for (i = 0; i < policyNum; i++) {
        IpcIoPushInt32(reply, policy[i].type);
        switch (policy[i].type) {
            case RANGE:
                IpcIoPushInt32(reply, policy[i].uidMin);
                IpcIoPushInt32(reply, policy[i].uidMax);
                break;
            case FIXED:
                TransmitFixedPolicy(reply, policy[i]);
                break;
            case BUNDLENAME:
                IpcIoPushInt32(reply, policy[i].fixedUid[0]);
                break;
            default:
                break;
        }
    }
}

static void TransmitFixedPolicy(IpcIo *reply, PolicyTrans policy)
{
    if (reply == NULL) {
        return;
    }
    uint32 i;
    for (i = 0; i < UID_SIZE; i++) {
        IpcIoPushInt32(reply, policy.fixedUid[i]);
    }
}

static int32 ProcGetFeature(SamgrServer *server, const void *origin, IpcIo *req, IpcIo *reply, SvcIdentity *identity)
{
    size_t len = 0;
    char *service = (char *)IpcIoPopString(req, &len);
    if (service == NULL || len == 0) {
        IpcIoPushInt32(reply, EC_INVALID);
        return EC_INVALID;
    }
    char *feature = IpcIoPopBool(req) ? NULL : (char *)IpcIoPopString(req, &len);
    MUTEX_Lock(server->mtx);
    *identity = SASTORA_Find(&server->store, service, feature);
    if (identity->handle == INVALID_INDEX) {
        MUTEX_Unlock(server->mtx);
        HILOG_DEBUG(HILOG_MODULE_SAMGR, "Cannot Find Feature<%s, %s> id<%u, %u> ret:%d",
                    service, feature, identity->handle, identity->token, EC_NOSERVICE);
        return EC_NOSERVICE;
    }

    PidHandle providerPid = SASTORA_FindPidHandleByIpcHandle(&server->store, identity->handle);
    MUTEX_Unlock(server->mtx);
    if (providerPid.pid == INVALID_INDEX || providerPid.uid == INVALID_INDEX) {
        HILOG_DEBUG(HILOG_MODULE_SAMGR, "Cannot Find PidHandle<%s, %s> id<%d, %d> ret:%d",
                    service, feature, identity->handle, identity->token, EC_FAILURE);
        return EC_FAILURE;
    }
    AuthParams authParams = {
        .providerService = service,
        .providerfeature = feature,
        .consumerPid = GetCallingPid(origin),
        .consumerUid = GetCallingUid(origin),
        .providerPid = providerPid.pid,
        .providerUid = providerPid.uid
    };
    int isAuth = g_server.ipcAuth->IsCommunicationAllowed(authParams);
    HILOG_DEBUG(HILOG_MODULE_SAMGR, "Judge Auth<%s, %s> ret:%d", service, feature, isAuth);

    int32 ret = (isAuth == EC_SUCCESS) ? AddServiceAccess(*identity, GetCallingTid(origin)) : EC_PERMISSION;
    HILOG_DEBUG(HILOG_MODULE_SAMGR, "Find Feature<%s, %s> id<%d, %d> ret:%d",
                service, feature, identity->handle, identity->token, ret);
    return ret;
}

static int ProcFeature(SamgrServer *server, int32 option, void *origin, IpcIo *req, IpcIo *reply)
{
    if (option != OP_PUT && option != OP_GET) {
        IpcIoPushInt32(reply, EC_INVALID);
        return EC_INVALID;
    }

    if (g_server.ipcAuth == NULL) {
        g_server.ipcAuth = GetIpcAuthInterface();
    }
    if (g_server.ipcAuth == NULL) {
        IpcIoPushInt32(reply, EC_NOINIT);
        return EC_NOINIT;
    }

    int ret = EC_SUCCESS;
    SvcIdentity identity = {INVALID_INDEX, INVALID_INDEX, INVALID_INDEX};
    if (option == OP_PUT) {
        ret = ProcPutFeature(server, origin, req, reply, &identity);
    }
    if (ret != EC_SUCCESS) {
        return ret;
    }

    if (option == OP_GET) {
        ret = ProcGetFeature(server, origin, req, reply, &identity);
        IpcIoPushInt32(reply, ret);
        if (ret == EC_SUCCESS) {
            IpcIoPushSvc(reply, &identity);
        }
    }
    return ret;
}

static int32 ProcAddSysCap(SamgrServer *server, IpcIo *req)
{
    size_t len = 0;
    char *sysCap = (char *)IpcIoPopString(req, &len);
    if (sysCap == NULL || len == 0 || len > MAX_SYSCAP_NAME_LEN) {
        HILOG_ERROR(HILOG_MODULE_SAMGR, "ProcAddSysCap sysCap invalid");
        return EC_INVALID;
    }
    MUTEX_Lock(server->sysCapMtx);
    Vector *sysCapablitys = &(server->sysCapabilitys);
    int16 pos = VECTOR_FindByKey(sysCapablitys, (void *)sysCap);
    if (pos < 0) {
        MUTEX_Unlock(server->sysCapMtx);
        return EC_FAILURE;
    }
    SysCapImpl *serviceImpl = (SysCapImpl *)VECTOR_At(sysCapablitys, pos);
    if (serviceImpl == NULL || serviceImpl->name == NULL) {
        MUTEX_Unlock(server->sysCapMtx);
        return EC_FAILURE;
    }
    serviceImpl->isRegister = TRUE;
    MUTEX_Unlock(server->sysCapMtx);
    return EC_SUCCESS;
}

static BOOL ProcGetSysCap(const SamgrServer *server, IpcIo *req)
{
    size_t len = 0;
    char *sysCap = (char *)IpcIoPopString(req, &len);
    if (sysCap == NULL || len == 0 || len > MAX_SYSCAP_NAME_LEN) {
        HILOG_ERROR(HILOG_MODULE_SAMGR, "ProcGetSysCap sysCap invalid");
        return FALSE;
    }
    MUTEX_Lock(server->sysCapMtx);
    Vector *sysCapablitys = &(server->sysCapabilitys);
    int16 pos = VECTOR_FindByKey(sysCapablitys, (void *)sysCap);
    if (pos < 0) {
        MUTEX_Unlock(server->sysCapMtx);
        return FALSE;
    }
    SysCapImpl *serviceImpl = (SysCapImpl *)VECTOR_At(sysCapablitys, pos);
    if (serviceImpl == NULL) {
        MUTEX_Unlock(server->sysCapMtx);
        return FALSE;
    }

    BOOL res = (serviceImpl->isRegister == TRUE);
    MUTEX_Unlock(server->sysCapMtx);
    return res;
}

static int32 GetReplyNumAndNextReqIdx(const Vector *sysCapablitys, int32 startIdx, int32 *nextRequestIdx)
{
    int32 registerNum = 0;
    int32 size = VECTOR_Num(sysCapablitys);
    int32 i = startIdx;
    for (; i < size && registerNum < MAX_SYSCAP_NUM_PER_REPLY; i++) {
        SysCapImpl *serviceImpl = (SysCapImpl *)VECTOR_At(sysCapablitys, i);
        if (serviceImpl->isRegister == FALSE) {
            continue;
        }
        registerNum++;
    }
    *nextRequestIdx = i;
    return registerNum;
}

void ProcGetAllSysCap(const SamgrServer *server, IpcIo *req, IpcIo *reply)
{
    uint32_t startIdx = IpcIoPopUint32(req);
    MUTEX_Lock(server->sysCapMtx);
    Vector *sysCapablitys = &(server->sysCapabilitys);
    int32 size = VECTOR_Num(sysCapablitys);
    if (size == INVALID_INDEX) {
        IpcIoPushInt32(reply, EC_FAILURE);
        IpcIoPushBool(reply, TRUE);
        IpcIoPushUint32(reply, startIdx);
        IpcIoPushUint32(reply, 0);
        MUTEX_Unlock(server->sysCapMtx);
        return;
    }
    int32 nextRequestIdx = startIdx;
    int32 replyNum = GetReplyNumAndNextReqIdx(sysCapablitys, startIdx, &nextRequestIdx);
    HILOG_DEBUG(HILOG_MODULE_SAMGR, "ProcGetAllSysCap replyNum: %d, size: %d, startIdx: %u, nextRequestIdx: %d",
                replyNum, size, startIdx, nextRequestIdx);
    IpcIoPushInt32(reply, EC_SUCCESS);
    // indicate is the last reply
    IpcIoPushBool(reply, nextRequestIdx == size);
    // indicate is the next start idx
    IpcIoPushUint32(reply, nextRequestIdx);
    IpcIoPushUint32(reply, replyNum);
    int32 cnt = 0;
    int32 i = startIdx;
    for (; i < size && cnt < replyNum; i++) {
        SysCapImpl *serviceImpl = (SysCapImpl *)VECTOR_At(sysCapablitys, i);
        if (serviceImpl->isRegister == FALSE) {
            continue;
        }
        IpcIoPushString(reply, serviceImpl->name);
        cnt++;
    }
    MUTEX_Unlock(server->sysCapMtx);
}

static int ProcSysCap(SamgrServer *server, int32 option, void *origin, IpcIo *req, IpcIo *reply)
{
    if (CanRequest(origin) == FALSE) {
        HILOG_ERROR(HILOG_MODULE_SAMGR, "ProcSysCap no permission");
        IpcIoPushInt32(reply, EC_PERMISSION);
        return EC_PERMISSION;
    }
    if (option != OP_PUT && option != OP_GET && option != OP_ALL) {
        IpcIoPushInt32(reply, EC_INVALID);
        return EC_INVALID;
    }
    HILOG_DEBUG(HILOG_MODULE_SAMGR, "ProcSysCap option: %d begin", option);
    if (option == OP_PUT) {
        int32 ret = ProcAddSysCap(server, req);
        IpcIoPushInt32(reply, ret);
    } else if (option == OP_GET) {
        BOOL ret = ProcGetSysCap(server, req);
        IpcIoPushInt32(reply, EC_SUCCESS);
        IpcIoPushBool(reply, ret);
    } else if (option == OP_ALL) {
        ProcGetAllSysCap(server, req, reply);
    } else {
        HILOG_WARN(HILOG_MODULE_SAMGR, "ProcSysCap error option: %d", option);
        IpcIoPushInt32(reply, EC_INVALID);
        return EC_INVALID;
    }
    HILOG_DEBUG(HILOG_MODULE_SAMGR, "ProcSysCap end");
    return EC_SUCCESS;
}

static int RegisterSamgrEndpoint(const IpcContext* context, SvcIdentity* identity)
{
    int ret = SetSaManager(context, MAX_SA_SIZE);
    if (ret != LITEIPC_OK) {
        HILOG_FATAL(HILOG_MODULE_SAMGR, "Set sa manager<%d> failed!", ret);
        // Set sa manager failed! We need restart to recover
        exit(-ret);
    }
    identity->handle = SAMGR_HANDLE;
    identity->token = SAMGR_TOKEN;
    identity->cookie = SAMGR_COOKIE;
    return EC_SUCCESS;
}

static int OnEndpointExit(const IpcContext *context, void* ipcMsg, IpcIo* data, void* argv)
{
    (void)data;
    if (ipcMsg != NULL) {
        FreeBuffer(context, ipcMsg);
    }
    pid_t pid = (pid_t)((uintptr_t)argv);
    Request request = {0};
    request.msgId = MSG_CLEAN;
    request.msgValue = pid;
    int retry = RETRY_TIMES;
    int ret = EC_INVALID;
    while (retry > 0) {
        ret = SAMGR_SendRequest(&g_server.identity, &request, NULL);
        if (ret == EC_SUCCESS) {
            break;
        }
        sleep(RETRY_INTERVAL);
        --retry;
    }
#ifdef __LINUX__
    PidHandle handle;
    int err = SASTORA_FindHandleByPid(&g_server.store, pid, &handle);
    if (err != INVALID_INDEX) {
        BinderRelease(context, handle.handle);
    }
#endif
    HILOG_ERROR(HILOG_MODULE_SAMGR, "IPC pid<%d> exit! send clean request retry(%d), ret(%d)!", pid, retry, ret);
    return EC_SUCCESS;
}

static IpcAuthInterface *GetIpcAuthInterface(void)
{
    IpcAuthInterface *ipcAuth = NULL;
    IUnknown *iUnknown = SAMGR_GetInstance()->GetFeatureApi(PERMISSION_SERVICE, IPCAUTH);
    if (iUnknown == NULL) {
        HILOG_ERROR(HILOG_MODULE_SAMGR, "Get IpcAuthInterface: IUnknown NULL");
        return NULL;
    }
    (void)iUnknown->QueryInterface(iUnknown, DEFAULT_VERSION, (void **)&ipcAuth);
    return ipcAuth;
}

static cJSON *GetJsonStream()
{
    const char *path = "/etc/system_capability.json";
    struct stat fileInfo;
    int32_t size = 0;

    if (stat(path, &fileInfo) != 0 || (size = fileInfo.st_size) == 0) {
        return NULL;
    }

    int32_t fp = open(path, O_RDONLY, S_IRUSR);
    if (fp < 0) {
        return NULL;
    }

    char *json = (char *)SAMGR_Malloc(size * sizeof(char));
    if (json == NULL) {
        close(fp);
        return NULL;
    }
    if (read(fp, json, size * sizeof(char)) != size * sizeof(char)) {
        SAMGR_Free(json);
        close(fp);
        return NULL;
    }
    close(fp);

    cJSON *root = cJSON_Parse(json);
    SAMGR_Free(json);
    json = NULL;
    return root;
}

static void ParseSysCap(void)
{
    cJSON *root = GetJsonStream();
    if (root == NULL) {
        HILOG_ERROR(HILOG_MODULE_SAMGR, "ParseSysCap GetJsonStream failed!");
        return;
    }
    cJSON *sysCaps = cJSON_GetObjectItem(root, "systemCapability");
    if (!cJSON_IsArray(sysCaps)) {
        cJSON_Delete(root);
        HILOG_ERROR(HILOG_MODULE_SAMGR, "ParseSysCap format failed!");
        return;
    }
    int32_t size = cJSON_GetArraySize(sysCaps);
    int32_t sysCapNum = 0;
    for (int32_t i = 0; i < size; i++) {
        if (sysCapNum >= MAX_SYSCAP_NUM) {
            HILOG_ERROR(HILOG_MODULE_SAMGR, "ParseSycCapMap system capability exceed");
            break;
        }
        cJSON *item = cJSON_GetArrayItem(sysCaps, i);
        if (!cJSON_IsObject(item)) {
            continue;
        }
        cJSON *name = cJSON_GetObjectItem(item, "name");
        cJSON *isRegister = cJSON_GetObjectItem(item, "register-on-startup");
        if (!cJSON_IsString(name) || !cJSON_IsBool(isRegister)) {
            continue;
        }
        char *nameStr = cJSON_GetStringValue(name);
        if (VECTOR_FindByKey(&(g_server.sysCapabilitys), nameStr) != INVALID_INDEX) {
            HILOG_WARN(HILOG_MODULE_SAMGR, "Duplicate system capability %s register!", nameStr);
            continue;
        }
        SysCapImpl *impl = (SysCapImpl *)SAMGR_Malloc(sizeof(SysCapImpl));
        if (impl == NULL) {
            continue;
        }
        if (strcpy_s(impl->name, sizeof(impl->name), cJSON_GetStringValue(name)) != EC_SUCCESS) {
            SAMGR_Free(impl);
            continue;
        }
        impl->isRegister = cJSON_IsTrue(isRegister);
        if (VECTOR_Add(&(g_server.sysCapabilitys), impl) == INVALID_INDEX) {
            SAMGR_Free(impl);
            HILOG_ERROR(HILOG_MODULE_SAMGR, "system capability %s register failed!", impl->name);
            continue;
        }
        sysCapNum++;
    }
    cJSON_Delete(root);
}
