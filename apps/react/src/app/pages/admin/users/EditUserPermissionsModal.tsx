import { BusyButton, CancelButton } from '@shared/components/buttons';
import { DefaultModalProps } from "@shared/components/modals";
import { PermissionTree } from '@shared/components/permission-tree';
import APIClient from '@shared/service-proxies/api-client';
import { UpdateUserPermissionDto, UserDto } from '@shared/service-proxies/identity-service-proxies';
import { SwalNotifyService } from '@shared/sweetalert2';
import { FormEvent, useEffect, useState } from 'react';
import { Modal } from 'react-bootstrap';

interface UserPermissionModalProps extends DefaultModalProps {
    userDto: UserDto;
}

const EditUserPermissionsModal = ({ userDto, show, handleClose }: UserPermissionModalProps) => {
    const [saving, setSaving] = useState(false);
    const [loading, setLoading] = useState(false);
    const [resetting, setResetting] = useState(false);
    const [grantedPermissions, setGrantedPermissions] = useState<string[]>([]);

    useEffect(() => {
        const abortController = new AbortController();
        const signal = abortController.signal;

        if (show) {
            setLoading(true);
            const userService = APIClient.getUserService();

            userService.getUserPermissions(userDto.id ?? 0, signal)
                .then((response) => {
                    if (response.items) {
                        setGrantedPermissions(response.items);
                    }
                }).finally(() => {
                    setLoading(false);
                });
        }

        return () => {
            abortController.abort();
        };
    }, [show, userDto]);

    const submitHandler = (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        setSaving(true);

        const input = new UpdateUserPermissionDto();
        input.grantedPermissions = grantedPermissions;

        const identityService = APIClient.getUserService();
        identityService.updateUserPermissions(userDto.id ?? 0, input)
            .then(() => {
                SwalNotifyService.info('Saved successfully');
                handleClose();
            }).finally(() => {
                setSaving(false);
            });
    }

    const closeHandler = () => {
        handleClose();
    }

    const resetPermissionHandler = () => {
        setResetting(true);

        const identityService = APIClient.getUserService();
        identityService.resetUserPermissions(userDto.id ?? 0)
            .then(() => {
                SwalNotifyService.info('Saved successfully');
                handleClose();
            }).finally(() => {
                setResetting(false);
            });
    }

    return (
        <Modal
            show={show}
            onHide={handleClose}
            backdrop="static"
            keyboard={false}
            size="lg"
        >
            <form
                onSubmit={submitHandler}
                className="form w-100 fv-plugins-bootstrap5 fv-plugins-framework"
                id="kt_edit_user_permissions_form"
            >
                <div className="modal-header">
                    <h5 className="modal-title">
                        <span>Permissions: {userDto.userName}</span>
                    </h5>
                    <button type="button" className="btn-close" onClick={closeHandler} aria-label="Close"></button>
                </div>
                <div className="modal-body">
                    <PermissionTree
                        loading={loading}
                        show={show}
                        grantedPermissions={grantedPermissions}
                        setGrantedPermissions={setGrantedPermissions}
                    />
                    <div className="notice d-flex bg-light-warning rounded border-warning border border-dashed mb-9 p-6 mt-5">
                        <i className="ki-duotone ki-design-1 fs-2tx text-primary me-4"></i>
                        <div className="d-flex flex-stack flex-grow-1">
                            <div className="fw-semibold">
                                <div className="fs-6 text-gray-700">
                                    If you are changing your own permissions, you may need to refresh page (F5) to take effect of permission changes on your own screen.
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div className="modal-footer">
                    <BusyButton
                        className="btn btn-primary fw-bold"
                        busyIf={resetting}
                        type="button"
                        onClick={resetPermissionHandler}
                        disabled={saving}
                    >
                        <i className="fa fa-sync me-2"></i>
                        <span>Reset</span>
                    </BusyButton>
                    <CancelButton disabled={saving} onClick={closeHandler} />
                    <BusyButton
                        className="btn btn-primary fw-bold"
                        busyIf={saving}
                        disabled={resetting}
                    >
                        <i className="fa fa-save me-2"></i>
                        <span>Save</span>
                    </BusyButton>
                </div>
            </form>
        </Modal>
    )
}

export default EditUserPermissionsModal