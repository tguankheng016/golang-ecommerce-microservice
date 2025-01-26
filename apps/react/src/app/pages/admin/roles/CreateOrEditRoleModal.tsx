import { zodResolver } from '@hookform/resolvers/zod';
import { BusyButton, CancelButton } from '@shared/components/buttons';
import { CustomMessage, ValidationMessage } from '@shared/components/form-validation';
import { DefaultModalProps } from '@shared/components/modals';
import { ChangePermissionWarningBox, PermissionTree } from '@shared/components/permission-tree';
import { APIClient } from "@shared/service-proxies";
import { HumaCreateRoleRequestBody, HumaUpdateRoleRequestBody } from '@shared/service-proxies/identity-service-proxies';
import { SwalNotifyService } from '@shared/sweetalert2';
import { InputText } from 'primereact/inputtext';
import { classNames } from 'primereact/utils';
import { useEffect, useRef, useState } from 'react';
import { Modal, Tab, Tabs } from 'react-bootstrap';
import { Controller, useForm } from 'react-hook-form';
import { z } from 'zod';

interface RoleModalProps extends DefaultModalProps {
    roleId?: number;
}

const CreateOrEditRoleDtoSchema = z.object({
    id: z.number().nullable().optional(),
    name: z.string().min(1, { message: CustomMessage.required }),
    isDefault: z.boolean().optional(),
});

type FormData = z.infer<typeof CreateOrEditRoleDtoSchema>;

const CreateOrEditRoleModal = ({ roleId, show, handleClose, handleSave }: RoleModalProps) => {
    const [saving, setSaving] = useState(false);
    const [isEdit, setIsEdit] = useState(false);
    const [loading, setLoading] = useState(false);
    const [grantedPermissions, setGrantedPermissions] = useState<string[]>([]);

    const { register, control, reset, handleSubmit, getValues, formState: { errors } } = useForm<FormData>({
        resolver: zodResolver(CreateOrEditRoleDtoSchema),
        mode: "onTouched",
        defaultValues: {
            name: "",
        },
    });

    useEffect(() => {
        const abortController = new AbortController();
        const signal = abortController.signal;

        if (show) {
            setLoading(true);
            const roleService = APIClient.getRoleService();

            roleService.getRoleById(roleId ?? 0, signal)
                .then((res) => {
                    setIsEdit(res.role?.id != undefined && res.role.id > 0);
                    reset({ ...res.role });
                    setGrantedPermissions(res.role?.grantedPermissions ?? []);
                }).finally(() => {
                    setLoading(false);
                });
        }

        return () => {
            abortController.abort();
        };
    }, [show]);

    const submitHandler = (data: FormData) => {
        setSaving(true);

        const roleService = APIClient.getRoleService();

        if (isEdit) {
            // Update role
            const input = HumaUpdateRoleRequestBody.fromJS(data);
            input.grantedPermissions = grantedPermissions;

            roleService.updateRole(input).then((res) => {
                SwalNotifyService.info('Saved successfully');
                closeHandler();
                handleSave?.();
            }).finally(() => {
                setSaving(false);
            });
        } else {
            // Create new role
            const input = HumaCreateRoleRequestBody.fromJS(data);
            input.grantedPermissions = grantedPermissions;

            roleService.createRole(input).then((res) => {
                SwalNotifyService.info('Saved successfully');
                closeHandler();
                handleSave?.();
            }).finally(() => {
                setSaving(false);
            });
        }
    };

    const closeHandler = () => {
        resetForm();
        handleClose();
    }

    const resetForm = () => {
        reset();
        setIsEdit(false);
        setSaving(false);
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
                onSubmit={handleSubmit(data => submitHandler(data))}
                className="form w-100 fv-plugins-bootstrap5 fv-plugins-framework"
                id="kt_create_or_edit_role_form"
            >
                <div className="modal-header">
                    <h5 className="modal-title">
                        {isEdit ? (
                            <span>Edit Role: {getValues("name")}</span>
                        ) : (
                            <span>Create New Role</span>
                        )}
                    </h5>
                    <button type="button" className="btn-close" onClick={closeHandler} aria-label="Close"></button>
                </div>
                <div className="modal-body">
                    <Tabs defaultActiveKey="general">
                        <Tab eventKey="general" title="General" className="p-3 pt-6">
                            <div className="mb-5">
                                <label className="form-label required">Role name</label>
                                <Controller
                                    name="name"
                                    control={control}
                                    render={({ field, fieldState }) => (
                                        <InputText
                                            id={field.name}
                                            {...field}
                                            className={classNames({ 'p-invalid': fieldState.invalid })}
                                            type="text"
                                            autoFocus
                                        />
                                    )}
                                />
                                <ValidationMessage errorMessage={errors?.name?.message} />
                            </div>
                            <div className="mb-5">
                                <label className="form-check form-check-custom form-check-solid py-2">
                                    <input
                                        {...register('isDefault')}
                                        type="checkbox"
                                        className="form-check-input"
                                    />
                                    <span className="fw-semibold ps-2 fs-6">
                                        Is Default?
                                    </span>
                                </label>
                            </div>
                        </Tab>
                        <Tab eventKey="permissions" title="Permissions" className="p-3 pt-6">
                            <PermissionTree
                                loading={loading}
                                show={show}
                                grantedPermissions={grantedPermissions}
                                setGrantedPermissions={setGrantedPermissions}
                            />
                            {
                                isEdit && <ChangePermissionWarningBox />
                            }
                        </Tab>
                    </Tabs>
                </div>
                <div className="modal-footer">
                    <CancelButton disabled={saving} onClick={closeHandler} />
                    <BusyButton
                        className="btn btn-primary fw-bold"
                        busyIf={saving}
                    >
                        <i className="fa fa-save me-2"></i>
                        <span>Save</span>
                    </BusyButton>
                </div>
            </form>
        </Modal>
    )
}

export default CreateOrEditRoleModal